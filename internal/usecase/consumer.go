package usecase

import (
	"context"
	"errors"
	"fmt"

	"gitlab.com/fluxx1on_group/event_message_service/internal/entity"
)

const DeletedError string = "Entity deleted error"

type ConsumerUseCase struct {
	msg      MessageRepo
	cli      ClientRepo
	mail     MailingRepo
	sender   Sender
	producer AdditionalProducer
}

func NewConsumer(
	msgRepo MessageRepo,
	cliRepo ClientRepo,
	mailRepo MailingRepo,
	sender Sender,
	producer AdditionalProducer,
) *ConsumerUseCase {
	return &ConsumerUseCase{
		msg:      msgRepo,
		cli:      cliRepo,
		mail:     mailRepo,
		sender:   sender,
		producer: producer,
	}
}

func (u *ConsumerUseCase) ConsumeGroup(ctx context.Context, mailing *entity.Mailing) (
	*entity.MailingStats, error,
) {
	var try = 0

	err := u.checkMailing(ctx, mailing)
	if err != nil {
		return nil, fmt.Errorf("ConsumerUseCase - ConsumeGroup() - checkMailing(): %w",
			errors.New(DeletedError))
	}

	clients, err := u.cli.ReadByFilter(ctx, mailing)
	if err != nil {
		return nil, fmt.Errorf("ConsumerUseCase - ConsumeGroup(): %w", err)
	}

	reserveClients := u.sendToClients(ctx, mailing, clients, try)
	if len(reserveClients) > 0 {
		err = u.publishReserved(ctx, &entity.MailingWithClients{
			Mailing: mailing,
			Clients: reserveClients,
			Try:     try,
		})
		if err != nil {
			return nil, fmt.Errorf("ConsumerUseCase - ConsumeGroup() - publishReserved(): %w", err)
		}
	}

	stats, err := u.mail.ReadWithMessages(ctx, mailing)
	if err != nil {
		return nil, fmt.Errorf("ConsumerUseCase - ConsumeGroup(): %w", err)
	}

	return stats, nil
}

func (u *ConsumerUseCase) ConsumePool(ctx context.Context, mwc *entity.MailingWithClients) (
	*entity.MailingStats, error,
) {
	err := u.checkMailing(ctx, mwc.Mailing)
	if err != nil {
		return nil, fmt.Errorf("ConsumerUseCase - ConsumePool() - checkMailing(): %w",
			errors.New(DeletedError))
	}

	clients, err := u.cli.ReadByFilter(ctx, mwc.Mailing)
	if err != nil {
		return nil, fmt.Errorf("ConsumerUseCase - ConsumePool(): %w", err)
	}

	reserveClients := u.sendToClients(ctx, mwc.Mailing, clients, mwc.Try)
	if len(reserveClients) > 0 {
		err = u.publishReserved(ctx, &entity.MailingWithClients{
			Mailing: mwc.Mailing,
			Clients: reserveClients,
			Try:     mwc.Try + 1,
		})
		if err != nil {
			return nil, fmt.Errorf("ConsumerUseCase - ConsumePool() - publishReserved(): %w", err)
		}
	}

	stats, err := u.mail.ReadWithMessages(ctx, mwc.Mailing)
	if err != nil {
		return nil, fmt.Errorf("ConsumerUseCase - ConsumePool(): %w", err)
	}

	return stats, nil
}

// sendToClients tryes to Send() mailing to clients and create
// new messages in DB for each
//
// sendToClients return aborted messages to resend it later
func (u *ConsumerUseCase) sendToClients(
	ctx context.Context, mailing *entity.Mailing, clients entity.Clients, try int,
) entity.Clients {
	var reserveClients entity.Clients = make(entity.Clients, 0)

	for _, client := range clients {
		var deliveryStatus bool = true

		if client.CheckTimeZone(mailing.IntervalStart, mailing.IntervalEnd) {
			err := u.sender.Send(ctx, &entity.SendRequest{
				ID:    client.ID,
				Phone: client.PhoneNumber,
				Text:  mailing.MessageText,
			})
			if err != nil {
				deliveryStatus = false
				reserveClients = append(reserveClients, client)
			}

			msg := &entity.Message{
				Try:            try,
				DeliveryStatus: deliveryStatus,
				MailingID:      mailing.ID,
				ClientID:       client.ID,
			}
			_ = u.msg.Create(ctx, msg)

		} else {
			reserveClients = append(reserveClients, client)
		}
	}

	return reserveClients
}

// checkMailing return error if mailing was deleted.
//
// Also it revert mailing with updated attrs (if these were updated)
func (u *ConsumerUseCase) checkMailing(ctx context.Context, mailing *entity.Mailing) error {
	receivedM, err := u.mail.Read(ctx, mailing)
	if err != nil {
		return err
	}

	if &receivedM != &mailing {
		mailing = receivedM
	}

	return nil
}

// publishReserver produce mwc to Nats to resend aborted messages
func (u *ConsumerUseCase) publishReserved(ctx context.Context, mwc *entity.MailingWithClients) error {
	err := u.producer.Publish(ctx, mwc)
	if err != nil {
		return err
	}

	return nil
}

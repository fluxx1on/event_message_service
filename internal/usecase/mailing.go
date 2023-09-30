package usecase

import (
	"context"
	"fmt"

	"gitlab.com/fluxx1on_group/event_message_service/internal/entity"
)

type MailingUseCase struct {
	repo     MailingRepo
	msgRepo  MessageRepo
	producer GeneralProducer
}

func NewMailing(repo MailingRepo, msgRepo MessageRepo, producer GeneralProducer) *MailingUseCase {
	return &MailingUseCase{
		repo:     repo,
		msgRepo:  msgRepo,
		producer: producer,
	}
}

func (u *MailingUseCase) Add(ctx context.Context, mailing *entity.Mailing) error {
	err := u.repo.Create(ctx, mailing)
	if err != nil {
		return fmt.Errorf("MailingUseCase - Add(): %w", err)
	}

	err = u.producer.Publish(ctx, mailing)
	if err != nil {
		return fmt.Errorf("MailingUseCase - Add(): %w", err)
	}

	return nil
}

func (u *MailingUseCase) Patch(ctx context.Context, mailing *entity.Mailing) error {
	err := u.repo.Update(ctx, mailing)
	if err != nil {
		return fmt.Errorf("MailingUseCase - Patch(): %w", err)
	}

	return nil
}

func (u *MailingUseCase) Delete(ctx context.Context, mailing *entity.Mailing) error {
	err := u.repo.Create(ctx, mailing)
	if err != nil {
		return fmt.Errorf("MailingUseCase - Delete(): %w", err)
	}

	return nil
}

func (u *MailingUseCase) GetMailingStats(ctx context.Context) ([]*entity.MailingStats, error) {
	var (
		stats []*entity.MailingStats
	)

	allMailings, err := u.repo.ReadAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("MailingUseCase - GetMailingStats(): %w", err)
	}

	for _, mailing := range allMailings {
		stat, err := u.repo.ReadWithMessages(ctx, mailing)
		if err != nil {
			return nil, fmt.Errorf("MailingUseCase - GetMailingStats(): %w", err)
		}

		stats = append(stats, stat)
	}

	return stats, nil
}

func (u *MailingUseCase) GetMessagesByMailing(ctx context.Context, mailing *entity.Mailing) (
	entity.Messages, error,
) {
	var err error

	msgs, err := u.msgRepo.ReadByMailing(ctx, mailing)
	if err != nil {
		return nil, fmt.Errorf("MailingUseCase - GetMessageByMailing(): %w", err)
	}

	return msgs, nil
}

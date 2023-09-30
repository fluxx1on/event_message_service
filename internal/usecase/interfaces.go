package usecase

import (
	"context"

	"gitlab.com/fluxx1on_group/event_message_service/internal/entity"
)

//go:generate mockgen -source=interfaces.go -destination=./mocks_test.go -package=usecase_test

// UseCases
type (
	// Client -
	Client interface {
		Add(context.Context, *entity.Client) error
		Patch(context.Context, *entity.Client) error
		Delete(context.Context, *entity.Client) error
	}

	// Mailing -
	Mailing interface {
		Add(context.Context, *entity.Mailing) error
		Patch(context.Context, *entity.Mailing) error
		Delete(context.Context, *entity.Mailing) error

		GetMailingStats(context.Context) ([]*entity.MailingStats, error)
		GetMessagesByMailing(context.Context, *entity.Mailing) (entity.Messages, error)
	}

	Consumer interface {
		ConsumeGroup(context.Context, *entity.Mailing) (*entity.MailingStats, error)
		ConsumePool(context.Context, *entity.MailingWithClients) (*entity.MailingStats, error)
	}
)

// Repositories
type (
	// ClientRepo -
	ClientRepo interface {
		Create(context.Context, *entity.Client) error
		Update(context.Context, *entity.Client) error
		Delete(context.Context, *entity.Client) error

		Read(context.Context, *entity.Client) (*entity.Client, error)
		ReadByFilter(context.Context, *entity.Mailing) (entity.Clients, error)
	}

	// MailingRepo -
	MailingRepo interface {
		Create(context.Context, *entity.Mailing) error
		Update(context.Context, *entity.Mailing) error
		Delete(context.Context, *entity.Mailing) error

		ReadWithMessages(context.Context, *entity.Mailing) (*entity.MailingStats, error)
		Read(context.Context, *entity.Mailing) (*entity.Mailing, error)
		ReadAll(context.Context) (entity.Mailings, error)
	}

	// MessageRepo -
	MessageRepo interface {
		Create(context.Context, *entity.Message) error

		ReadByMailing(context.Context, *entity.Mailing) (entity.Messages, error)
		Read(context.Context, *entity.Message) (*entity.Message, error)
	}
)

// MQ
type (
	// MailingProducer -
	GeneralProducer interface {
		Publish(context.Context, *entity.Mailing) error
	}

	// ClientProducer -
	AdditionalProducer interface {
		Publish(context.Context, *entity.MailingWithClients) error
	}
)

// External API
type (
	// SenderWebAPI -
	Sender interface {
		Send(context.Context, *entity.SendRequest) error
	}
)

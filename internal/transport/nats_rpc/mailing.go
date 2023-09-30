package nats_rpc

import (
	"context"
	"errors"
	"log/slog"

	"github.com/nats-io/nats.go"
	"gitlab.com/fluxx1on_group/event_message_service/internal/entity"
	"gitlab.com/fluxx1on_group/event_message_service/internal/usecase"
	"gitlab.com/fluxx1on_group/event_message_service/pkg/nats_server/mod"
)

type mailingConsumer struct {
	c usecase.Consumer

	ctx context.Context
}

func newMailingConsumer(r map[string]mod.HandlerGroup, c usecase.Consumer, subjects ...string) {
	m := &mailingConsumer{
		c:   c,
		ctx: context.Background(),
	}

	if len(subjects) < 2 {
		panic("not enough subjects to start nats router")
	}
	{
		r[subjects[0]] = mod.HandlerGroup{Get: m.sendGroup(), Clean: m.clean()}
		r[subjects[1]] = mod.HandlerGroup{Get: m.sendPool(), Clean: m.clean()}
	}
}

func (m *mailingConsumer) sendGroup() mod.MsgTimeHandler {
	return func(msg *nats.Msg) (bool, mod.Task) {
		var rtask mod.Task

		// Unmarshalling
		var mailing entity.Mailing
		err := mailing.UnmarshalJSON(msg.Data)
		if err != nil {
			slog.Error("Mailing unmarshal error",
				slog.String("Subject", msg.Subject),
				slog.String("ErrorMsg", err.Error()))
			m.clean()(msg)
			return true, rtask
		}

		// Delay
		rtask.SetIn(mailing.DateTimeStart, mailing.DateTimeEnd)
		if rtask.In() != 0 {
			rtask.Msg = msg
			return false, rtask
		}

		// Consumption
		s, err := m.c.ConsumeGroup(m.ctx, &mailing)
		if errors.Is(err, errors.New(usecase.DeletedError)) {
			err = msg.Term()
		} else if err != nil {
			err = msg.Nak()
		} else {
			err = msg.Ack()
		}

		if err != nil {
			slog.Error("Internal unexpected error",
				slog.String("Subject", msg.Subject),
				slog.String("ErrorMsg", err.Error()))
		}

		slog.Info("Messages sended",
			slog.String("Subject", msg.Subject),
			slog.Group(
				"Mailing stats",
				s.MailingID, s.DateTimeStart, s.DateTimeEnd, s.Succesed, s.Failed),
		)
		return true, rtask
	}
}

func (m *mailingConsumer) sendPool() mod.MsgTimeHandler {
	return func(msg *nats.Msg) (bool, mod.Task) {
		var rtask mod.Task

		// Unmarshalling
		var mwc entity.MailingWithClients
		err := mwc.UnmarshalJSON(msg.Data)
		if err != nil {
			slog.Error("Mailing unmarshal error",
				slog.String("Subject", msg.Subject),
				slog.String("ErrorMsg", err.Error()))
			m.clean()(msg)
			return true, rtask
		}

		// Delay
		rtask.SetIn(mwc.Mailing.DateTimeStart, mwc.Mailing.DateTimeEnd)
		if rtask.In() != 0 {
			rtask.Msg = msg
			return false, rtask
		}

		// Consumption
		s, err := m.c.ConsumePool(m.ctx, &mwc)
		if errors.Is(err, errors.New(usecase.DeletedError)) {
			err = msg.Term()
		} else if err != nil {
			err = msg.Nak()
		} else {
			err = msg.Ack()
		}

		if err != nil {
			slog.Error("Internal unexpected error",
				slog.String("Subject", msg.Subject),
				slog.String("ErrorMsg", err.Error()))
		}

		slog.Info("Messages sended",
			slog.String("Subject", msg.Subject),
			slog.Group(
				"Intermediate stats",
				s.MailingID, s.DateTimeStart, s.DateTimeEnd, s.Succesed, s.Failed),
		)
		return true, rtask
	}
}

func (m *mailingConsumer) clean() mod.MsgTermHandler {
	return func(msg *nats.Msg) {
		if err := msg.Term(); err != nil {
			slog.Error("mailingConsumer - clean()", slog.String("ErrorMsg", err.Error()))
		}
	}
}

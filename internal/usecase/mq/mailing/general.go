package mailing

import (
	"context"
	"fmt"

	"gitlab.com/fluxx1on_group/event_message_service/internal/entity"
	"gitlab.com/fluxx1on_group/event_message_service/pkg/nats_server"
)

type GeneralProducer struct {
	conn *nats_server.Connection

	subj string
}

func NewGeneral(conn *nats_server.Connection, subj string) *GeneralProducer {
	return &GeneralProducer{conn, subj}
}

func (p *GeneralProducer) Publish(ctx context.Context, general *entity.Mailing) error {
	data, err := general.MarshalJSON()
	if err != nil {
		return fmt.Errorf("GeneralProducer - Publish(): %w", err)
	}

	err = p.conn.Publish(p.subj, data)
	if err != nil {
		return fmt.Errorf("GeneralProducer - Publish(): %w", err)
	}

	return nil
}

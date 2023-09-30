package mailing

import (
	"context"
	"fmt"

	"gitlab.com/fluxx1on_group/event_message_service/internal/entity"
	"gitlab.com/fluxx1on_group/event_message_service/pkg/nats_server"
)

type AdditionalProducer struct {
	conn *nats_server.Connection

	subj string
}

func NewAdditional(conn *nats_server.Connection, subj string) *AdditionalProducer {
	return &AdditionalProducer{conn, subj}
}

func (p *AdditionalProducer) Publish(ctx context.Context, mwc *entity.MailingWithClients) error {
	data, err := mwc.MarshalJSON()
	if err != nil {
		return fmt.Errorf("AdditionalProducer - Publish(): %w", err)
	}

	err = p.conn.Publish(p.subj, data)
	if err != nil {
		return fmt.Errorf("AdditionalProducer - Publish(): %w", err)
	}

	return nil
}

package postgres

import (
	"context"
	"fmt"

	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5/pgxpool"
	"gitlab.com/fluxx1on_group/event_message_service/internal/entity"
)

const tableMessage = "message"

type MessageRepo struct {
	Builder squirrel.StatementBuilderType
	conn    *pgxpool.Pool
}

// NewMessage - MessageRepo constructor
func NewMessage(conn *pgxpool.Pool) *MessageRepo {
	return &MessageRepo{
		Builder: squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar),
		conn:    conn,
	}
}

// Create -.
func (r *MessageRepo) Create(ctx context.Context, message *entity.Message) error {
	query, args, err := r.Builder.
		Insert(tableMessage).
		Columns("try", "delivery_status", "mailing_id", "client_id").
		Values(
			message.Try,
			message.DeliveryStatus,
			message.MailingID,
			message.ClientID,
		).
		ToSql()
	if err != nil {
		return fmt.Errorf("MessageRepo - Create(): %w", err)
	}

	_, err = r.conn.Exec(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("MessageRepo - Create(): %w", err)
	}

	return nil
}

// ReadByMailing -.
func (r *MessageRepo) ReadByMailing(ctx context.Context, mailing *entity.Mailing) (
	entity.Messages, error,
) {
	query, _, err := r.Builder.
		Select("*").
		From(tableMessage).
		Where(squirrel.Eq{"mailing_id": mailing.ID}).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("MessageRepo - ReadByMailing(): %w", err)
	}

	rows, err := r.conn.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("MessageRepo - ReadByMailing(): %w", err)
	}

	var ms entity.Messages
	for rows.Next() {
		var m entity.Message
		err = rows.Scan(
			&m.ID, &m.DateTimeCreation, &m.Try, &m.DeliveryStatus, &m.MailingID, &m.ClientID,
		)
		if err != nil {
			return nil, fmt.Errorf("MessageRepo - ReadByMailing(): %w", err)
		}
		ms = append(ms, &m)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("MessageRepo - ReadByMailing(): %w", err)
	}

	return ms, nil
}

// Read -.
func (r *MessageRepo) Read(ctx context.Context, message *entity.Message) (*entity.Message, error) {
	query, _, err := r.Builder.
		Select("*").
		From(tableMessage).
		Where(squirrel.Eq{"id": message.ID}).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("MessageRepo - Read(): %w", err)
	}

	var m entity.Message
	err = r.conn.QueryRow(ctx, query).Scan(
		&m.ID, &m.DateTimeCreation, &m.Try, &m.DeliveryStatus, &m.MailingID, &m.ClientID,
	)
	if err != nil {
		return nil, fmt.Errorf("MessageRepo - Read(): %w", err)
	}

	return &m, nil
}

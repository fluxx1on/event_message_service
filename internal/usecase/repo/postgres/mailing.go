package postgres

import (
	"context"
	"fmt"

	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5/pgxpool"
	"gitlab.com/fluxx1on_group/event_message_service/internal/entity"
)

const tableMailing = "mailing"

type MailingRepo struct {
	Builder squirrel.StatementBuilderType
	conn    *pgxpool.Pool
}

func NewMailing(conn *pgxpool.Pool) *MailingRepo {
	return &MailingRepo{
		Builder: squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar),
		conn:    conn,
	}
}

func (r *MailingRepo) Create(ctx context.Context, mailing *entity.Mailing) error {
	query, args, err := r.Builder.
		Insert(tableMailing).
		Columns("message_text", "mobile_operator_code", "tag",
			"datetime_start", "datetime_end", "interval_start", "interval_end").
		Values(
			mailing.MessageText,
			mailing.MobileOperator,
			mailing.Tag,
			mailing.DateTimeStart,
			mailing.DateTimeEnd,
			mailing.IntervalStart,
			mailing.IntervalEnd,
		).
		ToSql()
	if err != nil {
		return fmt.Errorf("MailingRepo - Create(): %w", err)
	}

	_, err = r.conn.Exec(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("MailingRepo - Create(): %w", err)
	}

	return nil
}

func (r *MailingRepo) Update(ctx context.Context, mailing *entity.Mailing) error {
	builder := r.Builder.
		Update(tableMailing).
		Where(squirrel.Eq{"id": mailing.ID})

	if mailing.MessageText != "" {
		builder = builder.Set("message_text", mailing.MessageText)
	}
	if mailing.MobileOperator != "" {
		builder = builder.Set("mobile_operator_code", mailing.MobileOperator)
	}
	if mailing.Tag != "" {
		builder = builder.Set("tag", mailing.Tag)
	}
	if mailing.FilterChoice != "" {
		builder = builder.Set("filter_choice", mailing.FilterChoice)
	}
	if mailing.DateTimeStart.IsZero() {
		builder = builder.Set("datetime_start", mailing.DateTimeStart)
	}
	if mailing.DateTimeEnd.IsZero() {
		builder = builder.Set("datetime_end", mailing.DateTimeEnd)
	}
	if mailing.IntervalStart.IsZero() {
		builder = builder.Set("interval_start", mailing.IntervalStart)
	}
	if mailing.IntervalEnd.IsZero() {
		builder = builder.Set("interval_end", mailing.IntervalEnd)
	}

	query, args, err := builder.ToSql()
	if err != nil {
		return fmt.Errorf("MailingRepo - Update(): %w", err)
	}

	_, err = r.conn.Exec(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("MailingRepo - Update(): %w", err)
	}

	return nil
}

func (r *MailingRepo) Delete(ctx context.Context, mailing *entity.Mailing) error {
	query, args, err := r.Builder.
		Delete(tableMailing).
		Where(squirrel.Eq{"id": mailing.ID}).
		ToSql()
	if err != nil {
		return fmt.Errorf("MailingRepo - Delete(): %w", err)
	}

	_, err = r.conn.Exec(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("MailingRepo - Delete(): %w", err)
	}

	return nil
}

func (r *MailingRepo) ReadWithMessages(ctx context.Context, mailing *entity.Mailing) (*entity.MailingStats, error) {
	query, _, err := r.Builder.
		Select().
		Columns("m.id AS mailing_id",
			"MIN(msg.date_time_creation)",
			"MAX(msg.date_time_creation)",
			"COUNT(CASE WHEN msg.delivery_status = 1 THEN 1 ELSE 0 END) AS successed",
			"COUNT(CASE WHEN msg.delivery_status = 0 THEN 1 ELSE 0 END) AS failed",
		).
		From(tableMailing + " m").
		Where(squirrel.Eq{"m.id": mailing.ID}).
		LeftJoin("message msg ON msg.mailing_id = m.id").
		GroupBy("m.id, m.datetime_start, m.datetime_end").
		OrderBy("m.datetime_end").
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("MailingRepo - ReadWithMessages(): %w", err)
	}

	var m entity.MailingStats
	err = r.conn.QueryRow(ctx, query).Scan(
		&m.MailingID, &m.DateTimeStart, &m.DateTimeEnd, &m.Succesed, &m.Failed,
	)
	if err != nil {
		return nil, fmt.Errorf("MailingRepo - ReadWithMessages(): %w", err)
	}

	return &m, nil
}

// Read - Select all the fields from mailing table.
//
// Important: if current mailing.ID doesn't compare with any row
// in table then Read() return error. It needs to check mailing by deletion.
func (r *MailingRepo) Read(ctx context.Context, mailing *entity.Mailing) (*entity.Mailing, error) {
	query, _, err := r.Builder.
		Select("*").
		From(tableMailing).
		Where(squirrel.Eq{"id": mailing.ID}).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("MailingRepo - Read(): %w", err)
	}

	var m entity.Mailing
	err = r.conn.QueryRow(ctx, query).Scan(
		&m.ID, &m.MessageText, &m.MobileOperator, &m.Tag, &m.FilterChoice,
		&m.DateTimeStart, &m.DateTimeEnd, &m.IntervalStart, &m.IntervalEnd,
	)
	if err != nil {
		return nil, fmt.Errorf("MailingRepo - Read(): %w", err)
	}

	return &m, nil
}

func (r *MailingRepo) ReadAll(ctx context.Context) (entity.Mailings, error) {
	query, _, err := r.Builder.
		Select("*").
		From(tableMailing).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("MailingRepo - ReadAll(): %w", err)
	}

	rows, err := r.conn.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("MailingRepo - ReadAll(): %w", err)
	}

	var ms entity.Mailings
	for rows.Next() {
		var m entity.Mailing
		err = rows.Scan(
			&m.ID, &m.MessageText, &m.MobileOperator, &m.Tag, &m.FilterChoice,
			&m.DateTimeStart, &m.DateTimeEnd, &m.IntervalStart, &m.IntervalEnd,
		)
		if err != nil {
			return nil, fmt.Errorf("MailingRepo - ReadAll(): %w", err)
		}
		ms = append(ms, &m)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("MailingRepo - ReadAll(): %w", err)
	}

	return ms, nil
}

package postgres

import (
	"context"
	"fmt"

	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5/pgxpool"
	"gitlab.com/fluxx1on_group/event_message_service/internal/entity"
)

const tableClient string = "client"

type ClientRepo struct {
	Builder squirrel.StatementBuilderType
	conn    *pgxpool.Pool
}

func NewClient(conn *pgxpool.Pool) *ClientRepo {
	return &ClientRepo{
		Builder: squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar),
		conn:    conn,
	}
}

func (r *ClientRepo) Create(ctx context.Context, client *entity.Client) error {
	query, args, err := r.Builder.
		Insert(tableClient).
		Columns("mobile_operator_code", "phone_number", "tag", "time_zone").
		Values(
			client.MobileOperator,
			client.PhoneNumber,
			client.Tag,
			client.TimeZone,
		).
		ToSql()
	if err != nil {
		return fmt.Errorf("ClientRepo - Create(): %w", err)
	}

	_, err = r.conn.Exec(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("ClientRepo - Create(): %w", err)
	}

	return nil
}

func (r *ClientRepo) Update(ctx context.Context, client *entity.Client) error {
	builder := r.Builder.
		Update(tableClient).
		Where(squirrel.Eq{"id": client.ID})

	if client.PhoneNumber != 0 {
		builder = builder.Set("phone_number", client.PhoneNumber)
	}
	if client.MobileOperator != 0 {
		builder = builder.Set("mobile_operator_code", client.MobileOperator)
	}
	if client.Tag != "" {
		builder = builder.Set("tag", client.Tag)
	}
	if client.TimeZone != 0 {
		builder = builder.Set("time_zone", client.TimeZone)
	}

	query, args, err := builder.ToSql()
	if err != nil {
		return fmt.Errorf("ClientRepo - Update(): %w", err)
	}

	_, err = r.conn.Exec(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("ClientRepo - Update(): %w", err)
	}

	return nil
}

func (r *ClientRepo) Delete(ctx context.Context, client *entity.Client) error {
	query, args, err := r.Builder.
		Delete(tableClient).
		Where(squirrel.Eq{"id": client.ID}).
		ToSql()
	if err != nil {
		return fmt.Errorf("ClientRepo - Delete(): %w", err)
	}

	_, err = r.conn.Exec(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("ClientRepo - Delete(): %w", err)
	}

	return nil
}

func (r *ClientRepo) Read(ctx context.Context, client *entity.Client) (*entity.Client, error) {
	query, _, err := r.Builder.
		Select("*").
		From(tableClient).
		Where(squirrel.Eq{"id": client.ID}).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("ClientRepo - Read(): %w", err)
	}

	var c entity.Client
	err = r.conn.QueryRow(ctx, query).Scan(
		&c.ID, &c.MobileOperator, &c.PhoneNumber, &c.Tag, &c.TimeZone,
	)
	if err != nil {
		return nil, fmt.Errorf("ClientRepo - Read(): %w", err)
	}

	return &c, nil
}

func (r *ClientRepo) ReadByFilter(ctx context.Context, mailing *entity.Mailing) (entity.Clients, error) {
	var where map[string]interface{} = make(map[string]interface{}, 1)

	switch mailing.FilterChoice {
	case "code":
		where["mobile_operator_code"] = squirrel.Eq{"mobile_operator_code": mailing.MobileOperator}
	case "tag":
		where["tag"] = squirrel.Eq{"tag": mailing.Tag}
	}

	query, _, err := r.Builder.
		Select("id", "phone_number", "time_zone").
		From(tableClient).
		Where(where).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("ClientRepo - ReadByFilter(): %w", err)
	}

	var cs entity.Clients
	rows, err := r.conn.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("ClientRepo - ReadByFiler(): %w", err)
	}

	for rows.Next() {
		var c entity.Client
		err = rows.Scan(
			&c.ID, &c.PhoneNumber, &c.TimeZone,
		)
		if err != nil {
			return nil, fmt.Errorf("ClientRepo - ReadByFiler(): %w", err)
		}
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("ClientRepo - ReadByFilter(): %w", err)
	}

	return cs, nil
}

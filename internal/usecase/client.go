package usecase

import (
	"context"
	"fmt"

	"gitlab.com/fluxx1on_group/event_message_service/internal/entity"
)

type ClientUseCase struct {
	repo ClientRepo
}

func NewClient(repo ClientRepo) *ClientUseCase {
	return &ClientUseCase{repo}
}

func (u *ClientUseCase) Add(ctx context.Context, client *entity.Client) error {
	err := u.repo.Create(ctx, client)
	if err != nil {
		return fmt.Errorf("ClientUseCase - Add(): %w", err)
	}

	return nil
}

func (u *ClientUseCase) Patch(ctx context.Context, client *entity.Client) error {
	err := u.repo.Update(ctx, client)
	if err != nil {
		return fmt.Errorf("ClientUseCase - Patch(): %w", err)
	}

	return nil
}

func (u *ClientUseCase) Delete(ctx context.Context, client *entity.Client) error {
	err := u.repo.Delete(ctx, client)
	if err != nil {
		return fmt.Errorf("ClientUseCase - Delete(): %w", err)
	}

	return nil
}

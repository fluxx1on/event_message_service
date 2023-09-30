package external_test

import (
	"context"
	"testing"

	"github.com/go-playground/assert/v2"
	"gitlab.com/fluxx1on_group/event_message_service/internal/entity"
	"gitlab.com/fluxx1on_group/event_message_service/internal/usecase/external"
)

func TestSend(t *testing.T) {
	list := []*entity.SendRequest{
		{
			ID:    0,
			Phone: 0,
			Text:  "string",
		},
		{
			ID:    10000000,
			Phone: 79771203244,
			Text:  "string",
		},
	}

	sender := external.New()

	// test 1
	{
		err := sender.Send(context.Background(), list[0])
		assert.Equal(t, err, nil)
	}

	// test 2
	{
		err := sender.Send(context.Background(), list[1])
		assert.Equal(t, err, nil)
	}
}

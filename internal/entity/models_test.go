package entity_test

import (
	"log"
	"testing"
	"time"

	"github.com/go-playground/assert/v2"
	"gitlab.com/fluxx1on_group/event_message_service/internal/entity"
)

func TestCheckTimeZone(t *testing.T) {
	client := entity.Client{
		TimeZone: 12,
	}

	now := time.Now()
	// test 1
	{
		duration := 4 * time.Hour
		end := now.Add(duration)
		ok := client.CheckTimeZone(now, end)
		assert.Equal(t, ok, true)
	}
	// test 2
	{
		duration := 2 * time.Hour
		end := now.Add(duration)
		log.Printf("%s : %s", now, end)
		ok := client.CheckTimeZone(now, end)
		assert.Equal(t, ok, false)
	}
}

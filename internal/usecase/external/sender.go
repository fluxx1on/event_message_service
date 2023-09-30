package external

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"

	"gitlab.com/fluxx1on_group/event_message_service/internal/entity"
)

const (
	url = "https://probe.fbrq.cloud/v1/send/"

	token = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MjU3MDA1NjIsImlzcyI6ImZhYnJpcXVlIiwibmFtZSI6Imh0dHBzOi8vdC5tZS9sdW1vdXNmIn0.jMIi06w6qu2abk9g4mg-8BoWD8kGv5xy_TtzzRMmXlk" // Unsafe token storing
)

const Err400 = "Status BadRequest"

type Sender struct {
	*http.Client
}

func New() *Sender {
	return &Sender{
		Client: &http.Client{},
	}
}

func (s *Sender) Send(_ context.Context, body *entity.SendRequest) error {
	// Marshalling
	reqBody, err := body.MarshalJSON()
	if err != nil {
		return s.sendErr("Marshalling", err)
	}

	// Request building
	req, err := http.NewRequest(http.MethodPost, msg_url(body.ID), bytes.NewReader(reqBody))
	if err != nil {
		return s.sendErr("Request building", err)
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	// Getting response
	resp, err := s.Do(req)
	if err != nil {
		return s.sendErr("Getting response", err)
	}
	defer resp.Body.Close()

	// CheckStatus
	if resp.StatusCode == 400 {
		return s.sendErr("CheckStatus", errors.New(Err400))
	}

	// Parsing
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		slog.Error("", slog.String("ErrorMsg", s.sendErr("Parsing", err).Error()))
		return s.sendErr("Parsing", err)
	}

	// Unmarshalling
	var sendResp entity.SendResponse
	err = sendResp.UnmarshalJSON(respBody)
	if err != nil {
		slog.Error("", slog.String("ErrorMsg", s.sendErr("Unmarshalling", err).Error()))
		return s.sendErr("Unmarshalling", err)
	}

	// Checking response
	if sendResp.Code != 0 {
		return s.sendErr("Checking repsonse", errors.New("Unepected response {code}"))
	}

	return nil
}

func (s *Sender) sendErr(path string, err error) error {
	return fmt.Errorf("Sender - Send() - %s: %w", path, err)
}

func msg_url(msgId int64) string {
	return fmt.Sprintf("%s%d", url, msgId)
}

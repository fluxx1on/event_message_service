package logger

import (
	"context"
	"encoding/json"
	"io"
	"log"

	"log/slog"

	"github.com/fatih/color"
)

type ColorfulHandler struct {
	// implements base struct
	slog.Handler

	logger     *log.Logger
	logJournal *log.Logger
	attrs      []slog.Attr
}

func NewColorfulHandler(out, journal io.Writer, opts *slog.HandlerOptions) *ColorfulHandler {
	h := &ColorfulHandler{
		Handler:    slog.NewTextHandler(out, opts),
		logger:     log.New(out, "", 0),
		logJournal: log.New(journal, "", 0),
	}

	return h
}

func (h *ColorfulHandler) Handle(_ context.Context, r slog.Record) error {
	level := r.Level.String() + ":"

	switch r.Level {
	case slog.LevelDebug:
		level = color.HiBlackString(level)
	case slog.LevelInfo:
		level = color.GreenString(level)
	case slog.LevelWarn:
		level = color.YellowString(level)
	case slog.LevelError:
		level = color.RedString(level)
	}

	fields := make(map[string]interface{}, r.NumAttrs())

	r.Attrs(func(a slog.Attr) bool {
		if a.Value.Kind() == slog.KindGroup {
			nestedFields := map[string]interface{}{}
			for _, subAttr := range a.Value.Group() {
				nestedFields[subAttr.Key] = subAttr.Value.Any()
			}
			fields[a.Key] = nestedFields
		} else {
			fields[a.Key] = a.Value.Any()
		}
		return true
	})

	for _, a := range h.attrs {
		fields[a.Key] = a.Value.Any()
	}

	var b []byte
	var err error
	if len(fields) > 0 {
		b, err = json.MarshalIndent(fields, "", "  ")
		if err != nil {
			return err
		}
	}

	timeStr := r.Time.Format("[15:04:05.000]")
	msg := color.CyanString(r.Message)

	h.logJournal.Println(
		timeStr,
		level,
		r.Message,
		string(b),
	)

	h.logger.Println(
		timeStr,
		level,
		msg,
		color.WhiteString(string(b)),
	)

	return nil
}

func (h *ColorfulHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &ColorfulHandler{
		Handler:    h.Handler,
		logger:     h.logger,
		logJournal: h.logJournal,
		attrs:      attrs,
	}
}

func (h *ColorfulHandler) WithGroup(name string) slog.Handler {
	return &ColorfulHandler{
		Handler:    h.Handler.WithGroup(name),
		logger:     h.logger,
		logJournal: h.logJournal,
	}
}

package nats_server

import (
	"log/slog"
	"time"

	"github.com/nats-io/nats.go"
)

// Config - config to connect with nats.
type Config struct {
	URL string
}

// Connection - nats connection that using JetStream basically
type Connection struct {
	*nats.Conn
}

func OpenConnection(cfg Config) *Connection {
	options := nats.GetDefaultOptions()
	options.Url = cfg.URL
	options.AllowReconnect = true
	options.Timeout = 5 * time.Second
	options.MaxReconnect = 10
	options.RetryOnFailedConnect = true

	conn, err := options.Connect()
	if err != nil {
		panic("nats connection refused")
	}

	server := &Connection{conn}

	return server
}

func (c *Connection) Close(timeout time.Duration) {
	time.Sleep(timeout)

	c.Conn.Close()
	slog.Info("NATS disconnected with timeout")
}

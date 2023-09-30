package server

import (
	"time"

	server "gitlab.com/fluxx1on_group/event_message_service/pkg/nats_server"
	"gitlab.com/fluxx1on_group/event_message_service/pkg/nats_server/mod"
)

const _defaultTimeout = 2 * time.Second

type Server struct {
	conn *server.Connection
	stop chan struct{}

	manager mod.Manager

	timeout time.Duration
}

func New(
	conn *server.Connection,
	router map[string]mod.HandlerGroup,
	modtype mod.ManagerType,
	opts ...Option,
) *Server {

	stop := make(chan struct{})

	var manager mod.Manager

	switch modtype {
	case mod.TaskType:
		manager = mod.NewTaskManager(conn, router, stop)
	}

	server := &Server{
		conn:    conn,
		stop:    stop,
		manager: manager,
		timeout: _defaultTimeout,
	}

	// Custom options
	for _, opt := range opts {
		opt(server)
	}

	return server
}

func (s *Server) StartWorkers() {
	s.manager.Subscribe()
}

func (s *Server) Shutdown() {
	close(s.stop)

	s.conn.Close(s.timeout)
}

package server

import "time"

type Option func(*Server)

func Timeout(timeout time.Duration) Option {
	return func(s *Server) {
		s.timeout = timeout
	}
}

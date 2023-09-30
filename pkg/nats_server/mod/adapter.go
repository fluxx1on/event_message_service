package mod

import (
	"github.com/nats-io/nats.go"
)

type Manager interface {
	Subscribe()
}

type ManagerType string

const (
	TaskType ManagerType = "TaskManager"
)

// TaskManager

type HandlerGroup struct {
	Get   MsgTimeHandler
	Clean MsgTermHandler
}

type MsgTimeHandler func(*nats.Msg) (bool, Task)

type MsgTermHandler nats.MsgHandler

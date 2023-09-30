package nats_rpc

import (
	"gitlab.com/fluxx1on_group/event_message_service/internal/usecase"
	"gitlab.com/fluxx1on_group/event_message_service/pkg/nats_server/mod"
)

func NewRouter(consumer usecase.Consumer, subjects ...string) map[string]mod.HandlerGroup {
	routes := make(map[string]mod.HandlerGroup)
	{
		newMailingConsumer(routes, consumer, subjects...)
	}

	return routes
}

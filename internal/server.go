package internal

import (
	"context"
	"net/http"

	"log/slog"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"gitlab.com/fluxx1on_group/event_message_service/internal/config"
	v1 "gitlab.com/fluxx1on_group/event_message_service/internal/transport/http/v1"
	"gitlab.com/fluxx1on_group/event_message_service/internal/transport/nats_rpc"
	"gitlab.com/fluxx1on_group/event_message_service/internal/usecase"
	"gitlab.com/fluxx1on_group/event_message_service/internal/usecase/external"
	"gitlab.com/fluxx1on_group/event_message_service/internal/usecase/mq/mailing"
	"gitlab.com/fluxx1on_group/event_message_service/internal/usecase/repo/postgres"
	"gitlab.com/fluxx1on_group/event_message_service/pkg/nats_server"
	"gitlab.com/fluxx1on_group/event_message_service/pkg/nats_server/mod"
	"gitlab.com/fluxx1on_group/event_message_service/pkg/nats_server/server"
)

type Node struct {
	dbConn     *pgxpool.Pool
	httpServer *http.Server
	natsServer *server.Server
}

func (n *Node) Start(cfg *config.Config) {
	var err error

	// ___ Connections ___

	// PostgreSQL
	n.dbConn, err = pgxpool.New(context.Background(), cfg.PostgreSQL.URL)
	if err != nil {
		slog.Error("PostgreSQL unreached", slog.String("ErrorMsg", err.Error()))
		panic("startup")
	}

	// Nats
	conn := nats_server.OpenConnection(nats_server.Config{
		URL: cfg.Nats.URL,
	})

	// ___ Infrastructure Layer ___

	// Repositories
	var (
		clientRepo  *postgres.ClientRepo  = postgres.NewClient(n.dbConn)
		mailingRepo *postgres.MailingRepo = postgres.NewMailing(n.dbConn)
		messageRepo *postgres.MessageRepo = postgres.NewMessage(n.dbConn)
	)

	// Producers
	var (
		mailingProducer *mailing.GeneralProducer    = mailing.NewGeneral(conn, cfg.Nats.Subjects[0])
		clientProducer  *mailing.AdditionalProducer = mailing.NewAdditional(conn, cfg.Nats.Subjects[1])
	)

	// External API
	var sender *external.Sender = external.New()

	// ___ UseCase Layer ___

	// -
	var (
		client  *usecase.ClientUseCase  = usecase.NewClient(clientRepo)
		mailing *usecase.MailingUseCase = usecase.NewMailing(
			mailingRepo, messageRepo, mailingProducer,
		)
		consumer *usecase.ConsumerUseCase = usecase.NewConsumer(
			messageRepo, clientRepo, mailingRepo, sender, clientProducer,
		)
	)

	// ___ Transport Layer ___

	// NatsServer - Consumer server
	natsRouter := nats_rpc.NewRouter(consumer, cfg.Nats.Subjects...)
	n.natsServer = server.New(conn, natsRouter, mod.TaskType)

	// HTTP Server - API
	handler := gin.New()
	v1.NewRouter(handler, client, mailing)
	n.httpServer = &http.Server{
		Addr:    cfg.Addr,
		Handler: handler,
	}

	// Servers starting
	go func() {
		if err := n.httpServer.ListenAndServe(); err != nil {
			slog.Error("Failed to serve", slog.String("ErrorMsg", err.Error()))
		}
	}()

	slog.Info("HTTP server started.", slog.String("HTTP Address", cfg.Addr))

	// Start consumers
	go n.natsServer.StartWorkers()

	slog.Info("NATS server started.", slog.String("NATS Address", cfg.Nats.Host))
}

func (n *Node) Stop() {
	n.httpServer.Close()
	slog.Info("HTTP server shutted down")

	n.natsServer.Shutdown()

	n.dbConn.Close()
	slog.Info("PostgreSQL disconnected")
}

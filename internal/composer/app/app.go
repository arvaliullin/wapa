package app

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/arvaliullin/wapa/internal/broker"
	"github.com/arvaliullin/wapa/internal/delivery"
	"github.com/arvaliullin/wapa/internal/delivery/handlers"
	"github.com/arvaliullin/wapa/internal/persistence"
	"github.com/arvaliullin/wapa/internal/storage"
	"github.com/nats-io/nats.go"
)

// ComposerService - сервис composer.
type ComposerService struct {
	Config           *ComposerConfig
	HttpService      delivery.HttpService
	NATSConnection   *nats.Conn
	DesignPublisher  *broker.DesignPublisher
	ResultSubscriber *broker.ResultSubscriber
}

// NewComposerService создаёт новый экземпляр ComposerService.
func NewComposerService(config *ComposerConfig) (*ComposerService, error) {
	nc, err := nats.Connect(config.NatsURL)
	if err != nil {
		return nil, err
	}

	experimentRepository := persistence.NewExperimentRepository(config.DbConnection)

	return &ComposerService{
		Config:           config,
		HttpService:      delivery.NewEchoHttpService(),
		NATSConnection:   nc,
		DesignPublisher:  broker.NewDesignPublisher(nc, config.NatsSubjectRunner),
		ResultSubscriber: broker.NewResultSubscriber(nc, config.NatsSubjectResult, experimentRepository),
	}, nil
}

// Run запускает сервис и ожидает сигнала завершения.
func (service *ComposerService) Run(ctx context.Context) {
	go func() {
		designRepo := &persistence.DesignRepository{
			DbConnection: service.Config.DbConnection,
		}

		designStorage := &storage.DesignStorage{
			DataPath: service.Config.DataPath,
		}

		handlers.RegisterDesignHandler(service.HttpService,
			designRepo,
			designStorage)

		handlers.RegisterExperimentHandler(service.HttpService,
			designRepo,
			service.DesignPublisher)

		benchmarkRepo := &persistence.BenchmarkRepository{DbConnection: service.Config.DbConnection}

		handlers.RegisterBenchmarkHandler(service.HttpService, benchmarkRepo)

		log.Printf("Starting HTTP server on %s", service.Config.Address)
		if err := service.HttpService.Start(service.Config.Address); err != nil && err != http.ErrServerClosed {
			log.Printf("HTTP server error: %v", err)
		}
	}()

	go service.ResultSubscriber.Start()

	<-ctx.Done()

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := service.HttpService.Shutdown(shutdownCtx); err != nil {
		log.Printf("HTTP shutdown error: %v", err)
	}

	if err := service.NATSConnection.Drain(); err != nil {
		log.Printf("NATS drain error: %v", err)
	}
}

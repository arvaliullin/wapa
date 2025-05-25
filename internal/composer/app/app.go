package app

import (
	"log"

	"github.com/arvaliullin/wapa/internal/broker"
	"github.com/arvaliullin/wapa/internal/delivery"
	"github.com/arvaliullin/wapa/internal/delivery/handlers"
	"github.com/arvaliullin/wapa/internal/persistence"
	"github.com/arvaliullin/wapa/internal/storage"
	"github.com/nats-io/nats.go"
)

// ComposerService
type ComposerService struct {
	Config           *ComposerConfig
	HttpService      delivery.HttpService
	NATSConnection   *nats.Conn
	DesignPublisher  *broker.DesignPublisher
	ResultSubscriber *broker.ResultSubscriber
}

func NewComposerService(config *ComposerConfig) *ComposerService {

	nc, err := nats.Connect(config.NatsURL)
	if err != nil {
		log.Fatalf("Error connecting to NATS: %v", err)
	}

	experimentRepository := persistence.NewExperimentRepository(config.DbConnection)

	return &ComposerService{
		Config:           config,
		HttpService:      delivery.NewEchoHttpService(),
		NATSConnection:   nc,
		DesignPublisher:  broker.NewDesignPublisher(nc, config.NatsSubjectRunner),
		ResultSubscriber: broker.NewResultSubscriber(nc, config.NatsSubjectResult, experimentRepository),
	}
}

func (service *ComposerService) Run() {

	go func() {
		log.Printf("Starting HTTP server on %s", service.Config.Address)

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

		service.HttpService.Start(service.Config.Address)
	}()

	go func() {
		service.ResultSubscriber.Start()
	}()

	defer service.NATSConnection.Close()

	select {}
}

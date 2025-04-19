package app

import (
	"encoding/json"
	"log"

	"github.com/arvaliullin/wapa/internal/domain"
	"github.com/arvaliullin/wapa/internal/storage"
	"github.com/nats-io/nats.go"
)

type RunnerService struct {
	Config         *RunnerConfig
	Storage        *storage.ExperimentStorage
	NATSConnection *nats.Conn
}

func NewRunnerService(config *RunnerConfig) *RunnerService {
	st, err := storage.NewExperimentStorage(config.DataPath)
	if err != nil {
		log.Fatalf("Storage init error: %v", err)
	}
	nc, err := nats.Connect(config.NatsURL)
	if err != nil {
		log.Fatalf("Error connecting to NATS: %v", err)
	}
	return &RunnerService{
		Config:         config,
		Storage:        st,
		NATSConnection: nc,
	}
}

func (service *RunnerService) Publish(message []byte) error {
	return service.NATSConnection.Publish(service.Config.NatsSubjectResult, message)
}

func (service *RunnerService) Start() {
	_, err := service.NATSConnection.Subscribe(service.Config.NatsSubjectRunner,
		func(msg *nats.Msg) {
			var design domain.Design

			if err := json.Unmarshal(msg.Data, &design); err != nil {
				log.Printf("Ошибка декодирования задачи: %v", err)
				return
			}

			experiment := service.ExecuteTask(design)

			resBytes, err := json.Marshal(experiment)
			if err != nil {
				log.Printf("Ошибка сериализации результата: %v", err)
				return
			}
			if err := service.Publish(resBytes); err != nil {
				log.Printf("Ошибка публикации результата: %v", err)
			}
		})
	if err != nil {
		log.Fatalln(err)
	}
	log.Printf("Слушаем задания в %s ...", service.Config.NatsSubjectRunner)
}

func (service *RunnerService) ExecuteTask(design domain.Design) domain.Experiment {

	var wasmPath, jsPath string

	wasmPath, err := service.Storage.DownloadFile(design.ID, "wasm", service.Config.ComposerAddress)
	if err == nil {
		defer service.Storage.DeleteFile(wasmPath)
	}
	jsPath, err = service.Storage.DownloadFile(design.ID, "js", service.Config.ComposerAddress)
	if err == nil {
		defer service.Storage.DeleteFile(jsPath)
	}

	//TODO: Реализовать запуск Command

	log.Printf("Запустить тест для: %v\n", design)
	log.Printf("wasmPath: %s\n", wasmPath)
	log.Printf("wasmPath: %s\n", wasmPath)

	experiment := domain.Experiment{}

	return experiment
}

func (service *RunnerService) Run() {
	go service.Start()
	defer service.NATSConnection.Close()
	select {}
}

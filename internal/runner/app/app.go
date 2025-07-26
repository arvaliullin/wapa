package app

import (
	"encoding/json"
	"log"
	"os"

	"github.com/arvaliullin/wapa/internal/domain"
	"github.com/arvaliullin/wapa/internal/storage"
	"github.com/nats-io/nats.go"
)

// RunnerService получает план эксперимента, запускает его, отправляет результат в composer.
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
			var design domain.DesignPayload

			if err := json.Unmarshal(msg.Data, &design); err != nil {
				log.Printf("Ошибка декодирования задачи: %v", err)
				return
			}

			experiment, err := service.Execute(design)
			if err != nil {
				log.Printf("ошибка выполнения эксперимента: %v", err)
				return
			}
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

func (service *RunnerService) Execute(design domain.DesignPayload) (domain.Experiment, error) {
	var wasmPath, jsPath string

	wasmPath, err := service.Storage.DownloadFile(design, "wasm", service.Config.ComposerAddress)
	if err != nil {
		log.Printf("ошибка при скачивании %s ...", err)
	}
	jsPath, err = service.Storage.DownloadFile(design, "js", service.Config.ComposerAddress)
	if err != nil {
		log.Printf("ошибка при скачивании %s ...", err)
	}

	command := Command{
		DesignPayload:      design,
		HyperfinePath:      "hyperfine",
		HyperfineResultDir: os.TempDir(),
		NodePath:           "bun",
		ScriptPath:         "/opt/wapa/scripts/cpp.js",
		WasmPath:           wasmPath,
		JsPath:             jsPath,
	}

	experiment := command.Run()
	return experiment, nil
}

func (service *RunnerService) Run() {
	go service.Start()
	defer service.NATSConnection.Close()
	select {}
}

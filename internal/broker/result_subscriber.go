package broker

import (
	"encoding/json"
	"log"

	"github.com/arvaliullin/wapa/internal/domain"
	"github.com/arvaliullin/wapa/internal/persistence"
	"github.com/nats-io/nats.go"
)

type ResultSubscriber struct {
	NATSConnection       *nats.Conn
	Subject              string
	ExperimentRepository *persistence.ExperimentRepository
}

func NewResultSubscriber(natsConnection *nats.Conn,
	subject string, experimentRepository *persistence.ExperimentRepository) *ResultSubscriber {
	return &ResultSubscriber{
		NATSConnection:       natsConnection,
		Subject:              subject,
		ExperimentRepository: experimentRepository,
	}
}

func (es *ResultSubscriber) Start() {
	_, err := es.NATSConnection.Subscribe(es.Subject, func(msg *nats.Msg) {
		var exp domain.Experiment
		if err := json.Unmarshal(msg.Data, &exp); err != nil {
			log.Printf("[%s] Failed to unmarshal experiment: %v\n", es.Subject, err)
			return
		}

		if id, err := es.ExperimentRepository.Create(exp); err != nil {
			log.Printf("[%s] Failed to save experiment: %v\n", es.Subject, err)
			return
		} else {
			log.Printf("[%s] Experiment saved: %+v\n", es.Subject, id)
		}
	})

	if err != nil {
		log.Fatalln(err)
	}
}

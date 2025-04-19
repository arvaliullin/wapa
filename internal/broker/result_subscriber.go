package broker

import (
	"log"

	"github.com/nats-io/nats.go"
)

type ResultSubscriber struct {
	NATSConnection *nats.Conn
	Subject        string
}

func NewResultSubscriber(natsConnection *nats.Conn,
	subject string) *ResultSubscriber {
	return &ResultSubscriber{NATSConnection: natsConnection,
		Subject: subject}
}

func (es *ResultSubscriber) Start() {

	_, err := es.NATSConnection.Subscribe(es.Subject, func(msg *nats.Msg) {
		//TODO: Получает результаты эксперимента
		// Создает запись о результате в таблице composer.experiment
		log.Printf("[%s] Received a message: %s\n", es.Subject, string(msg.Data))
	})

	if err != nil {
		log.Fatalln(err)
	}
}

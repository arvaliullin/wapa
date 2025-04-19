package broker

import (
	"log"

	"github.com/nats-io/nats.go"
)

type DesignPublisher struct {
	NATSConnection *nats.Conn
	Subject        string
}

func NewDesignPublisher(conn *nats.Conn, subj string) *DesignPublisher {
	return &DesignPublisher{
		NATSConnection: conn,
		Subject:        subj,
	}
}

// Publish публикует данные о плане эксперимента исполнителям
func (dp *DesignPublisher) Publish(message []byte) error {
	return dp.NATSConnection.Publish(dp.Subject, message)
}

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

package broker

import (
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

// Publish публикует данные о плане эксперимента исполнителям.
func (dp *DesignPublisher) Publish(message []byte) error {
	return dp.NATSConnection.Publish(dp.Subject, message)
}

package main

import (
	"log"

	"github.com/nats-io/nats.go"
)

func main() {

	nc, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		log.Fatalln(err)
	}
	defer nc.Close()

	subj := "runners"
	_, err = nc.Subscribe(subj, func(msg *nats.Msg) {
		log.Printf("[%s] Received a message: %s\n", subj, string(msg.Data))
	})

	if err != nil {
		log.Fatalln(err)
	}

	log.Printf("[%s] Waiting for messages. Press Ctrl+C to exit.", subj)

	select {}
}

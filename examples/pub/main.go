package main

import (
	"encoding/json"
	"log"

	"github.com/nats-io/nats.go"
)

func main() {
	nc, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		log.Fatalf("Connection error %v", err)
	}
	defer nc.Close()

	message := map[string]string{
		"type":    "example",
		"content": "This is a JSON message",
	}

	jsonData, err := json.Marshal(message)
	if err != nil {
		log.Fatalf("JSON marshaling error: %v", err)
	}

	err = nc.Publish("runners", jsonData)
	if err != nil {
		log.Fatalf("Publish error: %v", err)
	}
}

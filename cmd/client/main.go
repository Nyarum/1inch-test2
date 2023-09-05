package main

import (
	"1inch/internal"
	"log"
	"os"

	"github.com/nats-io/nats.go"
)

func main() {
	natsURL := os.Getenv("NATS_URL")
	if natsURL == "" {
		natsURL = nats.DefaultURL
	}

	natsClient, err := internal.NewNATS(natsURL)
	if err != nil {
		log.Fatal(err)
	}

	externalMsg := natsClient.Pull(internal.DefaultPath)
	for msg := range externalMsg {
		log.Printf("Address: %v, Token0: %v, Token1: %v", msg.Prices.Pool, msg.Prices.Token0, msg.Prices.Token1)
	}
}

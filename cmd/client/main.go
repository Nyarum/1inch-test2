package main

import (
	"1inch/internal"
	"log"

	"github.com/nats-io/nats.go"
)

func main() {
	natsClient, err := internal.NewNATS(nats.DefaultURL)
	if err != nil {
		log.Fatal(err)
	}

	externalMsg := natsClient.Pull(internal.DefaultPath)
	for msg := range externalMsg {
		log.Printf("Address: %v, Token0: %v, Token1: %v", msg.Prices.Pool, msg.Prices.Token0, msg.Prices.Token1)
	}
}

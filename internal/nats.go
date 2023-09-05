package internal

import (
	"encoding/json"
	"log"

	"github.com/nats-io/nats.go"
)

type NATS struct {
	client *nats.Conn
}

func NewNATS(url string) (*NATS, error) {
	nc, err := nats.Connect(url)
	if err != nil {
		return nil, err
	}

	return &NATS{
		client: nc,
	}, nil
}

func (n *NATS) Push(path string, msg Message) error {
	buf, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	return n.client.Publish(path, buf)
}

func (n *NATS) Pull(path string) <-chan Message {
	ch := make(chan Message, 1)
	n.client.Subscribe(path, func(m *nats.Msg) {
		msg := Message{}
		err := json.Unmarshal(m.Data, &msg)
		if err != nil {
			log.Printf("can't unmarshal nats message, err - %v\n", err)
			return
		}

		ch <- msg
	})

	return ch
}

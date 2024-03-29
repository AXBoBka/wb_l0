package nats

import (
	"log"

	"github.com/nats-io/nats.go"
	"github.com/nats-io/stan.go"
)

func NewStanServer() stan.Conn {
	nc, err := nats.Connect("localhost:4222")
	if err != nil {
		panic(err)
	}
	log.Println("Успешно подключён к NATS")

	sc, err := stan.Connect("test-cluster", "clienttest", stan.NatsConn(nc))
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Успешно подключён к Streaming server")

	return sc
}

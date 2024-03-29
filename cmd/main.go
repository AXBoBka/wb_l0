package main

import (
	"context"
	"log"
	"os"
	"os/signal"

	"github.com/AXBoBka/wb_l0/internal/cache"
	"github.com/AXBoBka/wb_l0/internal/nats"
	"github.com/AXBoBka/wb_l0/internal/server"
	"github.com/AXBoBka/wb_l0/internal/store"
	"github.com/nats-io/stan.go"
)

func main() {
	//open DB connection
	os.Setenv("DATABASE_URL", "postgres://postgres:1@localhost:5432/wb")
	connection := store.OpenConnection()
	defer connection.Close(context.Background())

	//initialize cache
	cache := cache.New(connection)

	//connect to nats streaming
	nats_streaming := nats.NewStanServer()

	//subscribe to channel "orders" in nats streaming
	_, err := nats_streaming.Subscribe(
		"orders",
		func(msg *stan.Msg) {
			log.Printf("Найдено сообщение: %s", string(msg.Data))
			cache.AddOrder(string(msg.Data))
			store.AddOrder(msg.Data, connection)
		},
		stan.StartWithLastReceived(),
	)

	if err != nil {
		log.Fatal(err)
	}

	//start http server
	serv := server.New(cache)
	serv.Start()

	InterruptSignal := make(chan os.Signal, 1)
	cleanProcesses := make(chan bool)
	signal.Notify(InterruptSignal, os.Interrupt)
	go func() {
		for range InterruptSignal {
			log.Println("Получен сигнал прерывания, закрываю соединение...")

			nats_streaming.Close()
			cleanProcesses <- true
		}
	}()
	<-cleanProcesses
}

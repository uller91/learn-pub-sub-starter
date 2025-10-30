package main

import (
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
	"os"
	"os/signal"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/pubsub"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/routing"
)

func main() {
	const connectionString = "amqp://guest:guest@localhost:5672/"

	connection, err := amqp.Dial(connectionString)
	if err != nil {
		fmt.Printf("Could not connect to RabbitMQ. Error: %v\n", err)
		os.Exit(1)
	}
	defer connection.Close()

	fmt.Println("Peril server successfully connected to RabbitMQ")

	ch, err := connection.Channel()
	if err != nil {
		fmt.Printf("Could not creat connection channel. Error: %v\n", err)
		os.Exit(1)
	}

	message := routing.PlayingState{
		IsPaused: true,
	}
	err = pubsub.PublishJSON(ch, routing.ExchangePerilDirect, routing.PauseKey, message)

	// wait for ctrl+c
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)
	<-signalChan
	fmt.Println("\nRabbitMQ connection closed...")

	os.Exit(0)
}

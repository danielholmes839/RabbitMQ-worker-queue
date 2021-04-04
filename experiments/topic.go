package main

import (
	"fmt"
	"time"

	"github.com/streadway/amqp"
)

func checkError(err error) {
	if err != nil {
		panic(err)
	}
}

func producer(messages []string) {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	checkError(err)

	ch, err := conn.Channel()
	checkError(err)

	ch.ExchangeDeclare(
		"stocks",
		"topic",
		false,
		false,
		false,
		false,
		nil,
	)

	for _, routingKey := range messages {
		ch.Publish(
			"stocks",
			routingKey,
			false,
			false,
			amqp.Publishing{
				ContentType: "text/plain",
				Body:        []byte("stonks"),
			},
		)
	}

}

func consumer(name, key string) {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	checkError(err)

	ch, err := conn.Channel()
	checkError(err)

	err = ch.ExchangeDeclare(
		"stocks",
		"topic",
		false,
		false,
		false,
		false,
		nil,
	)
	checkError(err)

	q, err := ch.QueueDeclare(
		"",
		false,
		false,
		true,
		false,
		nil,
	)
	checkError(err)

	ch.QueueBind(
		q.Name,
		key,
		"stocks",
		false,
		nil,
	)

	messages, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto ack
		false,  // exclusive
		false,  // no local
		false,  // no wait
		nil,    // args
	)
	checkError(err)

	for message := range messages {
		fmt.Printf("'%s' received message with routing key: %s\n", name, message.RoutingKey)
	}

}

func main() {
	go consumer("TSX consumer", "tsx.#")
	go consumer("NYSE consumer", "nyse.#")
	go consumer("TSX.KXS consumer", "tsx.kxs.*")

	time.Sleep(time.Second)
	updates := []string{"tsx.kxs.update", "tsx.shop.update", "nyse.gme.update"}
	go producer(updates)

	forever := make(chan bool)
	<-forever

}

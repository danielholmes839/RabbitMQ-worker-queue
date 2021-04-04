package main

import (
	"fmt"
	"time"

	"github.com/streadway/amqp"
)

func main() {
	conn, _ := amqp.Dial("amqp://guest:guest@localhost:5672/")
	ch, _ := conn.Channel()
	q, _ := ch.QueueDeclare("worker_output", false, false, false, false, nil)

	// Consume messages from the workers
	messages, _ := ch.Consume("worker_input", "", true, false, false, false, nil)

	for message := range messages {
		fmt.Println("STARTED:", string(message.Body))
		time.Sleep(time.Second * 5)

		ch.Publish("", q.Name, false, false, amqp.Publishing{
			ContentType: "text/plain",
			Body:        message.Body,
		})
		fmt.Println("FINISHED:", string(message.Body))
	}
}

package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/streadway/amqp"
)

func checkError(err error) {
	if err == nil {
		return
	}
	panic(err)
}

func consumeWorkerMessages(ch *amqp.Channel, db *Database) {
	// Consume messages from the workers
	messages, _ := ch.Consume("worker_output", "", true, false, false, false, nil)
	for message := range messages {
		job := NewJob(message.Body)
		db.updateJob(job.Id)

		// Logging
		fmt.Printf("Job ID: %d, Name: '%s' COMPLETED\n", job.Id, job.Name)
	}
}

func submit(ch *amqp.Channel, db *Database, q *amqp.Queue) func(c *gin.Context) {
	// Create the function for the /submit endpoint

	return func(c *gin.Context) {
		// The name
		name := c.Query("name")

		// Publish message to worker_input queue
		job := db.createJob(name)
		message := amqp.Publishing{ContentType: "text/plain", Body: job.serialize()}

		err := ch.Publish("", q.Name, false, false, message)
		checkError(err)

		// Logging
		fmt.Printf("Job ID: %d, Name: '%s' STARTED\n", job.Id, job.Name)
		c.JSON(200, job)
	}
}

func main() {
	gin.SetMode(gin.ReleaseMode)

	// Connecting to RabbitMQ
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	checkError(err)

	ch, err := conn.Channel()
	checkError(err)

	q, err := ch.QueueDeclare("worker_input", false, false, false, false, nil)
	checkError(err)

	// In-memory "database"
	db := NewDatabase()
	go consumeWorkerMessages(ch, db)

	// Endpoints
	router := gin.Default()
	router.GET("/submit", submit(ch, db, &q))
	router.GET("/view", func(c *gin.Context) {
		c.JSON(200, db.read())
	})

	router.Run()
}

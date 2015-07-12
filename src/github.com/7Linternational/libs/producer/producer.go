package producer

import (
	"fmt"
	"log"
	"math/rand"
	// "os"
	"strconv"
	// "strings"
	// "time"

	"github.com/streadway/amqp"
)

var (
	amqpURI      = "amqp://guest:guest@localhost:5672/"
	queueName    = "test-exchange"
	exchangeType = "direct"
	routingKey   = "test-key"
	reliable     = true
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
		panic(fmt.Sprintf("%s: %s", msg, err))
	}
}

func echo(msgs string) string {
	return msgs
}

func randomString(l int) string {
	bytes := make([]byte, l)
	for i := 0; i < l; i++ {
		bytes[i] = byte(randInt(65, 90))
	}

	return string(bytes)
}

func randInt(min int, max int) int {
	return min + rand.Intn(max-min)
}

func Send(body string) (res int, err error) {
	// Connect to RabbitMQ
	conn, err := amqp.Dial(amqpURI)
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	// Open Channel
	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"",    // name
		false, // durable
		false, // delete when usused
		true,  // exclusive
		false, // noWait
		nil,   // arguments
	)
	failOnError(err, "Failed to declare a queue")

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	failOnError(err, "Failed to register a consumer")

	corrId := randomString(32)

	err = ch.Publish(
		"",          // exchange
		"rpc_queue", // routing key
		false,       // mandatory
		false,       // immediate
		amqp.Publishing{
			ContentType:   "text/plain",
			CorrelationId: corrId,
			ReplyTo:       q.Name,
			Body:          []byte(body),
		})
	failOnError(err, "Failed to publish a message")

	for d := range msgs {
		if corrId == d.CorrelationId {
			res, err = strconv.Atoi(string(d.Body))
			failOnError(err, "Failed to convert body to integer")
			break
		}
	}

	return

}

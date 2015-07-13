package mq

import (
	"github.com/streadway/amqp"
	"log"
	"math/rand"
)

/**
 * [checkErr description]
 * @param  {[type]} err error         [description]
 * @return {[type]}     [description]
 */
func checkErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

/**
 *
 */
func randInt(min int, max int) int {
	return min + rand.Intn(max-min)
}

/**
 *
 */
func randomString(l int) string {
	bytes := make([]byte, l)
	for i := 0; i < l; i++ {
		bytes[i] = byte(randInt(65, 90))
	}

	return string(bytes)
}

/**
 * [Consume description]
 */
func Consume(body string) string {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	checkErr(err)
	defer conn.Close()

	ch, err := conn.Channel()
	checkErr(err)
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"",    // name
		false, // durable
		false, // delete when usused
		true,  // exclusive
		false, // noWait
		nil,   // arguments
	)
	checkErr(err)

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consume
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	checkErr(err)

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
	checkErr(err)

	var response string
	for d := range msgs {
		if corrId == d.CorrelationId {
			////////////////////////
			// Send Response back //
			////////////////////////
			response = string(d.Body)
			break
		}
	}

	return response
}

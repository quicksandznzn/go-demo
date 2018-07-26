// Created by quicksandzn@gmail.com on 2018/7/25
package main

import (
	"github.com/streadway/amqp"
	"log"
	"github.com/gosexy/to"
	yml "github.com/quicksandznzn/go-yaml"
)

func main() {
	conf, err := yml.Open("./config/conf.yml")
	conn, err := amqp.Dial(to.String(conf.Get("rabbitmq", "url")))
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"hello", // name
		false,   // durable
		false,   // delete when unused
		false,   // exclusive
		false,   // no-wait
		nil,     // arguments
	)

	failOnError(err, "Failed to declare a queue")

	publish(ch, q, "hello word")

	consume(ch, q)

}

// receive msg
func consume(channel *amqp.Channel, queue amqp.Queue) {
	msgs, err := channel.Consume(
		queue.Name, // queue
		"",         // consumer
		true,       // auto-ack
		false,      // exclusive
		false,      // no-local
		false,      // no-wait
		nil,        // args
	)
	failOnError(err, "Failed to register a consumer")

	forever := make(chan bool)

	go func() {
		for d := range msgs {
			log.Printf("Received a message: %s", d.Body)
		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}

// send msg
func publish(channel *amqp.Channel, queue amqp.Queue, body string) {
	err := channel.Publish(
		"",         // exchange
		queue.Name, // routing key
		false,      // mandatory
		false,      // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(body),
		})
	log.Printf(" [x] Sent %s", body)
	failOnError(err, "Failed to publish a message")
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}
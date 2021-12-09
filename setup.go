package lograbbitmq

import (
	"github.com/coredns/caddy"
	"github.com/coredns/coredns/core/dnsserver"
	"github.com/coredns/coredns/plugin"

	"log"

	"github.com/streadway/amqp"

	"fmt"
)

type RabbitMqChannel struct {
	Channel amqp.Channel
	Queue   amqp.Queue
}

func (rabbitMqChannel *RabbitMqChannel) Send(msg string) {
	err := rabbitMqChannel.Channel.Publish(
		"",                         // exchange
		rabbitMqChannel.Queue.Name, // routing key
		false,                      // mandatory
		false,                      // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(msg),
		})

	failOnError(err, "Failed to publish a message")
}

// var RabbitMqChannelInstance = RabbitMqChannel{}

// init registers this plugin.
func init() { plugin.Register("lograbbitmq", setup) }

// setup is the function that gets called when the config parser see the token "example". Setup is responsible
// for parsing any extra options the example plugin may have. The first token this function sees is "example".
func setup(c *caddy.Controller) error {
	fmt.Println("Lograbbitmq started")
	c.Next() // Ignore "example" and give us the next token.
	if c.NextArg() {
		// If there was another token, return an error, because we don't have any configuration.
		// Any errors returned from this setup function should be wrapped with plugin.Error, so we
		// can present a slightly nicer error message to the user.
		return plugin.Error("example", c.ArgErr())
	}

	setupRabbitMqConnection()

	// Add the Plugin to CoreDNS, so Servers can use it in their plugin chain.
	dnsserver.GetConfig(c).AddPlugin(func(next plugin.Handler) plugin.Handler {
		return LogRabbitMQ{Next: next}
	})

	// All OK, return a nil error.
	return nil
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func setupRabbitMqConnection() {
	// conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	// failOnError(err, "Failed to connect to RabbitMQ")
	// defer conn.Close()

	// ch, err := conn.Channel()
	// failOnError(err, "Failed to open a channel")
	// defer ch.Close()

	// q, err := ch.QueueDeclare(
	// 	"hello", // name
	// 	false,   // durable
	// 	false,   // delete when unused
	// 	false,   // exclusive
	// 	false,   // no-wait
	// 	nil,     // arguments
	// )
	// failOnError(err, "Failed to declare a queue")

	// RabbitMqChannelInstance.Channel = ch
	// RabbitMqChannelInstance.Queue = q
}

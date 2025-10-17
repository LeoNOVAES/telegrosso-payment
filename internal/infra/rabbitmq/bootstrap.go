package rabbitmq

import (
	"fmt"
	"log"

	"github.com/rabbitmq/amqp091-go"
	"github.com/throindev/payments/cmd/config"
)

type RabbitMQRepository struct {
	Connection *amqp091.Connection
	Channel    *amqp091.Channel
}

func NewRabbitMQRepository() *RabbitMQRepository {
	conn, channel, err := CreateChannel()
	SetupExchange(channel)
	SetupQueusAndBinds(channel)

	if err != nil {
		log.Fatalf("Falha ao conectar: %v", err)
	}

	return &RabbitMQRepository{Connection: conn, Channel: channel}
}

func CreateChannel() (*amqp091.Connection, *amqp091.Channel, error) {
	conn, err := amqp091.Dial(config.AppConfig.RabbitMQ)
	if err != nil {
		log.Fatalf("Falha ao conectar: %v", err)
	}

	ch, err := conn.Channel()

	if err != nil {
		log.Fatalf("Falha ao abrir canal: %v", err)
		return nil, nil, err
	}

	return conn, ch, nil
}

func SetupExchange(ch *amqp091.Channel) {
	err := ch.ExchangeDeclare(
		config.AppConfig.Exchange,
		"topic",
		true,
		false,
		false,
		false,
		nil,
	)

	if err != nil {
		log.Fatalf("Error to create exchange: %v", err)
	}

	log.Println("Exchange Created successfully!")
}

func SetupQueusAndBinds(ch *amqp091.Channel) error {
	for queueName, routingKeys := range config.AppConfig.Queues {
		fmt.Printf("Configuring Queue: %s, Binds: %v\n", queueName, routingKeys)

		q, err := ch.QueueDeclare(queueName, true, false, false, false, nil)
		if err != nil {
			log.Fatal(err)
		}

		for _, key := range routingKeys {
			if err := ch.QueueBind(q.Name, key, config.AppConfig.Exchange, false, nil); err != nil {
				log.Fatal(err)
			}
		}
	}

	log.Println("Queues Created successfully!")
	return nil
}

package rabbitmq

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/rabbitmq/amqp091-go"
	"github.com/throindev/payments/cmd/config"
	"github.com/throindev/payments/internal/domain"
)

type RabbitMQRepository struct {
	Connection *amqp091.Connection
	Channel    *amqp091.Channel
}

func NewRabbitMQRepository() *RabbitMQRepository {
	conn, channel, err := createChannel()
	setupExchange(channel)
	setupQueusAndBinds(channel)

	if err != nil {
		log.Fatalf("Falha ao conectar: %v", err)
	}

	return &RabbitMQRepository{Connection: conn, Channel: channel}
}

func createChannel() (*amqp091.Connection, *amqp091.Channel, error) {
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

func setupExchange(ch *amqp091.Channel) {
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

func setupQueusAndBinds(ch *amqp091.Channel) error {
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

func (r *RabbitMQRepository) Consume(queueName string, handler func(body domain.Event)) error {
	msgs, err := r.Channel.Consume(
		queueName,
		"",
		true,
		false, false, false, nil,
	)

	if err != nil {
		fmt.Printf("erro to consume queue %s", queueName)
		return err
	}

	go func() {
		for msg := range msgs {
			var event domain.Event
			if err := json.Unmarshal(msg.Body, &event); err != nil {
				log.Printf("❌ Erro ao parsear evento: %v | Body: %s", err, string(msg.Body))
				continue
			}
			handler(event)
		}
	}()

	log.Printf("⚡ Consumindo fila: %s", queueName)
	return nil
}

func (r *RabbitMQRepository) Publish(queueName string, event *domain.Event) error {
	body, err := json.Marshal(event)

	if err != nil {
		return err
	}

	return r.Channel.Publish(
		config.AppConfig.Exchange, // exchange
		queueName,                 // routing key
		false,                     // mandatory
		false,                     // immediate
		amqp091.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)
}

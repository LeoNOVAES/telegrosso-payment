package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	RabbitMQ      string
	ProviderToken string
	ProviderUrl   string
	MongoURI      string
	DbName        string
	Port          string
	Exchange      string
	Queues        map[string][]string
}

var AppConfig Config

func Load() {

	queues := map[string][]string{
		"payment_queue":      {"payment.*"},
		"bot_queue":          {"payment.confirmed", "subscription.reminder"},
		"subscription_queue": {"subscription.*"},
	}

	if err := godotenv.Load(); err != nil {
		log.Fatal("Aviso: .env não encontrado, usando variáveis do sistema")
	}

	AppConfig = Config{
		RabbitMQ:      os.Getenv("RABBITMQ_URL"),
		ProviderToken: os.Getenv("PROVIDER_TOKEN"),
		ProviderUrl:   os.Getenv("PROVIDER_URL"),
		MongoURI:      os.Getenv("MONGODB_URI"),
		DbName:        os.Getenv("DB_NAME"),
		Port:          os.Getenv("PORT"),
		Exchange:      "events_topic",
		Queues:        queues,
	}
}

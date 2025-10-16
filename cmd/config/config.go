package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	MercadoPagoToken string
	MongoURI         string
	DbName           string
	Port             string
}

var AppConfig Config

func Load() {
	if err := godotenv.Load("./../.env"); err != nil {
		log.Fatal("Aviso: .env não encontrado, usando variáveis do sistema")
	}

	AppConfig = Config{
		MercadoPagoToken: os.Getenv("MERCADO_PAGO_TOKEN"),
		MongoURI:         os.Getenv("MONGODB_URI"),
		DbName:           os.Getenv("DB_NAME"),
		Port:             os.Getenv("PORT"),
	}
}

package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	ProviderToken string
	ProviderUrl   string
	MongoURI      string
	DbName        string
	Port          string
}

var AppConfig Config

func Load() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Aviso: .env não encontrado, usando variáveis do sistema")
	}

	AppConfig = Config{
		ProviderToken: os.Getenv("PROVIDER_TOKEN"),
		ProviderUrl:   os.Getenv("PROVIDER_URL"),
		MongoURI:      os.Getenv("MONGODB_URI"),
		DbName:        os.Getenv("DB_NAME"),
		Port:          os.Getenv("PORT"),
	}
}

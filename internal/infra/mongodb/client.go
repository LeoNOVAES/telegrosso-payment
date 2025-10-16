package mongodb

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoClient struct {
	Client *mongo.Client
	DB     *mongo.Database
}

func NewMongoClient(uri, dbName string) *MongoClient {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		log.Fatalf("Erro ao conectar no MongoDB: %v", err)
	}

	return &MongoClient{
		Client: client,
		DB:     client.Database(dbName),
	}
}

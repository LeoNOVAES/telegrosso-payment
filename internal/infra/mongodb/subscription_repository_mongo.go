package mongodb

import (
	"context"
	"time"

	"github.com/throindev/payments/internal/domain"
	"go.mongodb.org/mongo-driver/bson"
)

type SubscriptionMongoRepository struct {
	collectionName string
	client         *MongoClient
}

func NewSubscriptionMongoRepository(client *MongoClient) *SubscriptionMongoRepository {
	return &SubscriptionMongoRepository{
		collectionName: "subscriptions",
		client:         client,
	}
}

func (r *SubscriptionMongoRepository) Create(subscription *domain.Subscription) error {
	subscription.CreatedAt = time.Now().UTC()
	subscription.UpdatedAt = time.Now().UTC()

	_, err := r.client.DB.Collection(r.collectionName).InsertOne(context.Background(), subscription)
	return err
}

func (r *SubscriptionMongoRepository) Update(subscription *domain.Subscription) error {
	subscription.UpdatedAt = time.Now().UTC()

	_, err := r.client.DB.Collection(r.collectionName).UpdateOne(
		context.Background(),
		bson.M{"_id": subscription.ID},
		bson.M{"$set": subscription},
	)
	return err
}

func (r *SubscriptionMongoRepository) FindByID(id string) (*domain.Subscription, error) {
	var sub domain.Subscription
	err := r.client.DB.Collection(r.collectionName).FindOne(
		context.Background(),
		bson.M{"_id": id},
	).Decode(&sub)
	if err != nil {
		return nil, err
	}
	return &sub, nil
}

func (r *SubscriptionMongoRepository) FindByUserID(userID string) (*domain.Subscription, error) {
	var sub domain.Subscription
	err := r.client.DB.Collection(r.collectionName).FindOne(
		context.Background(),
		bson.M{"user_id": userID, "status": "active"},
	).Decode(&sub)
	if err != nil {
		return nil, err
	}
	return &sub, nil
}

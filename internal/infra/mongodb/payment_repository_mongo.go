package mongodb

import (
	"context"
	"time"

	"github.com/throindev/payments/internal/domain"
	"go.mongodb.org/mongo-driver/bson"
)

type PaymentMongoRepository struct {
	collectionName string
	client         *MongoClient
}

func NewPaymentMongoRepository(client *MongoClient) *PaymentMongoRepository {
	return &PaymentMongoRepository{
		collectionName: "payments",
		client:         client,
	}
}

func (r *PaymentMongoRepository) Create(payment *domain.Payment) error {
	payment.CreatedAt = time.Now().UTC()
	payment.UpdatedAt = time.Now().UTC()

	_, err := r.client.DB.Collection(r.collectionName).InsertOne(context.Background(), payment)
	return err
}

func (r *PaymentMongoRepository) Update(payment *domain.Payment) error {
	payment.UpdatedAt = time.Now().UTC()

	_, err := r.client.DB.Collection(r.collectionName).UpdateOne(
		context.Background(),
		bson.M{"_id": payment.ID},
		bson.M{"$set": payment},
	)
	return err
}

func (r *PaymentMongoRepository) UpdateByExternalId(payment *domain.Payment) error {
	payment.UpdatedAt = time.Now().UTC()
	update := bson.M{
		"status":        payment.Status,
		"updated_at":    time.Now().UTC(),
		"date_approved": payment.DateApproved,
	}
	_, err := r.client.DB.Collection(r.collectionName).UpdateOne(
		context.Background(),
		bson.M{"external_id": payment.ExternalId},
		bson.M{"$set": update},
	)
	return err
}

func (r *PaymentMongoRepository) FindByID(id string) (*domain.Payment, error) {
	var pay domain.Payment
	err := r.client.DB.Collection(r.collectionName).FindOne(
		context.Background(),
		bson.M{"_id": id},
	).Decode(&pay)
	if err != nil {
		return nil, err
	}
	return &pay, nil
}

func (r *PaymentMongoRepository) FindByExternalId(id string) (*domain.Payment, error) {
	var pay domain.Payment
	err := r.client.DB.Collection(r.collectionName).FindOne(
		context.Background(),
		bson.M{"external_id": id},
	).Decode(&pay)
	if err != nil {
		return nil, err
	}
	return &pay, nil
}

func (r *PaymentMongoRepository) FindByUserID(userID string) (*domain.Payment, error) {
	var pay domain.Payment
	err := r.client.DB.Collection(r.collectionName).FindOne(
		context.Background(),
		bson.M{"user_id": userID},
	).Decode(&pay)
	if err != nil {
		return nil, err
	}
	return &pay, nil
}

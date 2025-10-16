package mongodb

import (
	"context"
	"time"

	"github.com/throindev/payments/internal/domain"
	"go.mongodb.org/mongo-driver/bson"
)

type PlanMongoRepository struct {
	collectionName string
	client         *MongoClient
}

func NewPlanMongoRepository(client *MongoClient) *PlanMongoRepository {
	return &PlanMongoRepository{
		collectionName: "plans",
		client:         client,
	}
}

func (r *PlanMongoRepository) Create(Plan *domain.Plan) error {
	Plan.CreatedAt = time.Now().UTC()
	Plan.UpdatedAt = time.Now().UTC()

	_, err := r.client.DB.Collection(r.collectionName).InsertOne(context.Background(), Plan)
	return err
}

func (r *PlanMongoRepository) Update(plan *domain.Plan) error {
	plan.UpdatedAt = time.Now().UTC()

	_, err := r.client.DB.Collection(r.collectionName).UpdateOne(
		context.Background(),
		bson.M{"_id": plan.ID},
		bson.M{"$set": plan},
	)
	return err
}

func (r *PlanMongoRepository) FindByID(id string) (*domain.Plan, error) {
	var pay domain.Plan
	err := r.client.DB.Collection(r.collectionName).FindOne(
		context.Background(),
		bson.M{"_id": id},
	).Decode(&pay)
	if err != nil {
		return nil, err
	}
	return &pay, nil
}

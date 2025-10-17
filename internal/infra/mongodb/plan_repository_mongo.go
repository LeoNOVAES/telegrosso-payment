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

func (r *PlanMongoRepository) FindAll() ([]domain.Plan, error) {
	cursor, err := r.client.DB.Collection(r.collectionName).Find(context.Background(), bson.D{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())

	var results []domain.Plan
	for cursor.Next(context.Background()) {
		var item domain.Plan
		if err := cursor.Decode(&item); err != nil {
			return nil, err
		}
		results = append(results, item)
	}

	return results, cursor.Err()
}

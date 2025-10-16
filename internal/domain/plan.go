package domain

import (
	"time"

	"github.com/google/uuid"
)

type Plan struct {
	ID                string    `bson:"_id" json:"id"`
	Name              string    `bson:"name" json:"name"`
	Description       string    `bson:"description" json:"description"`
	Price             float64   `bson:"price" json:"price"`
	IntervalCountDays int       `bson:"interval_count" json:"interval_count"` // days
	GroupID           string    `bson:"group_id" json:"group_id"`
	Benefits          []string  `bson:"benefits" json:"benefits"`
	Active            bool      `bson:"active" json:"active"`
	CreatedAt         time.Time `bson:"created_at" json:"created_at"`
	UpdatedAt         time.Time `bson:"updated_at" json:"updated_at"`
}

func NewPlan(name, description string, price float64, intervalCount int, groupID string, benefits []string) *Plan {
	now := time.Now().UTC()

	return &Plan{
		ID:                uuid.NewString(),
		Name:              name,
		Description:       description,
		Price:             price,
		IntervalCountDays: intervalCount,
		GroupID:           groupID,
		Benefits:          benefits,
		Active:            true,
		CreatedAt:         now,
		UpdatedAt:         now,
	}
}

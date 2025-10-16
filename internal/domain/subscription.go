package domain

import (
	"time"

	"github.com/google/uuid"
)

type SubscriptionStatus string

const (
	SubscriptionStatusActive         SubscriptionStatus = "active"
	SubscriptionStatusPendingPayment SubscriptionStatus = "pending_payment"
	SubscriptionStatusExpired        SubscriptionStatus = "expired"
	SubscriptionStatusCanceled       SubscriptionStatus = "canceled"
)

type Subscription struct {
	ID              string             `bson:"_id" json:"id"`
	UserID          string             `bson:"user_id" json:"user_id"`
	PlanID          string             `bson:"plan_id" json:"plan_id"`
	PlanName        string             `bson:"plan_name" json:"plan_name"`
	Price           float64            `bson:"price" json:"price"`
	Currency        string             `bson:"currency" json:"currency"`
	PaymentProvider string             `bson:"payment_provider" json:"payment_provider"`
	NextDueAt       time.Time          `bson:"next_due_at" json:"next_due_at"`
	PaidAt          time.Time          `bson:"paid_at" json:"paid_at"`
	LastPayment     Payment            `bson:"last_payment" json:"last_payment"`
	Status          SubscriptionStatus `bson:"status" json:"status"`
	Metadata        map[string]any     `bson:"metadata" json:"metadata"`
	CreatedAt       time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt       time.Time          `bson:"updated_at" json:"updated_at"`
}

func NewSubscription(
	userId string,
	planId string,
	planName string,
	price float64,
	currency string,
	paymentProvider string,
	nextDue time.Time,
	paidAt time.Time,
	lastPayment Payment,
	status SubscriptionStatus,
	metadata map[string]any,
) Subscription {
	return Subscription{
		ID:              uuid.NewString(),
		UserID:          userId,
		PlanID:          planId,
		PlanName:        planName,
		Price:           price,
		Currency:        currency,
		PaymentProvider: paymentProvider,
		NextDueAt:       nextDue,
		PaidAt:          paidAt,
		LastPayment:     lastPayment,
		Status:          status,
		Metadata:        metadata,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}
}

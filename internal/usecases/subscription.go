package usecases

import (
	"errors"
	"fmt"
	"time"

	"github.com/throindev/payments/internal/domain"
)

type SubscriptionUsecases struct {
	repository   domain.SubscriptionRepository
	planUsecases PlanUsecases
}

func NewSubscriptionUsecases(repository domain.SubscriptionRepository, planUsecases PlanUsecases) SubscriptionUsecases {
	return SubscriptionUsecases{repository, planUsecases}
}

func (s *SubscriptionUsecases) CreateSubscription(
	userId string,
	price float64,
	currency string,
	paymentProvider string,
	paidAt time.Time,
	lastPayment domain.Payment,
	metadata map[string]any,
	planId string,
) (domain.Subscription, error) {
	plan, err := s.planUsecases.FindById(planId)

	if err != nil {
		return domain.Subscription{}, err
	}

	if plan == nil {
		return domain.Subscription{}, errors.New("plan not found")
	}

	nextDue := time.Now().Add(time.Duration(plan.IntervalCountDays) * 24 * time.Hour)

	subscription := domain.NewSubscription(
		userId,
		plan.ID,
		plan.Name,
		price,
		currency,
		paymentProvider,
		nextDue,
		paidAt,
		lastPayment,
		"active",
		metadata,
	)

	err = s.repository.Create(&subscription)

	if err != nil {
		fmt.Printf("error to create subscription %v", err)
		return domain.Subscription{}, err
	}

	return subscription, nil
}

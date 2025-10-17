package usecases

import (
	"fmt"

	"github.com/throindev/payments/internal/domain"
)

type PlanUsecases struct {
	repository domain.PlanRepository
}

func NewPlanUsecases(repository domain.PlanRepository) PlanUsecases {
	return PlanUsecases{repository}
}

func (s *PlanUsecases) CreatePlan(
	name,
	description string,
	price float64,
	intervalCount int,
	groupID string,
	benefits []string,
) (*domain.Plan, error) {
	plan := domain.NewPlan(
		name,
		description,
		price,
		intervalCount,
		groupID,
		benefits,
	)

	err := s.repository.Create(plan)

	if err != nil {
		fmt.Printf("error to create Plan %v", err)
		return &domain.Plan{}, err
	}

	return plan, nil
}

func (s *PlanUsecases) FindAll() ([]domain.Plan, error) {
	plans, err := s.repository.FindAll()

	if err != nil {
		fmt.Printf("error to GET Plans %v", err)
		return []domain.Plan{}, err
	}

	return plans, nil
}

func (s *PlanUsecases) FindById(id string) (*domain.Plan, error) {
	plan, err := s.repository.FindByID(id)

	if err != nil {
		fmt.Printf("error to GET Plan %v", err)
		return &domain.Plan{}, err
	}

	return plan, nil
}

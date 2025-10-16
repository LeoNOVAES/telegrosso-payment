package domain

type PlanRepository interface {
	Create(plan *Plan) error
	Update(plan *Plan) error
	FindByID(id string) (*Plan, error)
}

package domain

type PaymentRepository interface {
	Create(payment *Payment) error
	Update(payment *Payment) error
	UpdateByExternalId(payment *Payment) error
	FindByID(id string) (*Payment, error)
	FindByUserID(userID string) (*Payment, error)
}

package interfaces

import "github.com/throindev/payments/internal/domain"

type PaymentProvider interface {
	CreatePayment(amount float64, description string, method string, userId, planId string) (domain.Payment, error)
	GetPayment(id string) (domain.Payment, error)
}

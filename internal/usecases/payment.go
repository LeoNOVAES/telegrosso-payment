package usecases

import (
	"errors"
	"fmt"
	"strings"

	"github.com/throindev/payments/internal/domain"
	"github.com/throindev/payments/internal/domain/interfaces"
)

type PaymentUsecases struct {
	repository          domain.PaymentRepository
	provider            interfaces.PaymentProvider
	subscriptionUsecase SubscriptionUsecases
	planUsecase         PlanUsecases
}

func NewPaymentUsecases(repository domain.PaymentRepository, provider interfaces.PaymentProvider, subscriptionUsecase SubscriptionUsecases, planUsecase PlanUsecases) PaymentUsecases {
	return PaymentUsecases{repository, provider, subscriptionUsecase, planUsecase}
}

func (s *PaymentUsecases) ConfirmPayment(external_id string, provider string) (domain.Payment, error) {
	payment, err := s.GetPaymentFromProvider(external_id)

	if err != nil {
		fmt.Printf("Error to get payment %s -> %v", external_id, err)
		return domain.Payment{}, err
	}

	if payment.Status != "approved" {
		return domain.Payment{}, errors.New("pagamento nao aprovado")
	}

	s.repository.UpdateByExternalId(&payment)

	parts := strings.Split(payment.ExternalReference, ":")

	planId := parts[1]
	userId := parts[2]

	if planId == "" || userId == "" {
		return domain.Payment{}, errors.New("usuario ou plano invalidos")
	}

	plan, err := s.planUsecase.FindById(planId)

	if err != nil || plan == nil {
		fmt.Printf("Plano nao existe")
		return domain.Payment{}, errors.New("plano nao existe")
	}

	_, errSubscription := s.subscriptionUsecase.CreateSubscription(
		userId,
		payment.TransactionAmount,
		"BRL",
		provider,
		payment.DateApproved,
		payment,
		map[string]any{},
		plan.ID,
	)

	if errSubscription != nil {
		fmt.Printf("Error to create subscription %s -> %v", external_id, err)
		return domain.Payment{}, err
	}

	return payment, nil
}

func (s *PaymentUsecases) CreatePayment(typePayment string, userId string, planId string) (domain.Payment, error) {
	plan, err := s.planUsecase.FindById(planId)

	if err != nil || plan == nil {
		fmt.Printf("Plano nao existe")
		return domain.Payment{}, errors.New("plano nao existe")
	}

	description := fmt.Sprintf("Plan paid for %s to plan %s", userId, plan.Name)

	paymentFromProvider, err := s.provider.CreatePayment(plan.Price, description, typePayment, userId, planId)
	paymentFromProvider.PlanId = planId

	if err != nil {
		fmt.Printf("error to create payment in provider %v", err)
		return domain.Payment{}, err
	}

	if err := s.repository.Create(&paymentFromProvider); err != nil {
		return domain.Payment{}, err
	}

	return paymentFromProvider, nil
}

func (s *PaymentUsecases) GetPaymentFromProvider(id string) (domain.Payment, error) {
	payment, err := s.provider.GetPayment(id)

	if err != nil {
		return domain.Payment{}, err
	}

	return payment, nil
}

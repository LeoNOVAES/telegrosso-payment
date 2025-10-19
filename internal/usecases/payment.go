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
	eventbusRepository  interfaces.EventBus
}

func NewPaymentUsecases(repository domain.PaymentRepository, provider interfaces.PaymentProvider, subscriptionUsecase SubscriptionUsecases, planUsecase PlanUsecases, eventbus interfaces.EventBus) PaymentUsecases {
	return PaymentUsecases{repository, provider, subscriptionUsecase, planUsecase, eventbus}
}

func (s *PaymentUsecases) ConfirmPayment(external_id string, provider string) (domain.Payment, error) {
	paymentFromProvider, err := s.GetPaymentFromProvider(external_id)

	if err != nil {
		fmt.Printf("Error to get payment %s -> %v", external_id, err)
		return domain.Payment{}, err
	}

	if paymentFromProvider.Status != "approved" {
		return domain.Payment{}, errors.New("pagamento nao aprovado")
	}

	errUpdate := s.repository.UpdateByExternalId(&paymentFromProvider)
	payment, _ := s.repository.FindByExternalId(external_id)
	fmt.Printf("AAAAAAAAAAAAAAAAA %v", payment)

	if errUpdate != nil {
		return domain.Payment{}, err
	}

	parts := strings.Split(paymentFromProvider.ExternalReference, ":")
	userId := parts[0]
	planId := parts[1]

	if planId == "" || userId == "" {
		return domain.Payment{}, errors.New("usuario ou plano invalidos")
	}

	plan, err := s.planUsecase.FindById(planId)

	if err != nil || plan == nil {
		fmt.Printf("Plano nao existe")
		return domain.Payment{}, errors.New("plano nao existe")
	}

	subscription, errSubscription := s.subscriptionUsecase.CreateSubscription(
		userId,
		payment.TransactionAmount,
		"BRL",
		provider,
		payment.DateApproved,
		*payment,
		map[string]any{},
		plan.ID,
	)

	if errSubscription != nil {
		fmt.Printf("Error to create subscription %s -> %v", external_id, err)
		return domain.Payment{}, err
	}

	event := domain.NewEvent(domain.EventPaymentConfirmed, subscription)

	s.eventbusRepository.Publish(event.EventType, event)
	return *payment, nil
}

func (s *PaymentUsecases) CreatePayment(typePayment string, userId string, planId string, chatId string) (domain.Payment, error) {
	plan, err := s.planUsecase.FindById(planId)

	if err != nil || plan == nil {
		fmt.Printf("Plano nao existe")
		return domain.Payment{}, errors.New("plano nao existe")
	}

	description := fmt.Sprintf("Plan paid for %s to plan %s", userId, plan.Name)

	paymentFromProvider, err := s.provider.CreatePayment(plan.Price, description, typePayment, userId, planId)
	paymentFromProvider.PlanId = planId
	paymentFromProvider.ChatId = chatId

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

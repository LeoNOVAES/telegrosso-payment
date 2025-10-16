package domain

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

type Payer struct {
	Email string `json:"email"`
}

type Payment struct {
	ID                string    `bson:"_id" json:"id"`
	ExternalId        string    `bson:"external_id" json:"external_id"`
	TransactionAmount float64   `bson:"transaction_amount" json:"transaction_amount"`
	Description       string    `bson:"description" json:"description"`
	PaymentMethodID   string    `bson:"payment_method_id" json:"payment_method_id"`
	Payer             Payer     `bson:"payer" json:"payer"`
	PlanId            string    `bson:"plan_id" json:"plan_id"`
	ExternalReference string    `bson:"external_reference" json:"external_reference"`
	QRCode            string    `bson:"qr_code" json:"qr_code"`
	QRCodeBase64      string    `bson:"qr_code_base64" json:"qr_code_base64"`
	TicketURL         string    `bson:"ticket_url" json:"ticket_url"`
	Status            string    `bson:"status" json:"status"`
	DateApproved      time.Time `bson:"date_approved" json:"date_approved"`
	CreatedAt         time.Time `bson:"created_at" json:"created_at"`
	UpdatedAt         time.Time `bson:"updated_at" json:"updated_at"`
}

func NewPayment(
	external_id string,
	amount float64,
	description string,
	paymentMethod string,
	external_code string,
	qrCode string,
	qrCodeBase64 string,
	ticketUrl string,
	dateApproved time.Time,
	status string,
	planId string,
) Payment {
	payer := fmt.Sprintf("%s@noreply.com", external_code)

	return Payment{
		ID:                uuid.NewString(),
		ExternalId:        external_id,
		TransactionAmount: amount,
		Description:       description,
		PaymentMethodID:   paymentMethod,
		Payer:             Payer{Email: payer},
		ExternalReference: external_code,
		QRCode:            qrCode,
		QRCodeBase64:      qrCodeBase64,
		DateApproved:      dateApproved,
		TicketURL:         ticketUrl,
		Status:            status,
		PlanId:            planId,
	}
}

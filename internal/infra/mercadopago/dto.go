package mercadopago

import "time"

type PaymentResponse struct {
	ID                int64     `json:"id"`
	TransactionAmount float64   `json:"transaction_amount"`
	Description       string    `json:"description"`
	ExternalReference string    `json:"external_reference"`
	Status            string    `json:"status"`
	DateApproved      time.Time `json:"date_approved"`
	paymentMethod     struct {
		ID       string `json:"id"`
		Type     string `json:"type"`
		IssuerID string `json:"issuer_id"`
	}

	Payer struct {
		Email string `json:"email"`
	} `json:"payer"`

	PointOfInteraction struct {
		TransactionData struct {
			QRCode       string `json:"qr_code"`
			QRCodeBase64 string `json:"qr_code_base64"`
			TicketURL    string `json:"ticket_url"`
		} `json:"transaction_data"`
	} `json:"point_of_interaction"`
}

type PaymentCallback struct {
	Action     string `json:"action"`
	APIVersion string `json:"api_version"`
	Data       struct {
		ID string `json:"id"`
	} `json:"data"`
	DateCreated time.Time `json:"date_created"`
	ID          int64     `json:"id"`
	LiveMode    bool      `json:"live_mode"`
	Type        string    `json:"type"`
	UserID      string    `json:"user_id"`
}

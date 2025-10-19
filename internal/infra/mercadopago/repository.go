package mercadopago

import (
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/throindev/payments/cmd/config"
	"github.com/throindev/payments/internal/domain"
	"github.com/throindev/payments/internal/infra"
)

type MercadoPagoRepository struct {
	client infra.HTTPClient
}

func NewClient() *MercadoPagoRepository {
	headers := map[string]string{
		"Authorization": fmt.Sprintf("Bearer %s", config.AppConfig.ProviderToken),
	}

	client := infra.NewHTTPClient(headers, 30*time.Second)
	return &MercadoPagoRepository{*client}
}

func (c *MercadoPagoRepository) CreatePayment(amount float64, description string, method string, userId, planId string) (domain.Payment, error) {
	fmt.Println("Chamando Mercado Pago para criar Payment de:", amount)
	var response PaymentResponse
	external_code := fmt.Sprintf("%s:%s", userId, planId)

	customHeader := map[string]string{
		"X-Idempotency-Key": uuid.New().String(),
	}

	payload := map[string]interface{}{
		"transaction_amount": amount,
		"description":        description,
		"payment_method_id":  method,
		"payer": map[string]interface{}{
			"email": fmt.Sprintf("%s_%s@noreply.com", userId, planId),
		},
		"external_reference": external_code,
	}

	err := c.client.Post(config.AppConfig.ProviderUrl, payload, &response, customHeader)

	if err != nil {
		log.Printf("Erro ao criar Payment para PaymentID %s: %v", external_code, err)
		return domain.Payment{}, err
	}

	return c.parseToPayment(response), nil
}

func (c *MercadoPagoRepository) GetPayment(id string) (domain.Payment, error) {
	fmt.Println("Chamando Mercado Pago para buscar um Payment:", id)
	var response PaymentResponse

	err := c.client.Get(fmt.Sprintf("%s/%s", config.AppConfig.ProviderUrl, id), &response)

	if err != nil {
		log.Printf("Erro ao pegar pagamento %v", err)
		return c.parseToPayment(response), err
	}

	return c.parseToPayment(response), nil
}

func (c *MercadoPagoRepository) parseToPayment(response PaymentResponse) domain.Payment {
	return domain.NewPayment(
		fmt.Sprintf("%d", response.ID),
		response.TransactionAmount,
		response.Description,
		response.paymentMethod.ID,
		response.ExternalReference,
		response.PointOfInteraction.TransactionData.QRCode,
		response.PointOfInteraction.TransactionData.QRCodeBase64,
		response.PointOfInteraction.TransactionData.TicketURL,
		response.DateApproved,
		response.Status,
		"",
		"",
	)
}

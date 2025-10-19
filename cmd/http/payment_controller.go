package http

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/throindev/payments/internal/infra/mercadopago"
	"github.com/throindev/payments/internal/usecases"
)

type PaymentController struct {
	usecase usecases.PaymentUsecases
}

func NewPaymentController(usecase usecases.PaymentUsecases) *PaymentController {
	return &PaymentController{usecase: usecase}
}

func (h *PaymentController) CreatePayment(c *gin.Context) {
	body, err := io.ReadAll(c.Request.Body)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid body"})
		return
	}

	var payload struct {
		Method string `json:"method"`
		UserId string `json:"user_id"`
		PlanId string `json:"plan_id"`
		ChatId string `json:"chat_id"`
	}

	if err := json.Unmarshal(body, &payload); err != nil {
		fmt.Println("Erro ao parsear JSON:", err)
		return
	}

	payment, err := h.usecase.CreatePayment(payload.Method, payload.UserId, payload.PlanId, payload.ChatId)

	if err != nil {
		fmt.Sprintf("error to create payment %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "ERROR TO CREATE"})
		return
	}

	c.JSON(http.StatusOK, payment)
}

func (h *PaymentController) CallbackfromMercadoPago(c *gin.Context) {
	body, err := io.ReadAll(c.Request.Body)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid body"})
		return
	}

	var payload mercadopago.PaymentCallback

	if err := json.Unmarshal(body, &payload); err != nil {
		fmt.Println("Erro ao parsear JSON:", err)
		return
	}

	if payload.Action != "payment.updated" {
		return
	}

	payment, errSubs := h.usecase.ConfirmPayment(payload.Data.ID, "mercadopago")

	if errSubs != nil {
		fmt.Sprintf("error to confirm payment %v", errSubs)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "ERROR TO CONFIRM"})
		return
	}

	c.JSON(http.StatusOK, payment)
}

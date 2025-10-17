package http

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/throindev/payments/internal/usecases"
)

type PlanController struct {
	usecase usecases.PlanUsecases
}

func NewPlanController(usecase usecases.PlanUsecases) *PlanController {
	return &PlanController{usecase: usecase}
}

func (h *PlanController) CreatePlan(c *gin.Context) {
	body, err := io.ReadAll(c.Request.Body)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid body"})
		return
	}

	var payload struct {
		Name              string   `json:"name"`
		Description       string   `json:"description"`
		Price             float64  `json:"price"`
		IntervalCountDays int      `json:"interval_count"`
		GroupID           string   `json:"group_id"`
		Benefits          []string `json:"benefits"`
	}

	if err := json.Unmarshal(body, &payload); err != nil {
		fmt.Println("Erro ao parsear JSON:", err)
		return
	}

	plan, err := h.usecase.CreatePlan(payload.Name, payload.Description, payload.Price, payload.IntervalCountDays, payload.GroupID, payload.Benefits)

	if err != nil {
		fmt.Sprintf("error to create Plan %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "ERROR TO CREATE"})
		return
	}

	c.JSON(http.StatusOK, plan)
}

func (h *PlanController) GetPlans(c *gin.Context) {
	results, err := h.usecase.FindAll()

	if err != nil {
		fmt.Sprintf("error to GET Plans %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "ERROR TO GTE PLANS"})
		return
	}

	c.JSON(http.StatusOK, results)
}

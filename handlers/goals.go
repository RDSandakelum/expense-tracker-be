package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type GoalResponse struct {
	ID            uuid.UUID `json:"id"`
	Name          string    `json:"name"`
	Description   string    `json:"description"`
	TargetAmount  float64   `json:"target_amount"`
	CurrentAmount float64   `json:"current_amount"`
	Saved         float64   `json:"saved"`
	Deadline      string    `json:"deadline"`
	Category      string    `json:"category"`
	Status        string    `json:"status"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

func GetGoals(c *gin.Context) {
	// Get user ID from context (set by AuthMiddleware)
	_, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "User ID not found in context"})
		return
	}

	// Dummy goals data
	now := time.Now()
	goals := []GoalResponse{
		{
			ID:            uuid.New(),
			Name:          "Vacation Fund",
			Description:   "Summer vacation to Europe",
			TargetAmount:  5000.00,
			CurrentAmount: 2350.75,
			Saved:         2350.75,
			Deadline:      now.AddDate(0, 4, 0).Format("2006-01-02"),
			Category:      "Travel",
			Status:        "in_progress",
			CreatedAt:     now.AddDate(0, -3, 0),
			UpdatedAt:     now,
		},
		{
			ID:            uuid.New(),
			Name:          "Emergency Fund",
			Description:   "3-6 months of living expenses",
			TargetAmount:  20000.00,
			CurrentAmount: 12500.00,
			Saved:         12500.00,
			Deadline:      now.AddDate(1, 0, 0).Format("2006-01-02"),
			Category:      "Security",
			Status:        "in_progress",
			CreatedAt:     now.AddDate(-1, 0, 0),
			UpdatedAt:     now,
		},
		{
			ID:            uuid.New(),
			Name:          "New Car",
			Description:   "Down payment for new vehicle",
			TargetAmount:  10000.00,
			CurrentAmount: 4750.00,
			Saved:         4750.00,
			Deadline:      now.AddDate(1, 6, 0).Format("2006-01-02"),
			Category:      "Transportation",
			Status:        "in_progress",
			CreatedAt:     now.AddDate(-6, 0, 0),
			UpdatedAt:     now,
		},
		{
			ID:            uuid.New(),
			Name:          "Home Renovation",
			Description:   "Kitchen and bathroom updates",
			TargetAmount:  25000.00,
			CurrentAmount: 8200.50,
			Saved:         8200.50,
			Deadline:      now.AddDate(2, 0, 0).Format("2006-01-02"),
			Category:      "Home",
			Status:        "in_progress",
			CreatedAt:     now.AddDate(-12, 0, 0),
			UpdatedAt:     now,
		},
	}

	c.JSON(http.StatusOK, gin.H{
		"goals": goals,
	})
}

type AddFundsToGoalRequest struct {
	GoalID    uuid.UUID `json:"goal_id" binding:"required"`
	AccountID uuid.UUID `json:"account_id" binding:"required"`
	Amount    float64   `json:"amount" binding:"required"`
}

func AddFundsToGoal(c *gin.Context) {
	// Get user ID from context (set by AuthMiddleware)
	_, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "User ID not found in context"})
		return
	}

	var req AddFundsToGoalRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid request body"})
		return
	}

	if req.Amount <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Amount must be greater than 0"})
		return
	}

	// Dummy response - simulating funds added to goal
	updatedGoal := GoalResponse{
		ID:            req.GoalID,
		Name:          "Vacation Fund",
		Description:   "Summer vacation to Europe",
		TargetAmount:  5000.00,
		CurrentAmount: 2350.75 + req.Amount,
		Saved:         2350.75 + req.Amount,
		Deadline:      time.Now().AddDate(0, 4, 0).Format("2006-01-02"),
		Category:      "Travel",
		Status:        "in_progress",
		CreatedAt:     time.Now().AddDate(0, -3, 0),
		UpdatedAt:     time.Now(),
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Funds added successfully",
		"goal":    updatedGoal,
	})
}

package handlers

import (
	"expense-tracker-be/storage"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type GoalResponse struct {
	ID           uuid.UUID `json:"id"`
	Name         string    `json:"name"`
	TargetAmount float64   `json:"target_amount"`
	Saved        float64   `json:"saved"`
	Deadline     string    `json:"deadline"`
	Completed    bool      `json:"status"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

func GetGoals(c *gin.Context) {
	// Get user ID from context (set by AuthMiddleware)
	_, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "User ID not found in context"})
		return
	}

	userIDInterface, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found in context"})
		return
	}

	userIDStr, ok := userIDInterface.(string)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID format"})
		return
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	goals, err := storage.GetGoalsByUserID(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "internal server error",
		})
		return
	}

	goalResponse := []GoalResponse{}
	for _, goal := range goals {
		goalResp := GoalResponse{
			ID:           goal.ID,
			Name:         goal.Name,
			TargetAmount: goal.TargetAmount,
			Saved:        goal.SavedAmount,
			Deadline:     goal.TargetDate.Format("2026-01-02"),
			Completed:    goal.IsCompleted,
		}
		goalResponse = append(goalResponse, goalResp)
	}

	c.JSON(http.StatusOK, gin.H{
		"goals": goalResponse,
	})
}

package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type TransactionDetailResponse struct {
	ID            uuid.UUID `json:"id"`
	UserID        uuid.UUID `json:"user_id"`
	AccountID     uuid.UUID `json:"account_id"`
	CategoryID    uuid.UUID `json:"category_id"`
	SubcategoryID uuid.UUID `json:"subcategory_id"`
	Category      string    `json:"category"`
	SubCategory   string    `json:"sub_category"`
	Amount        float64   `json:"amount"`
	Type          string    `json:"type"`
	Description   string    `json:"description"`
	Date          string    `json:"date"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

func GetTransaction(c *gin.Context) {
	// Get user ID from context (set by AuthMiddleware)
	_, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "User ID not found in context"})
		return
	}

	// Dummy single transaction
	transaction := TransactionDetailResponse{
		ID:            uuid.New(),
		UserID:        uuid.New(),
		AccountID:     uuid.New(),
		CategoryID:    uuid.New(),
		SubcategoryID: uuid.New(),
		Category:      "Food & Dining",
		SubCategory:   "Groceries",
		Amount:        125.50,
		Type:          "DEBIT",
		Description:   "Walmart - Weekly groceries",
		Date:          time.Now().AddDate(0, 0, -5).Format("2006-01-02"),
		CreatedAt:     time.Now().AddDate(0, 0, -5),
		UpdatedAt:     time.Now().AddDate(0, 0, -5),
	}

	c.JSON(http.StatusOK, transaction)
}

type UpdateTransactionRequest struct {
	Amount      float64 `json:"amount"`
	Type        string  `json:"type"`
	Description string  `json:"description"`
	Date        string  `json:"date"`
	CategoryID  uuid.UUID `json:"category_id"`
	SubcategoryID uuid.UUID `json:"subcategory_id"`
}

func UpdateTransaction(c *gin.Context) {
	// Get user ID from context (set by AuthMiddleware)
	_, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "User ID not found in context"})
		return
	}

	var req UpdateTransactionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid request body"})
		return
	}

	// Dummy updated transaction
	updated := TransactionDetailResponse{
		ID:            uuid.New(),
		UserID:        uuid.New(),
		AccountID:     uuid.New(),
		CategoryID:    req.CategoryID,
		SubcategoryID: req.SubcategoryID,
		Category:      "Food & Dining",
		SubCategory:   "Groceries",
		Amount:        req.Amount,
		Type:          req.Type,
		Description:   req.Description,
		Date:          req.Date,
		CreatedAt:     time.Now().AddDate(0, 0, -5),
		UpdatedAt:     time.Now(),
	}

	c.JSON(http.StatusOK, gin.H{
		"message":     "Transaction updated successfully",
		"transaction": updated,
	})
}

func DeleteTransaction(c *gin.Context) {
	// Get user ID from context (set by AuthMiddleware)
	_, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "User ID not found in context"})
		return
	}

	transactionID := c.Param("transactionId")

	c.JSON(http.StatusOK, gin.H{
		"message": "Transaction deleted successfully",
		"id":      transactionID,
	})
}

type CreateGoalRequest struct {
	Name          string    `json:"name" binding:"required"`
	Description   string    `json:"description"`
	TargetAmount  float64   `json:"target_amount" binding:"required"`
	Deadline      string    `json:"deadline"`
	Category      string    `json:"category"`
}

func CreateGoal(c *gin.Context) {
	// Get user ID from context (set by AuthMiddleware)
	_, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "User ID not found in context"})
		return
	}

	var req CreateGoalRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid request body"})
		return
	}

	// Dummy created goal
	goal := GoalResponse{
		ID:            uuid.New(),
		Name:          req.Name,
		Description:   req.Description,
		TargetAmount:  req.TargetAmount,
		CurrentAmount: 0.00,
		Saved:         0.00,
		Deadline:      req.Deadline,
		Category:      req.Category,
		Status:        "in_progress",
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Goal created successfully",
		"goal":    goal,
	})
}

type UpdateGoalRequest struct {
	Name          string  `json:"name"`
	Description   string  `json:"description"`
	TargetAmount  float64 `json:"target_amount"`
	CurrentAmount float64 `json:"current_amount"`
	Deadline      string  `json:"deadline"`
	Status        string  `json:"status"`
}

func UpdateGoal(c *gin.Context) {
	// Get user ID from context (set by AuthMiddleware)
	_, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "User ID not found in context"})
		return
	}

	var req UpdateGoalRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid request body"})
		return
	}

	// Dummy updated goal
	goal := GoalResponse{
		ID:            uuid.New(),
		Name:          req.Name,
		Description:   req.Description,
		TargetAmount:  req.TargetAmount,
		CurrentAmount: req.CurrentAmount,
		Saved:         req.CurrentAmount,
		Deadline:      req.Deadline,
		Category:      "Custom",
		Status:        req.Status,
		CreatedAt:     time.Now().AddDate(0, -3, 0),
		UpdatedAt:     time.Now(),
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Goal updated successfully",
		"goal":    goal,
	})
}

func DeleteGoal(c *gin.Context) {
	// Get user ID from context (set by AuthMiddleware)
	_, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "User ID not found in context"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Goal deleted successfully",
	})
}

package handlers

import (
	"bytes"
	"expense-tracker-be/dto"
	"expense-tracker-be/storage"
	"io"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func CreateTransaction(c *gin.Context) {
	// Get user ID from context (set by AuthMiddleware)
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
	bodyBytes, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Unable to read request body"})
		return
	}

	// 2. LOG THE RAW JSON STRING HERE
	log.Printf("📥 Raw Incoming Request Body: %s", string(bodyBytes))

	// 3. CRITICAL: Restore the body stream so c.ShouldBindJSON can read it next
	c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
	// Parse request body
	var req dto.CreateTransactionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}
	// Create transaction in database
	transaction, err := storage.CreateTransactionRecord(
		userID,
		req.AccountID,
		&req.SubcategoryID,
		req.Type,
		req.Amount,
		req.Note,
		req.Date,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create transaction"})
		return
	}

	err = storage.AddToBudget(userID, req.SubcategoryID, req.Date, req.Amount)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update budget"})
		return
	}
	// Return created transaction
	response := dto.TransactionResponse{
		ID:            transaction.ID,
		UserID:        transaction.UserID,
		UserAccountID: transaction.AccountID,
		Type:          transaction.Type,
		Amount:        transaction.Amount,
		SubcategoryID: *transaction.SubCategoryID,
		CreatedAt:     transaction.TransactionDate.String(),
	}

	c.JSON(http.StatusCreated, response)
}

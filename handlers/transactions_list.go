package handlers

import (
	"expense-tracker-be/storage"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type TransactionListResponse struct {
	ID          uuid.UUID `json:"id"`
	Date        string    `json:"date"`
	Category    string    `json:"category"`
	SubCategory string    `json:"sub_category"`
	Description string    `json:"description"`
	Amount      float64   `json:"amount"`
	Type        string    `json:"type"`
	Account     string    `json:"account"`
	CreatedAt   string    `json:"created_at"`
	UpdatedAt   string    `json:"updated_at"`
}

type TransactionsListWrapper struct {
	Transactions []TransactionListResponse `json:"transactions"`
	Total        int                       `json:"total"`
	Limit        int                       `json:"limit"`
	Offset       int                       `json:"offset"`
	HasMore      bool                      `json:"has_more"`
}

func GetCategoryTransactionsList(c *gin.Context) {
	// 1. Extract user credentials from context
	userIDInterface, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "User ID not found in context"})
		return
	}

	userIDStr, ok := userIDInterface.(string)
	userID, err := uuid.Parse(userIDStr)
	if !ok || err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user identification"})
		return
	}

	// 2. Parse pagination and date context parameters
	now := time.Now()
	monthStr := c.DefaultQuery("month", strconv.Itoa(int(now.Month())))
	yearStr := c.DefaultQuery("year", strconv.Itoa(now.Year()))
	limitStr := c.DefaultQuery("limit", "50")
	offsetStr := c.DefaultQuery("offset", "0")

	month, err := strconv.Atoi(monthStr)
	if err != nil || month < 1 || month > 12 {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid month parameter"})
		return
	}

	year, err := strconv.Atoi(yearStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid year parameter"})
		return
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 1 || limit > 100 {
		limit = 50
	}

	offset, err := strconv.Atoi(offsetStr)
	if err != nil || offset < 0 {
		offset = 0
	}

	// 3. Compute exact date boundaries for the targeted month
	startOfMonth := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.Local)
	// Add exactly 1 month to the start date and subtract 1 nanosecond to hit the exact end bound
	endOfMonth := startOfMonth.AddDate(0, 1, 0).Add(-time.Nanosecond)

	// 4. Fetch the real records from database
	dbTransactions, totalCount, err := storage.GetTransactionsByMonth(userID, startOfMonth, endOfMonth, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch database transactions"})
		return
	}

	// 5. Map internal structural records into the output DTO array safely
	responseList := []TransactionListResponse{}
	for _, tx := range dbTransactions {
		// Enforce fallbacks for missing nested relational references
		accountName := "Unknown Account"
		if tx.Account != nil {
			accountName = tx.Account.Name
		}

		categoryName := "Uncategorized"
		subCategoryName := "General"
		if tx.SubCategory != nil {
			subCategoryName = tx.SubCategory.Name
			if tx.SubCategory.Category != nil {
				categoryName = tx.SubCategory.Category.Name
			}
		}

		responseList = append(responseList, TransactionListResponse{
			ID:          tx.ID,
			Date:        tx.TransactionDate.Format("2006-01-02"),
			Category:    categoryName,
			SubCategory: subCategoryName,
			Description: tx.Note,
			Amount:      tx.Amount,
			Type:        strings.ToUpper(tx.Type), // Normalizes 'Debit'/'Credit' strings to 'DEBIT'/'CREDIT' matching frontend types
			Account:     accountName,
			CreatedAt:   tx.TransactionDate.Format(time.RFC3339),
			UpdatedAt:   tx.TransactionDate.Format(time.RFC3339), // Fallback map if explicit field isn't declared
		})
	}

	// 6. Calculate progressive pagination states
	hasMore := int64(offset+limit) < totalCount

	c.JSON(http.StatusOK, gin.H{
		"transactions": responseList,
		"total":        totalCount,
		"limit":        limit,
		"offset":       offset,
		"hasMore":      hasMore,
	})
}

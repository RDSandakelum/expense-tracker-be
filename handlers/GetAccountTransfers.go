package handlers

import (
	"expense-tracker-be/storage"
	"net/http"
	"sort"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type TransferResponse struct {
	Type AccountTransactionType `json:"type"`

	// Account-to-Account Transfer
	FromAccountName string `json:"from_account_name,omitempty"`
	ToAccountName   string `json:"to_account_name,omitempty"`

	// Savings Transfer
	AccountName string `json:"account_name,omitempty"`
	GoalName    string `json:"goal_name,omitempty"`
	Direction   string `json:"direction,omitempty"`

	// Common
	Amount float64 `json:"amount"`
	Date   string  `json:"date"` // Formatted as YYYY-MM-DD
}

type AccountTransactionType string

const (
	AccTransactionTypeAccount AccountTransactionType = "ACCOUNT"
	AccTransactionTypeSavings AccountTransactionType = "SAVINGS"
)

func GetAccountTransfers(c *gin.Context) {
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

	accountTransfers, err := storage.GetAllAccountTransfers(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch account transfers"})
		return
	}

	result := make([]TransferResponse, 0)

	const dateLayout = "2006-01-02"

	for _, at := range accountTransfers {
		fromAccount, err := storage.GetAccountByID(at.FromAccountID)
		if err != nil {
			continue
		}

		toAccount, err := storage.GetAccountByID(at.ToAccountID)
		if err != nil {
			continue
		}
		result = append(result, TransferResponse{
			Type:            AccTransactionTypeAccount,
			FromAccountName: fromAccount.Name,
			ToAccountName:   toAccount.Name,
			Amount:          at.Amount,
			Date:            at.TransferDate.Format(dateLayout),
		})
	}

	saveWithdrawals, err := storage.GetAllSavingWithdrawals(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch savings withdrawals"})
		return
	}

	for _, sw := range saveWithdrawals {
		account, err := storage.GetAccountByID(sw.AccountID)
		if err != nil {
			continue
		}

		goal, err := storage.GetGoalByID(sw.GoalID)
		if err != nil {
			continue
		}

		result = append(result, TransferResponse{
			Type:        AccTransactionTypeSavings,
			AccountName: account.Name,
			GoalName:    goal.Name,
			Direction:   sw.Direction,
			Amount:      sw.Amount,
			Date:        sw.WithdrawnAt.Format(dateLayout),
		})
	}

	sort.Slice(result, func(i, j int) bool {
		return result[i].Date > result[j].Date
	})

	c.JSON(http.StatusOK, result)
}

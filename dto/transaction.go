package dto

import (
	"time"

	"github.com/google/uuid"
)

type CategoryResponse struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}

type SubcategoryResponse struct {
	ID         uuid.UUID `json:"id"`
	CategoryID uuid.UUID `json:"category_id"`
	Name       string    `json:"name"`
}

type TransactionResponse struct {
	ID            uuid.UUID `json:"id"`
	UserID        uuid.UUID `json:"user_id"`
	UserAccountID uuid.UUID `json:"user_account_id"`
	Type          string    `json:"type"`
	Amount        float64   `json:"amount"`
	CategoryID    uuid.UUID `json:"category_id"`
	SubcategoryID uuid.UUID `json:"subcategory_id"`
	CreatedAt     string    `json:"created_at"`
}

type CreateTransactionRequest struct {
	Type          string    `json:"type" binding:"required"`
	AccountID     uuid.UUID `json:"accountId" binding:"required"`
	CategoryID    uuid.UUID `json:"categoryId" binding:"required"`
	SubcategoryID uuid.UUID `json:"subcategoryId" binding:"required"`
	Amount        float64   `json:"amount" binding:"required"`
	Note          string    `json:"note"` // Optional, no binding required
	Date          time.Time `json:"date" binding:"required"`
}

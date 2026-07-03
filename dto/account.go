package dto

import (
	"time"

	"github.com/google/uuid"
)

type AccountResponse struct {
	ID               uuid.UUID `json:"id"`
	Name             string    `json:"name"`
	Type             string    `json:"type"`
	SpendableBalance float64   `json:"spendable_balance"`
	SavingsBalance   float64   `json:"savings_balance"`
	Currency         string    `json:"currency"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}

package service

import (
	"expense-tracker-be/dto"
	"expense-tracker-be/storage"
	"fmt"

	"github.com/google/uuid"
)

func GetAllAccountInfo(userID uuid.UUID) []dto.AccountResponse {
	userAccounts, err := storage.GetAccountsByUserID(userID)
	if err != nil {
		fmt.Println("accounts not found")
		return []dto.AccountResponse{}
	}

	accounts := make([]dto.AccountResponse, 0, len(userAccounts))
	for _, account := range userAccounts {
		accountResponse := dto.AccountResponse{
			ID:               account.ID,
			Name:             account.Name,
			SpendableBalance: account.SpendableBalance,
			SavingsBalance:   account.SavingsBalance,
			Currency:         account.Currency,
			CreatedAt:        account.CreatedAt,
		}

		accounts = append(accounts, accountResponse)
	}
	fmt.Println(accounts)
	return accounts
}

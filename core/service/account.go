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

func TransferFunds(userID uuid.UUID, transferDirection string, amount float64, goalID uuid.UUID) (bool, string) {
	accounts, err := storage.GetAccountsByUserID(userID)
	if err != nil {
		return false, "failed to fetch accounts"
	}

	var capital *storage.Account
	var tax *storage.Account
	for _, account := range accounts {
		switch account.Name {
		case "Capital":
			capital = &account
		case "Tax":
			tax = &account
		}
	}

	switch transferDirection {
	case "C2T":
		if amount > capital.SpendableBalance {
			return false, "Insufficient spendable balance"
		}
		err := storage.TransferFundsBetweenAccounts(capital.ID, tax.ID, amount)
		if err != nil {
			return false, "failed to transfer"
		}
	case "T2C":
		if amount > tax.SpendableBalance {
			return false, "Insufficient spendable balance"
		}
		err := storage.TransferFundsBetweenAccounts(tax.ID, capital.ID, amount)
		if err != nil {
			return false, "failed to transfer"
		}
	case "Sp2Sv":
		err := storage.AddFundsToGoalTransaction(capital.ID, goalID, amount)
		if err != nil {
			return false, "failed to transfer"
		}
	case "Sv2Sp":
		err := storage.WithdrawFundsFromGoalTransaction(capital.ID, goalID, amount)
		if err != nil {
			return false, "failed to transfer"
		}
	}
	return true, ""
}

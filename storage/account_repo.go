package storage

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func GetAccountsByUserID(userID uuid.UUID) ([]Account, error) {
	var accounts []Account
	result := DB.Where("user_id = ?", userID).Find(&accounts)
	return accounts, result.Error
}

func GetAccountByID(id uuid.UUID) (*Account, error) {
	var account Account
	result := DB.First(&account, "id = ?", id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &account, nil
}

func AddToAccountSpendableBalance(accountID uuid.UUID, amount float64) error {
	result := DB.Model(&Account{}).
		Where("id = ?", accountID).
		Update("spendable_balance", gorm.Expr("spendable_balance + ?", amount))
	return result.Error
}

func WithdrawAccountSpendableBalance(accountID uuid.UUID, amount float64) error {
	result := DB.Model(&Account{}).
		Where("id = ?", accountID).
		Update("spendable_balance", gorm.Expr("spendable_balance - ?", amount))
	return result.Error
}

func AddToCapitalSavingsBalance(amount float64) error {
	return DB.Model(&Account{}).
		Where("name = ?", "Capital").
		Update("savings_balance", gorm.Expr("savings_balance + ?", amount)).Error
}

func CreateAccount(userID uuid.UUID, name string, spendableBalance, savingsBalance float64, currency string) (*Account, error) {
	account := Account{
		ID:               uuid.New(),
		UserID:           userID,
		Name:             name,
		SpendableBalance: spendableBalance,
		SavingsBalance:   savingsBalance,
		Currency:         currency,
	}
	result := DB.Create(&account)
	return &account, result.Error
}

package storage

import (
	"errors"

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

// Transactional method for drawing from savings to cover spendable deficit
func WithdrawFromSavings(accountID uuid.UUID, goalID uuid.UUID, amount float64) error {
	return DB.Transaction(func(tx *gorm.DB) error {
		// 1. Deduct from Goal balance
		var goal Goal
		if err := tx.First(&goal, "id = ?", goalID).Error; err != nil {
			return err
		}
		if goal.SavedAmount < amount {
			return errors.New("insufficient funds in targeted goal")
		}
		tx.Model(&goal).Update("saved_amount", gorm.Expr("saved_amount - ?", amount))

		// 2. Adjust Account Balances
		var account Account
		if err := tx.First(&account, "id = ?", accountID).Error; err != nil {
			return err
		}
		if account.SavingsBalance < amount {
			return errors.New("insufficient savings pool balance in account")
		}

		err := tx.Model(&account).Updates(map[string]interface{}{
			"savings_balance":   gorm.Expr("savings_balance - ?", amount),
			"spendable_balance": gorm.Expr("spendable_balance + ?", amount),
		}).Error
		if err != nil {
			return err
		}

		// 3. Log the history record
		withdrawal := SavingsWithdrawal{
			ID:        uuid.New(),
			AccountID: account.ID,
			GoalID:    goal.ID,
			Amount:    amount,
		}
		return tx.Create(&withdrawal).Error
	})
}

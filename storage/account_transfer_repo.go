package storage

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

func TransferFundsBetweenAccounts(fromAccountID, toAccountID uuid.UUID, amount float64) error {
	return DB.Transaction(func(tx *gorm.DB) error {
		var fromAccount Account
		if err := tx.First(&fromAccount, "id = ?", fromAccountID).Error; err != nil {
			return err
		}
		if fromAccount.SpendableBalance < amount {
			return fmt.Errorf("insufficient spendable funds in origin account")
		}

		// Deduct from sender spendable pool
		if err := tx.Model(&fromAccount).Update("spendable_balance", gorm.Expr("spendable_balance - ?", amount)).Error; err != nil {
			return err
		}

		// Add to receiver spendable pool
		if err := tx.Model(&Account{}).Where("id = ?", toAccountID).Update("spendable_balance", gorm.Expr("spendable_balance + ?", amount)).Error; err != nil {
			return err
		}

		// Log configuration log event
		transferRecord := AccountTransfer{
			ID:            uuid.New(),
			FromAccountID: fromAccountID,
			ToAccountID:   toAccountID,
			Amount:        amount,
			TransferDate:  time.Now(),
		}
		return tx.Create(&transferRecord).Error
	})
}

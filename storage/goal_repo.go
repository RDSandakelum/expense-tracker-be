package storage

import (
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

func GetGoalsByUserID(userID uuid.UUID) ([]Goal, error) {
	var goals []Goal
	result := DB.Where("user_id = ?", userID).Find(&goals)
	return goals, result.Error
}

func GetGoalByID(id uuid.UUID) (*Goal, error) {
	var goal Goal
	result := DB.First(&goal, "id = ?", id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &goal, nil
}

func AddFundsToGoalTransaction(accountID uuid.UUID, goalID uuid.UUID, amount float64, userID uuid.UUID) error {
	return DB.Transaction(func(tx *gorm.DB) error {
		// 1. Verify spendable balance availability
		var account Account
		if err := tx.First(&account, "id = ?", accountID).Error; err != nil {
			return err
		}
		if account.SpendableBalance < amount {
			return errors.New("insufficient spendable funds to move to savings goal")
		}

		// 2. Deduct from spendable pool, move to savings pool
		err := tx.Model(&account).Updates(map[string]interface{}{
			"spendable_balance": gorm.Expr("spendable_balance - ?", amount),
			"savings_balance":   gorm.Expr("savings_balance + ?", amount),
		}).Error
		if err != nil {
			return err
		}

		// 3. Add to targeted goal tracking balance
		var goal Goal
		if err := tx.First(&goal, "id = ?", goalID).Error; err != nil {
			return err
		}

		newSavedAmount := goal.SavedAmount + amount
		isCompleted := newSavedAmount >= goal.TargetAmount

		err = tx.Model(&goal).Updates(map[string]interface{}{
			"saved_amount": newSavedAmount,
			"is_completed": isCompleted,
		}).Error

		if err != nil {
			return err
		}

		savingWithdrawRecord := &SavingsWithdrawal{
			AccountID: accountID,
			GoalID:    goalID,
			UserID:    userID,
			Amount:    amount,
			Direction: "Saved",
		}
		return tx.Create(savingWithdrawRecord).Error
	})
}

func WithdrawFundsFromGoalTransaction(accountID uuid.UUID, goalID uuid.UUID, amount float64, userID uuid.UUID) error {
	return DB.Transaction(func(tx *gorm.DB) error {
		// 1. Verify spendable balance availability
		var account Account
		if err := tx.First(&account, "id = ?", accountID).Error; err != nil {
			return err
		}
		if account.SavingsBalance < amount {
			return errors.New("insufficient spendable funds to move to savings goal")
		}

		// 2. Deduct from spendable pool, move to savings pool
		err := tx.Model(&account).Updates(map[string]interface{}{
			"spendable_balance": gorm.Expr("spendable_balance + ?", amount),
			"savings_balance":   gorm.Expr("savings_balance - ?", amount),
		}).Error
		if err != nil {
			return err
		}

		// 3. Add to targeted goal tracking balance
		var goal Goal
		if err := tx.First(&goal, "id = ?", goalID).Error; err != nil {
			return err
		}

		newSavedAmount := goal.SavedAmount - amount
		isCompleted := newSavedAmount >= goal.TargetAmount

		err = tx.Model(&goal).Updates(map[string]interface{}{
			"saved_amount": newSavedAmount,
			"is_completed": isCompleted,
		}).Error

		if err != nil {
			return err
		}

		savingWithdrawRecord := &SavingsWithdrawal{
			AccountID: accountID,
			GoalID:    goalID,
			UserID:    userID,
			Amount:    amount,
			Direction: "Withdrawn",
		}
		return tx.Create(savingWithdrawRecord).Error
	})
}

package storage

import "github.com/google/uuid"

func GetAllSavingWithdrawals(userID uuid.UUID) ([]SavingsWithdrawal, error) {
	var savings []SavingsWithdrawal
	result := DB.Where("user_id = ?", userID).Find(&savings)
	return savings, result.Error
}

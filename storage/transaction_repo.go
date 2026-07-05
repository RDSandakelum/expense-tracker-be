package storage

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

func GetTransactionsByMonth(userID uuid.UUID, startOfMonth, endOfMonth time.Time, limit, offset int) ([]Transaction, int64, error) {
	var transactions []Transaction
	var totalCount int64

	// Base query matching your filtering criteria
	baseQuery := DB.Model(&Transaction{}).
		Where("user_id = ? AND transaction_date BETWEEN ? AND ?", userID, startOfMonth, endOfMonth)

	// 1. Fetch total count matching criteria before applying pagination limits
	if err := baseQuery.Count(&totalCount).Error; err != nil {
		return nil, 0, err
	}

	// 2. Fetch paginated records with preloaded relations
	result := baseQuery.
		Preload("Account").
		Preload("SubCategory.Category").
		Order("transaction_date DESC").
		Limit(limit).
		Offset(offset).
		Find(&transactions)

	return transactions, totalCount, result.Error
}

func GetTransactionsByMonthWithoutPagination(userID uuid.UUID, startOfMonth, endOfMonth time.Time) ([]Transaction, error) {
	var transactions []Transaction

	// Base query matching your filtering criteria
	baseQuery := DB.Model(&Transaction{}).
		Where("user_id = ? AND transaction_date BETWEEN ? AND ?", userID, startOfMonth, endOfMonth)

	// 2. Fetch paginated records with preloaded relations
	result := baseQuery.
		Preload("Account").
		Preload("SubCategory").
		Preload("SubCategory.Category").
		Find(&transactions)

	return transactions, result.Error
}

func GetTransactionByID(id uuid.UUID) (*Transaction, error) {
	var transaction Transaction
	result := DB.First(&transaction, "id = ?", id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &transaction, nil
}

func CreateTransactionRecordWithoutAccountUpdate(userID uuid.UUID, accountID uuid.UUID, subCatID *uuid.UUID, txType string, amount float64, note string, transactionDate time.Time) (*Transaction, error) {
	txRecord := Transaction{
		ID:              uuid.New(),
		UserID:          userID,
		AccountID:       accountID,
		SubCategoryID:   subCatID,
		Type:            txType,
		Amount:          amount,
		Note:            note,
		TransactionDate: transactionDate,
	}

	// Wrap execution in an internal transaction to alter standard balance changes alongside log placement

	if err := DB.Create(&txRecord).Error; err != nil {
		return nil, err
	}

	// Standard user transaction affects the spendable pool of the designated account

	return &txRecord, nil
}
func CreateTransactionRecord(userID uuid.UUID, accountID uuid.UUID, subCatID *uuid.UUID, txType string, amount float64, note string, transactionDate time.Time) (*Transaction, error) {
	txRecord := Transaction{
		ID:              uuid.New(),
		UserID:          userID,
		AccountID:       accountID,
		SubCategoryID:   subCatID,
		Type:            txType,
		Amount:          amount,
		Note:            note,
		TransactionDate: transactionDate,
	}

	// Wrap execution in an internal transaction to alter standard balance changes alongside log placement
	err := DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&txRecord).Error; err != nil {
			return err
		}

		// Dynamic modifier expression based on Debit/Credit rules
		balanceChange := amount
		if txType == "Debit" {
			balanceChange = -amount
		}

		// Standard user transaction affects the spendable pool of the designated account
		return tx.Model(&Account{}).Where("id = ?", accountID).
			Update("spendable_balance", gorm.Expr("spendable_balance + ?", balanceChange)).Error
	})

	return &txRecord, err
}

func UpdateTransactionRecord(txRecord *Transaction) error {
	return DB.Save(txRecord).Error
}

func DeleteTransactionRecord(id uuid.UUID) error {
	result := DB.Delete(&Transaction{}, "id = ?", id)
	return result.Error
}

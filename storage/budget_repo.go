package storage

import (
	"log"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

func GetBudgetsByUserIDAndMonth(userID uuid.UUID, month time.Time) ([]Budget, error) {
	var budgets []Budget
	firstOfMonth := time.Date(month.Year(), month.Month(), 1, 0, 0, 0, 0, month.Location())

	// Chain Preload with the exact name of the struct field ('SubCategory')
	result := DB.Preload("SubCategory").
		Where("user_id = ? AND budget_month = ?", userID, firstOfMonth).
		Find(&budgets)

	return budgets, result.Error
}

func GetBudgetBySubCategoryID(userID uuid.UUID, subCategoryID uuid.UUID, month time.Time) (*Budget, error) {
	var budget Budget
	firstOfMonth := time.Date(month.Year(), month.Month(), 1, 0, 0, 0, 0, month.Location())
	result := DB.Where("user_id = ? AND sub_category_id = ? AND budget_month = ?", userID, subCategoryID, firstOfMonth).First(&budget)
	if result.Error != nil {
		return nil, result.Error
	}
	return &budget, nil
}

func AddToBudget(userID uuid.UUID, subCategoryID uuid.UUID, month time.Time, amount float64, txType string) error {
	// Normalize date context to the 1st of the target month
	firstOfMonth := time.Date(month.Year(), month.Month(), 1, 0, 0, 0, 0, month.Location())

	// Atomically increment the current_spend column for the matching user record
	changeAmount := amount
	if txType == "Credit" {
		changeAmount = -amount // Subtract for credits
	}
	result := DB.Model(&Budget{}).
		Where("user_id = ? AND sub_category_id = ? AND budget_month = ?", userID, subCategoryID, firstOfMonth).
		Update("current_spend", gorm.Expr("current_spend + ?", changeAmount))

	if result.Error != nil {
		return result.Error
	}

	// Optional: Handle edge case where a transaction is logged before the monthly budget row was initialized
	if result.RowsAffected == 0 {
		// If you want to automatically initialize the budget from the template if it's missing:
		if err := InitializeMonthlyBudgetFromTemplate(userID, firstOfMonth); err != nil {
			return err
		}
	}

	return nil
}

func InitializeMonthlyBudgetFromTemplate(userID uuid.UUID, targetMonth time.Time) error {
	// 1. Normalize dates to the 1st of the month
	currentMonthStart := time.Date(targetMonth.Year(), targetMonth.Month(), 1, 0, 0, 0, 0, targetMonth.Location())
	nextMonthStart := currentMonthStart.AddDate(0, 1, 0)
	prevMonthStart := currentMonthStart.AddDate(0, -1, 0)

	return DB.Transaction(func(tx *gorm.DB) error {
		// 2. Fetch the user's baseline budget templates
		var templates []BudgetTemplate
		if err := tx.Where("user_id = ?", userID).Find(&templates).Error; err != nil {
			log.Printf("[BudgetInit] ERROR: Failed to fetch templates for user %s: %v", userID, err)
			return err
		}
		log.Printf("[BudgetInit] Found %d budget templates configuration rules", len(templates))

		// 3. Fetch the previous month's actual budgets to calculate carry-overs
		var prevBudgets []Budget
		if err := tx.Where("user_id = ? AND budget_month = ?", userID, prevMonthStart).Find(&prevBudgets).Error; err != nil {
			log.Printf("[BudgetInit] ERROR: Failed to fetch previous month (%s) budgets: %v", prevMonthStart.Format("2006-01"), err)
			return err
		}
		log.Printf("[BudgetInit] Found %d existing budget snapshots from previous month", len(prevBudgets))

		// Map previous budgets by SubCategoryID for fast structural lookups
		prevBudgetMap := make(map[uuid.UUID]Budget)
		for _, pb := range prevBudgets {
			prevBudgetMap[pb.SubCategoryID] = pb
		}

		// 4. Process each template entry to initialize the target month
		for _, tmpl := range templates {
			var carryOverAmount float64 = 0.00

			// Calculate rollover if allowed
			if tmpl.AllowCarryOver {
				if oldBudget, exists := prevBudgetMap[tmpl.SubCategoryID]; exists {
					totalPool := oldBudget.AllocatedAmount + oldBudget.CarriedOverAmount
					leftOver := totalPool - oldBudget.CurrentSpend

					if leftOver > 0 {
						carryOverAmount = leftOver
						log.Printf("[BudgetInit][SubCat: %s] Carry-over detected: %2f (Pool: %2f, Spent: %2f)",
							tmpl.SubCategoryID, carryOverAmount, totalPool, oldBudget.CurrentSpend)
					}
				} else {
					log.Printf("[BudgetInit][SubCat: %s] Rollover enabled but no historical record found for previous month", tmpl.SubCategoryID)
				}
			} else {
				if oldBudget, exists := prevBudgetMap[tmpl.SubCategoryID]; exists {
					totalPool := oldBudget.AllocatedAmount + oldBudget.CarriedOverAmount
					leftOver := totalPool - oldBudget.CurrentSpend
					if leftOver > 0 {
						savingAmount := leftOver
						err := AddToCapitalSavingsBalance(savingAmount)
						if err != nil {
							log.Printf("[BudgetInit][SubCat: %s] ERROR: Failed to add to capital savings balance: %v", tmpl.SubCategoryID, err)
						}
						log.Printf("[BudgetInit][SubCat: %s] added to capital savings balance: %v", tmpl.SubCategoryID, err)
					}
				}
			}

			// Aggregate pre-existing transactions for this specific subcategory for the current month
			var preExistingSpend float64 = 0.00
			err := tx.Model(&Transaction{}).
				Where("user_id = ? AND sub_category_id = ? AND type = 'Debit' AND transaction_date >= ? AND transaction_date < ?",
					userID, tmpl.SubCategoryID, currentMonthStart, nextMonthStart).
				Select("COALESCE(SUM(amount), 0)").
				Scan(&preExistingSpend).Error
			if err != nil {
				log.Printf("[BudgetInit][SubCat: %s] ERROR: Failed to aggregate pre-existing transactions: %v", tmpl.SubCategoryID, err)
				return err
			}
			if preExistingSpend > 0 {
				log.Printf("[BudgetInit][SubCat: %s] Found pre-existing transactions totaling: %.2f", tmpl.SubCategoryID, preExistingSpend)
			}

			// 5. Build the new monthly record
			newBudget := Budget{
				UserID:            userID,
				SubCategoryID:     tmpl.SubCategoryID,
				BudgetMonth:       currentMonthStart,
				AllocatedAmount:   tmpl.AllocatedAmount,
				CarriedOverAmount: carryOverAmount,
				CurrentSpend:      preExistingSpend,
			}

			// Upsert strategy: If it already exists for this specific month, update it; otherwise, insert it.
			err = tx.Where(Budget{UserID: userID, SubCategoryID: tmpl.SubCategoryID, BudgetMonth: currentMonthStart}).
				Attrs(Budget{ID: uuid.New()}).
				FirstOrCreate(&newBudget).Error
			if err != nil {
				log.Printf("[BudgetInit][SubCat: %s] ERROR: FirstOrCreate transaction step failed: %v", tmpl.SubCategoryID, err)
				return err
			}

			// If it already existed when FirstOrCreate ran, update fields accurately
			if tx.RowsAffected == 0 {
				log.Printf("[BudgetInit][SubCat: %s] Record already exists for this month. Forcing update synchronization...", tmpl.SubCategoryID)
				err = tx.Model(&Budget{}).
					Where("user_id = ? AND sub_category_id = ? AND budget_month = ?", userID, tmpl.SubCategoryID, currentMonthStart).
					Updates(map[string]interface{}{
						"allocated_amount":    tmpl.AllocatedAmount,
						"carried_over_amount": carryOverAmount,
						"current_spend":       preExistingSpend,
					}).Error
				if err != nil {
					log.Printf("[BudgetInit][SubCat: %s] ERROR: Fallback synchronization update failed: %v", tmpl.SubCategoryID, err)
					return err
				}
			} else {
				log.Printf("[BudgetInit][SubCat: %s] Successfully created a fresh monthly budget entry record.", tmpl.SubCategoryID)
			}
		}

		log.Printf("[BudgetInit] Transaction block completed successfully for user %s", userID)
		return nil
	})
}

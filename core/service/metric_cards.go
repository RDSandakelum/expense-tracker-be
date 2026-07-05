package service

import (
	"expense-tracker-be/storage"
	"time"

	"github.com/google/uuid"
)

type DashboardSummary struct {
	MonthIncome       float64 `json:"month_income"`
	MonthExpenses     float64 `json:"month_expenses"`
	RemainingExpenses float64 `json:"remaining_expenses"`
	MonthSavings      float64 `json:"month_savings"`
	CarriedOver       float64 `json:"carried_over"`
}

func GetMetricCards(userID uuid.UUID) DashboardSummary {
	cards := DashboardSummary{
		MonthIncome:       0.00,
		MonthExpenses:     0.00,
		RemainingExpenses: 0.00,
		MonthSavings:      0.00,
		CarriedOver:       0.00,
	}
	now := time.Now()
	startOfMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	endOfMonth := startOfMonth.AddDate(0, 1, 0).Add(-time.Nanosecond)
	transctions, err := storage.GetTransactionsByMonthWithoutPagination(userID, startOfMonth, endOfMonth)
	if err != nil {
		return cards
	}

	expenses := 0.00
	income := 0.00
	carriedOver := 0.00
	allocatedExpenses := 0.00
	for _, transaction := range transctions {
		if transaction.SubCategory.Category.Name == "Income" {
			income += transaction.Amount
		}
	}

	cards.MonthIncome = income

	budgets, err := storage.GetBudgetsByUserIDAndMonth(userID, startOfMonth)
	if err != nil {
		return cards
	}

	for _, budget := range budgets {
		carriedOver += budget.CarriedOverAmount
		allocatedExpenses += budget.AllocatedAmount
		expenses += budget.CurrentSpend
	}

	cards.MonthExpenses = expenses
	cards.CarriedOver = carriedOver
	cards.MonthSavings = income - expenses
	cards.RemainingExpenses = allocatedExpenses - expenses

	return cards
}

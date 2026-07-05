package handlers

import (
	"expense-tracker-be/core/service"
	"expense-tracker-be/storage"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type DashboardSummary struct {
	MonthIncome       float64 `json:"month_income"`
	MonthExpenses     float64 `json:"month_expenses"`
	RemainingExpenses float64 `json:"remaining_expenses"`
	MonthSavings      float64 `json:"month_savings"`
	CarriedOver       float64 `json:"carried_over"`
}

type MonthlyCategoryTrend struct {
	Month    string  `json:"month"`
	Income   float64 `json:"income"`
	Expenses float64 `json:"expenses"`
	Balance  float64 `json:"balance"`
}

type BudgetItem struct {
	ID         uuid.UUID `json:"id"`
	Category   string    `json:"category"`
	Allocated  float64   `json:"allocated"`
	Spent      float64   `json:"spent"`
	Remaining  float64   `json:"remaining"`
	Percentage float64   `json:"percentage"`
}

func GetSummaryCards(c *gin.Context) {

	userIDInterface, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found in context"})
		return
	}

	userIDStr, ok := userIDInterface.(string)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID format"})
		return
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	dashboardSummary := service.GetMetricCards(userID)

	// Dummy summary data
	summary := DashboardSummary{
		MonthIncome:       dashboardSummary.MonthIncome,
		MonthExpenses:     dashboardSummary.MonthExpenses,
		RemainingExpenses: dashboardSummary.RemainingExpenses,
		MonthSavings:      dashboardSummary.MonthSavings,
		CarriedOver:       dashboardSummary.CarriedOver,
	}

	c.JSON(http.StatusOK, summary)
}

func GetMonthlyTrend(c *gin.Context) {
	// Get user ID from context (set by AuthMiddleware)
	_, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "User ID not found in context"})
		return
	}

	// Dummy monthly trend data (last 12 months)
	now := time.Now()
	trends := []MonthlyCategoryTrend{}

	for i := 11; i >= 0; i-- {
		month := now.AddDate(0, -i, 0)
		trends = append(trends, MonthlyCategoryTrend{
			Month:    month.Format("Jan 2006"),
			Income:   4500.00,
			Expenses: 900.00 + float64(i*50),
			Balance:  3600.00 - float64(i*50),
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"trends": trends,
	})
}

func GetCategoryBreakdownDashboard(c *gin.Context) {
	// Get user ID from context (set by AuthMiddleware)
	userIDInterface, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "User ID not found in context"})
		return
	}

	userIDStr, ok := userIDInterface.(string)
	_, err := uuid.Parse(userIDStr)
	if !ok || err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid request user credentials"})
		return
	}

	// 1. Extract query parameters with fallbacks to the current date values
	now := time.Now()
	yearStr := c.DefaultQuery("year", strconv.Itoa(now.Year()))
	monthStr := c.DefaultQuery("month", strconv.Itoa(int(now.Month())))

	// 2. Parse year value
	year, err := strconv.Atoi(yearStr)
	if err != nil || year < 1970 {
		year = now.Year() // Fallback if corrupt
	}

	// 3. Parse month value
	month, err := strconv.Atoi(monthStr)
	if err != nil || month < 1 || month > 12 {
		month = int(now.Month()) // Fallback if corrupt
	}

	// 4. Construct a time.Time object pointing explicitly to the 1st of that month
	// Using time.Local matches your connection string timezone setup context
	requestedMonthStart := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.Local)

	budgets, err := storage.GetBudgetsByUserIDAndMonth(uuid.MustParse(userIDStr), requestedMonthStart)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to fetch budgets: %v", err)})
		return
	}

	budgetItems := []BudgetItem{}
	for _, budget := range budgets {
		budgetItems = append(budgetItems, BudgetItem{
			ID:         budget.ID,
			Category:   budget.SubCategory.Name,
			Allocated:  budget.AllocatedAmount,
			Spent:      budget.CurrentSpend,
			Remaining:  (budget.AllocatedAmount - budget.CurrentSpend),
			Percentage: (budget.CurrentSpend / budget.AllocatedAmount) * 100,
		})
	}

	// Dummy category breakdown data
	c.JSON(http.StatusOK, gin.H{
		"categories": budgetItems,
	})
}

func GetBudgets(c *gin.Context) {
	// 1. Get user ID from context
	userIDInterface, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "User ID not found in context"})
		return
	}

	userIDStr, ok := userIDInterface.(string)
	userID, err := uuid.Parse(userIDStr)
	if !ok || err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user ID"})
		return
	}

	// 2. Parse and normalize query parameters (defaulting to current month/year)
	now := time.Now()
	monthStr := c.DefaultQuery("month", strconv.Itoa(int(now.Month())))
	yearStr := c.DefaultQuery("year", strconv.Itoa(now.Year()))

	month, err := strconv.Atoi(monthStr)
	if err != nil || month < 1 || month > 12 {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid month parameter"})
		return
	}

	year, err := strconv.Atoi(yearStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid year parameter"})
		return
	}

	targetMonth := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.Local)

	// 3. Performance Optimization: Pull all budgets + preloaded subcategories in one batch
	budgets, err := storage.GetBudgetsByUserIDAndMonth(userID, targetMonth)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch budgets"})
		return
	}

	// 4. Aggregate subcategory budgets up to their parent categories dynamically
	// This ensures that if multiple subcategories exist under "Food & Dining", they roll up cleanly.
	type catTotals struct {
		Allocated float64
		Spent     float64
	}
	categoryRollups := make(map[string]catTotals)

	for _, b := range budgets {
		if b.SubCategory == nil || b.SubCategory.Category == nil {
			continue // Skip orphans lacking relational mapping contexts
		}

		categoryName := b.SubCategory.Category.Name

		// Use total pool (Allocation + Rollover) if you want tracking against absolute limits
		current := categoryRollups[categoryName]
		current.Allocated += b.AllocatedAmount
		current.Spent += b.CurrentSpend

		categoryRollups[categoryName] = current
	}

	// 5. Build final response list using your existing BudgetItem structural signature
	budgetItems := []BudgetItem{}
	for catName, totals := range categoryRollups {
		remaining := totals.Allocated - totals.Spent

		var percentage float64 = 0.00
		if totals.Allocated > 0 {
			percentage = (totals.Spent / totals.Allocated) * 100
		}

		budgetItems = append(budgetItems, BudgetItem{
			ID:         uuid.New(), // Group baseline reference identifier
			Category:   catName,
			Allocated:  totals.Allocated,
			Spent:      totals.Spent,
			Remaining:  remaining,
			Percentage: percentage,
		})
	}

	// 6. Return dynamic tracking states to frontend
	c.JSON(http.StatusOK, gin.H{
		"budgets": budgetItems,
	})
}

func GetSettings(c *gin.Context) {
	// Get user ID from context (set by AuthMiddleware)
	userIDInterface, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "User ID not found in context"})
		return
	}

	// Dummy user settings data
	settings := gin.H{
		"user_id":       userIDInterface,
		"theme":         "light",
		"currency":      "USD",
		"language":      "en",
		"notifications": true,
		"email_reports": true,
		"two_factor":    false,
	}

	c.JSON(http.StatusOK, settings)
}

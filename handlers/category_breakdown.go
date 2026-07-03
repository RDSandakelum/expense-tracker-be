package handlers

import (
	"expense-tracker-be/storage"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type CategoryBreakdownResponse struct {
	ID            uuid.UUID              `json:"id"`
	Name          string                 `json:"name"`
	Allocated     float64                `json:"allocated"`
	Spent         float64                `json:"spent"`
	Remaining     float64                `json:"remaining"`
	Percentage    float64                `json:"percentage"`
	Subcategories []SubcategoryBreakdown `json:"subcategories"`
	CreatedAt     time.Time              `json:"created_at"`
	UpdatedAt     time.Time              `json:"updated_at"`
}

type SubcategoryBreakdown struct {
	ID         uuid.UUID `json:"id"`
	Name       string    `json:"name"`
	Allocated  float64   `json:"allocated"`
	Spent      float64   `json:"spent"`
	Remaining  float64   `json:"remaining"`
	Percentage float64   `json:"percentage"`
}

func GetCategoryBreakdown(c *gin.Context) {
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

	// 2. Parse and normalize query parameters
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

	// 3. Performance Optimization: Pull all budgets + preloaded relations at once
	// This replaces your inner loop N+1 database queries!
	budgets, err := storage.GetBudgetsByUserIDAndMonth(userID, targetMonth)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch budget breakdowns"})
		return
	}

	// 4. Group data into a relational map structure in memory
	// Key: Main Category ID, Value: List of calculated subcategory structures
	subBreakdownMap := make(map[uuid.UUID][]SubcategoryBreakdown)
	categoryObjectsMap := make(map[uuid.UUID]storage.Category)

	var globalAllocated float64 = 0.00
	var globalSpent float64 = 0.00

	for _, b := range budgets {
		if b.SubCategory == nil {
			continue // Skip orphans if data integrity is broken
		}

		subCat := b.SubCategory

		// Math: Pool includes both what was allocated this month + rollover funds
		totalAllocatedPool := b.AllocatedAmount + b.CarriedOverAmount
		remaining := totalAllocatedPool - b.CurrentSpend

		// Prevent division by zero runtime NaN errors
		var percentage float64 = 0.00
		if totalAllocatedPool > 0 {
			percentage = (b.CurrentSpend / totalAllocatedPool) * 100
		}

		// Accumulate global totals
		globalAllocated += b.AllocatedAmount // or totalAllocatedPool based on your business metric preference
		globalSpent += b.CurrentSpend

		// Build subcategory response item
		subItem := SubcategoryBreakdown{
			ID:         subCat.ID,
			Name:       subCat.Name,
			Allocated:  b.AllocatedAmount,
			Spent:      b.CurrentSpend,
			Remaining:  remaining,
			Percentage: percentage,
		}

		// Keep track of parent category mapping contexts
		parentID := subCat.CategoryID
		subBreakdownMap[parentID] = append(subBreakdownMap[parentID], subItem)

		// Ensure we have access to the category name later
		if subCat.Category != nil {
			categoryObjectsMap[parentID] = *subCat.Category
		} else {
			category, err := storage.GetCategoryByID(parentID)
			if err == nil {
				categoryObjectsMap[parentID] = *category
			}
		}
	}

	// 5. Construct the final hierarchy payload
	categoriesBreakdownResponse := []CategoryBreakdownResponse{}
	for catID, subs := range subBreakdownMap {
		var catAllocated float64 = 0
		var catSpent float64 = 0

		// Calculate rollups dynamically for the main category grouping
		for _, s := range subs {
			catAllocated += s.Allocated
			catSpent += s.Spent
		}

		catRemaining := catAllocated - catSpent
		var catPercentage float64 = 0
		if catAllocated > 0 {
			catPercentage = (catSpent / catAllocated) * 100
		}

		catMeta := categoryObjectsMap[catID]

		categoriesBreakdownResponse = append(categoriesBreakdownResponse, CategoryBreakdownResponse{
			ID:            catID,
			Name:          catMeta.Name,
			Allocated:     catAllocated,
			Spent:         catSpent,
			Remaining:     catRemaining,
			Percentage:    catPercentage,
			Subcategories: subs,
			CreatedAt:     catMeta.CreatedAt,
			UpdatedAt:     time.Now(),
		})
	}

	// 6. Send perfectly calculated state back to your React frontend context!
	c.JSON(http.StatusOK, gin.H{
		"month":          month,
		"year":           year,
		"categories":     categoriesBreakdownResponse,
		"totalAllocated": globalAllocated,
		"totalSpent":     globalSpent,
		"totalRemaining": (globalAllocated - globalSpent),
	})
}

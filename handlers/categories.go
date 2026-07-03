package handlers

import (
	"expense-tracker-be/core/service"
	"expense-tracker-be/dto"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func GetCategories(c *gin.Context) {
	// Get user ID from context (set by AuthMiddleware)
	userIDInf, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "User ID not found in context"})
		return
	}

	userIDStr := userIDInf.(string)
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid request body"})
		return
	}
	categories := service.GetCategoriesAndSubCategories(userID)

	c.JSON(http.StatusOK, gin.H{
		"categories": categories,
	})
}

func GetSubcategories(c *gin.Context) {
	// Get user ID from context (set by AuthMiddleware)
	_, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "User ID not found in context"})
		return
	}

	// Dummy subcategories
	subcategories := []dto.SubcategoryResponse{
		{
			ID:         uuid.New(),
			CategoryID: uuid.New(),
			Name:       "Groceries",
		},
		{
			ID:         uuid.New(),
			CategoryID: uuid.New(),
			Name:       "Restaurants",
		},
		{
			ID:         uuid.New(),
			CategoryID: uuid.New(),
			Name:       "Gas",
		},
		{
			ID:         uuid.New(),
			CategoryID: uuid.New(),
			Name:       "Movies",
		},
		{
			ID:         uuid.New(),
			CategoryID: uuid.New(),
			Name:       "Electricity",
		},
	}

	c.JSON(http.StatusOK, gin.H{
		"subcategories": subcategories,
	})
}

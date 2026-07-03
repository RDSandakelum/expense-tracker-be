package handlers

import (
	"expense-tracker-be/core/service"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func GetAccounts(c *gin.Context) {
	_, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "User ID not found in context"})
		return
	}

	userIDInf, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "User ID not found in context"})
		return
	}

	userIDstr := userIDInf.(string)
	userID, err := uuid.Parse(userIDstr)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid request body"})
		return
	}

	accounts := service.GetAllAccountInfo(userID)

	c.JSON(http.StatusOK, gin.H{
		"accounts": accounts,
	})
}

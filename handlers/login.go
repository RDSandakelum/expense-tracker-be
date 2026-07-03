package handlers

import (
	"expense-tracker-be/core/domain"
	"expense-tracker-be/core/service"
	"expense-tracker-be/dto"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func LoginHandler(c *gin.Context) {
	var req dto.LoginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Email and password are required.",
		})
		return
	}

	user := service.GetUserByEmail(req.Email)
	if user == nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "Invalid email or password.",
		})
		return
	}
	if req.Password != user.PasswordHash {
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "Invalid email or password.",
		})
		return
	}

	token, err := createJWT(*user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Could not create login session.",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token": token,
		"user": gin.H{
			"id":    user.ID,
			"email": user.Email,
			"name":  user.FirstName,
		},
	})
}

func createJWT(user domain.User) (string, error) {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = "dev-secret-change-me"
	}

	claims := jwt.MapClaims{
		"sub":   user.ID,
		"email": user.Email,
		"exp":   time.Now().Add(24 * time.Hour).Unix(),
		"iat":   time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

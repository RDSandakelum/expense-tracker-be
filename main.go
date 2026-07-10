package main

import (
	"expense-tracker-be/handlers"
	"expense-tracker-be/middleware"
	"expense-tracker-be/storage"
	"os"

	"github.com/gin-gonic/gin"
)

func main() {
	storage.ConnectDatabase()
	r := gin.Default()
	r.Use(middleware.CorsMiddleware())

	// Public routes
	r.POST("/api/login", handlers.LoginHandler)

	// Protected routes - grouped with auth middleware
	protected := r.Group("/api/users/:userId")
	protected.Use(middleware.AuthMiddleware())

	// Dashboard endpoints
	protected.GET("/summary-cards", handlers.GetSummaryCards)
	protected.GET("/monthly-trend", handlers.GetMonthlyTrend)
	protected.GET("/category-breakdown", handlers.GetCategoryBreakdownDashboard)
	protected.GET("/budgets", handlers.GetBudgets)
	protected.GET("/goals", handlers.GetGoals)
	protected.GET("/settings", handlers.GetSettings)

	// Categories endpoints
	protected.GET("/categories", handlers.GetCategories)
	protected.GET("/categories/breakdown", handlers.GetCategoryBreakdown)

	// Accounts endpoints
	protected.GET("/accounts", handlers.GetAccounts)
	protected.POST("/accounts/transfer", handlers.TransferFunds)
	protected.GET("/account/transfers", handlers.GetAccountTransfers)

	// Transactions endpoints
	protected.GET("/transactions", handlers.GetCategoryTransactionsList)
	protected.POST("/transactions", handlers.CreateTransaction)
	protected.PUT("/transactions/:transactionId", handlers.UpdateTransaction)
	protected.DELETE("/transactions/:transactionId", handlers.DeleteTransaction)

	//initialize budget
	protected.GET("/budget/initialize", handlers.InitializeBudget)
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	r.Run("0.0.0.0:" + port)
}

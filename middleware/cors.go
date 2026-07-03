package middleware

import (
	"os"
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func CorsMiddleware() gin.HandlerFunc {
	// 1. Start with the safe library default config structure
	config := cors.DefaultConfig()

	var allowedOrigins []string
	envOrigins := os.Getenv("ALLOWED_ORIGINS")

	if envOrigins != "" {
		// Only split if the environment variable actually has content
		for _, origin := range strings.Split(envOrigins, ",") {
			trimmed := strings.TrimSpace(origin)
			if trimmed != "" {
				allowedOrigins = append(allowedOrigins, trimmed)
			}
		}
	} else {
		// Fallback for local development if you forget to set the env variable
		allowedOrigins = []string{"http://localhost:5173"}
	}

	// 2. Assign your values to the safe config object
	config.AllowOrigins = allowedOrigins
	config.AllowMethods = []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"}
	config.AllowHeaders = []string{"Origin", "Content-Type", "Accept", "Authorization", "X-Requested-With"}
	config.ExposeHeaders = []string{"Content-Length"}
	config.AllowCredentials = true
	config.MaxAge = 12 * time.Hour

	return cors.New(config)
}

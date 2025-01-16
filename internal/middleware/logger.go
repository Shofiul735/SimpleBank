package middleware

import (
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// Logger middleware for Gin
func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		// Process request
		c.Next()

		// Log after request is processed
		log.Printf(
			"%s %s %s",
			c.Request.Method,
			c.Request.RequestURI,
			time.Since(start),
		)
	}
}

// Recoverer middleware for Gin
func Recoverer() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				c.AbortWithStatus(http.StatusInternalServerError)
			}
		}()
		c.Next()
	}
}

// CORS middleware for Gin
func CORS() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusOK)
			return
		}

		c.Next()
	}
}

// Authentication middleware for Gin
func Authentication() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get token from Authorization header
		token := c.GetHeader("Authorization")
		if token == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Unauthorized",
			})
			return
		}

		// TODO: Implement proper token validation
		// This is just a placeholder. You should implement proper JWT validation

		c.Next()
	}
}

// ChainGin chains multiple Gin middleware functions together
func ChainGin(middlewares ...gin.HandlerFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Create a new slice to hold the middleware chain
		chain := make([]gin.HandlerFunc, len(middlewares))
		copy(chain, middlewares)

		// Execute each middleware in the chain
		for _, middleware := range chain {
			middleware(c)
			// If the request was aborted in a middleware, stop the chain
			if c.IsAborted() {
				return
			}
		}
	}
}

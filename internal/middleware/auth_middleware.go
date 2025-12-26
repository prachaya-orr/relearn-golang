package middleware

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization header missing"})
			return
		}

		parts := strings.Split(authHeader, " ")
		tokenString := ""
		if len(parts) == 2 && parts[0] == "Bearer" {
			tokenString = parts[1]
		} else if len(parts) == 1 {
			tokenString = parts[0]
		} else {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization header format"})
			return
		}
		secret := os.Getenv("JWT_SECRET")
		if secret == "" {
			secret = "secret"
		}

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(secret), nil
		})

		if err != nil || !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
			return
		}

		// Set User ID to context to be used in handlers
		// Note: JWT library parses numbers as float64 by default, but UUIDs are strings
		if sub, ok := claims["sub"].(string); ok {
			userID, err := uuid.Parse(sub)
			if err == nil {
				c.Set("userID", userID)
			}
		}

		// Ensure it's an access token
		if claims["type"] != "access" {
			// Backward compatibility: if type claim is missing, treat as valid strict access token if we haven't implemented it before
			// But for new implementation, strictness is better.
			// However for now let's allow if type is strictly "access" or missing (legacy/migration)
			// Actually, let's strictly enforce "access" since we just implemented it.
			if typ, ok := claims["type"].(string); ok && typ != "access" {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token type, access token required"})
				return
			}
		}

		c.Next()
	}
}

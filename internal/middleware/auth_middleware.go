package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"kantin-mardira-api/internal/utils"
)

type TokenRevocationChecker interface {
	IsRevoked(jti string) (bool, error)
}

func AuthMiddleware(tokenChecker TokenRevocationChecker) gin.HandlerFunc {
	return func(c *gin.Context) {
		authorizationHeader := c.GetHeader("Authorization")
		if authorizationHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"message": "unauthorized",
			})
			c.Abort()
			return
		}

		parts := strings.SplitN(authorizationHeader, " ", 2)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
			c.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"message": "unauthorized",
			})
			c.Abort()
			return
		}

		claims, err := utils.ParseToken(parts[1])
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"message": "unauthorized",
			})
			c.Abort()
			return
		}

		if tokenChecker != nil {
			isRevoked, err := tokenChecker.IsRevoked(claims.ID)
			if err != nil || isRevoked {
				c.JSON(http.StatusUnauthorized, gin.H{
					"success": false,
					"message": "unauthorized",
				})
				c.Abort()
				return
			}
		}

		c.Set("user_id", claims.UserID)
		c.Set("role", claims.Role)
		c.Set("token_jti", claims.ID)
		c.Set("claims", claims)
		c.Next()
	}
}

func RoleMiddleware(roles ...string) gin.HandlerFunc {
	allowedRoles := make(map[string]struct{}, len(roles))
	for _, role := range roles {
		allowedRoles[strings.ToLower(strings.TrimSpace(role))] = struct{}{}
	}

	return func(c *gin.Context) {
		roleValue, exists := c.Get("role")
		if !exists {
			c.JSON(http.StatusForbidden, gin.H{
				"success": false,
				"message": "forbidden",
			})
			c.Abort()
			return
		}

		currentRole, ok := roleValue.(string)
		if !ok {
			c.JSON(http.StatusForbidden, gin.H{
				"success": false,
				"message": "forbidden",
			})
			c.Abort()
			return
		}

		if _, allowed := allowedRoles[strings.ToLower(strings.TrimSpace(currentRole))]; !allowed {
			c.JSON(http.StatusForbidden, gin.H{
				"success": false,
				"message": "forbidden",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}
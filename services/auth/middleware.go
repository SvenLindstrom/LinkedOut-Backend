package auth

import (
	"linkedout/services/auth/utils/JWT"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func TokenMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		bearerToken := c.GetHeader("Authorization")

		if bearerToken == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "unauthorized"})
			return
		}

		reqToken := strings.Split(bearerToken, " ")
		if len(reqToken) < 2 {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "unauthorized"})
			return
		}

		token, err := jwt.Verify(reqToken[1], jwt.Access)

		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "unauthorized"})
			return
		}
		id := token.Subject
		c.Set("x-user-id", id)
		c.Next()
	}
}

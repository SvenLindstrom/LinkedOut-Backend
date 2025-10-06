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
			c.JSON(http.StatusUnauthorized, gin.H{"message": "unauthorized"})
		}

		reqToken := strings.Split(bearerToken, " ")[1]

		token, err := jwt.Verify(reqToken, jwt.Access)

		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "unauthorized"})
		}

		id := token.Subject
		c.Set("x-user-id", id)
		c.Next()
	}
}

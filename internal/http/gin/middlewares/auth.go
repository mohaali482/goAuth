package middlewares

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mohaali482/goAuth/auth"
)

func AuthMiddleware(s auth.UserService) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString, err := c.Cookie("access_token")
		if err != nil || tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "request does not contain an access token"})
			c.Abort()
			return
		}

		_, err = s.ValidateJWT(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid access token"})
			c.Abort()
			return
		}

		c.Next()

	}

}

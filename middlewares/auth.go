package middlewares

import (
	"doctorrank_go/helpers"
	"doctorrank_go/responses"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

func Authentication() gin.HandlerFunc {
	return func(c *gin.Context) {
		auth := c.Request.Header.Get("Authorization")
		if auth == "" {
			c.JSON(http.StatusUnauthorized, responses.UserResponse{Status: http.StatusUnauthorized, Message: "error", Data: "No Authorization header provided"})
			c.Abort()
			return
		}
		token := strings.TrimPrefix(auth, "Bearer ")
		if token == auth {
			c.JSON(http.StatusUnauthorized, responses.UserResponse{Status: http.StatusUnauthorized, Message: "error", Data: "Could not find bearer token in Authorization header"})
			c.Abort()
			return
		}
		claims, err := helpers.ValidateToken(token)
		if err != "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err})
			c.Abort()
			return
		}
		fmt.Println(claims)
		c.Set("email", claims.Email)
		c.Set("first_name", claims.FirstName)
		c.Set("last_name", claims.LastName)
		c.Set("_id", claims.Id)

		c.Next()
	}
}

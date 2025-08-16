package middlewares

import (
	"context"
	"fmt"

	"Taskly.com/m/package/utils/auth"

	"github.com/gin-gonic/gin"
)

func AuthenMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Thay vì đọc từ header, đọc từ cookie
		jwtToken, err := c.Cookie("token")
		if err != nil {
			c.AbortWithStatusJSON(401, gin.H{"code": 400011, "eror": "Unauthorized", "description": "token not found in cookie"})
			fmt.Println("error", err)
			return
		}

		// validate jwt token by subject
		claims, err := auth.VerifyTokenSubject(jwtToken)
		if err != nil {
			c.AbortWithStatusJSON(401, gin.H{"code": 400012, "error": "InvalidToken", "description": ""})
			fmt.Println("error", err)
			return
		}
		// update claims to context
		ctx := context.WithValue(c.Request.Context(), "claims", claims)
		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}

}

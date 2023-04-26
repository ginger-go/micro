package auth

import (
	"strings"

	"github.com/gin-gonic/gin"
)

func GetJWTTokenFromHeader(c *gin.Context) string {
	token := c.Request.Header.Get("Authorization")
	token = strings.ReplaceAll(token, "Bearer ", "")
	token = strings.ReplaceAll(token, " ", "")
	return token
}

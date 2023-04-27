package auth

import (
	"strings"

	"github.com/gin-gonic/gin"
)

// GetAuthToken returns the token from the Authorization header: (Bearer <token>)
func GetAuthToken(ctx *gin.Context) string {
	t := ctx.GetHeader("Authorization")
	t = strings.ReplaceAll(t, "Bearer ", "")
	t = strings.ReplaceAll(t, " ", "")
	return t
}

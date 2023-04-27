package auth

import (
	"github.com/gin-gonic/gin"
	"github.com/ginger-go/micro"
)

func updatePublicPemHandler(ctx *gin.Context) {
	req := micro.GinRequest[AuthPublicPem](ctx)
	SYSTEM_TOKEN_PUBLIC_PEM = req.SystemPem
	USER_TOKEN_PUBLIC_PEM = req.UserPem
	ctx.JSON(200, nil)
}

func getSystemInfoHandler(ctx *gin.Context) {
	ctx.JSON(200, &SystemInfo{
		UUID: SYSTEM_ID,
		Name: SYSTEM_NAME,
	})
}

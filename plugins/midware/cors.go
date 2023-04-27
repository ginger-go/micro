package midware

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func CORS(config cors.Config) gin.HandlerFunc {
	return cors.New(config)
}

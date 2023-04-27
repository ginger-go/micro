package midware

import (
	"time"

	"github.com/gin-gonic/gin"
	limiter "github.com/ulule/limiter/v3"
	_gin "github.com/ulule/limiter/v3/drivers/middleware/gin"
	"github.com/ulule/limiter/v3/drivers/store/memory"
)

func RateLimited(duration time.Duration, rate int64) gin.HandlerFunc {
	return _gin.NewMiddleware(limiter.New(memory.NewStore(), limiter.Rate{Period: duration, Limit: rate}))
}

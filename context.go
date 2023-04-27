package micro

import (
	"net/http/httptest"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/ginger-go/sql"
)

type Trace struct {
	Success    bool          `json:"success"`
	Time       time.Time     `json:"time"`
	SystemID   string        `json:"system_id"`
	SystemName string        `json:"system_name"`
	TraceID    string        `json:"trace_id"`
	Error      ResponseError `json:"error"`
}

type Context[T any] struct {
	GinContext *gin.Context
	TraceID    string
	Request    *T
	Page       *sql.Pagination
	Sort       *sql.Sort
	Response   interface{}
}

type MockContextParams[T any] struct {
	Request   *T
	Page      *sql.Pagination
	Sort      *sql.Sort
	Method    string
	Path      string
	ClientIP  string
	UserAgent string
	Headers   map[string]string
}

func NewMockContext[T any](param MockContextParams[T]) *Context[T] {
	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)

	if param.Method == "" {
		param.Method = "GET"
	}
	if param.Path == "" {
		param.Path = "/"
	}
	ctx.Request = httptest.NewRequest(param.Method, param.Path, nil)
	if param.ClientIP != "" {
		ctx.Request.RemoteAddr = param.ClientIP
	}

	if param.UserAgent != "" {
		ctx.Request.Header.Set("User-Agent", param.UserAgent)
	} else {
		ctx.Request.Header.Set("User-Agent", "GingerMockContext")
	}
	if param.Headers != nil {
		for k, v := range param.Headers {
			ctx.Request.Header.Set(k, v)
		}
	}

	return &Context[T]{
		GinContext: ctx,
		Request:    param.Request,
		Page:       param.Page,
		Sort:       param.Sort,
	}
}

func (ctx *Context[T]) ClientIP() string {
	return ctx.GinContext.ClientIP()
}

func (ctx *Context[T]) UserAgent() string {
	return ctx.GinContext.Request.UserAgent()
}

func (ctx *Context[T]) OK(data interface{}, traceID string, traces []Trace, page ...*sql.Pagination) {
	var p *sql.Pagination
	if len(page) > 0 {
		p = page[0]
	}
	resp := &Response{
		Success:    true,
		Data:       data,
		Pagination: p,
		TraceID:    traceID,
		Traces:     traces,
	}
	ctx.Response = resp // for testing
	ctx.GinContext.JSON(200, resp)
}

func (ctx *Context[T]) Error(err Error, traceID string, traces []Trace) {
	resp := &Response{
		Success: false,
		Error: &ResponseError{
			Code:    err.Code(),
			Message: err.Error(),
		},
		TraceID: traceID,
		Traces:  traces,
	}
	ctx.Response = resp // for testing
	ctx.GinContext.JSON(200, resp)
}

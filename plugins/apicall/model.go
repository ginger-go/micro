package apicall

import (
	"github.com/ginger-go/micro"
	"github.com/ginger-go/sql"
)

type Response[T any] struct {
	Success    bool                 `json:"success"`
	Error      *micro.ResponseError `json:"error,omitempty"`
	Pagination *sql.Pagination      `json:"pagination,omitempty"`
	Data       *T                   `json:"data,omitempty"`
	Traces     []micro.Trace        `json:"traces,omitempty"`
}

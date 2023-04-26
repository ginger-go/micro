package caller

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/ginger-go/micro"
	"github.com/ginger-go/sql"
)

type Response[T any] struct {
	Success    bool            `json:"success"`
	Error      *ResponseError  `json:"error,omitempty"`
	Pagination *sql.Pagination `json:"pagination,omitempty"`
	Data       *T              `json:"data,omitempty"`
	Traces     []micro.Trace   `json:"traces,omitempty"`
}

type ResponseError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func GET[T any](url string, params map[string]string, headers map[string]string, traceID string, traces []micro.Trace) (*Response[T], error) {

	return get[T](url, params, headers, traceID, traces)
}

func POST[T any](url string, body interface{}, headers map[string]string, traceID string, traces []micro.Trace) (*Response[T], error) {
	return nonGet[T](url, "POST", body, headers, traceID, traces)
}

func PUT[T any](url string, body interface{}, headers map[string]string, traceID string, traces []micro.Trace) (*Response[T], error) {
	return nonGet[T](url, "PUT", body, headers, traceID, traces)
}

func DELETE[T any](url string, body interface{}, headers map[string]string, traceID string, traces []micro.Trace) (*Response[T], error) {
	return nonGet[T](url, "DELETE", body, headers, traceID, traces)
}

func get[T any](url string, params map[string]string, headers map[string]string, traceID string, traces []micro.Trace) (*Response[T], error) {
	if headers == nil {
		headers = make(map[string]string)
	}
	headers["Micro-TraceID"] = traceID
	tracesStr, _ := json.Marshal(traces)
	headers["Micro-Traces"] = string(tracesStr)
	url += "?"
	for k, v := range params {
		url += k + "=" + v + "&"
	}
	url = url[:len(url)-1]

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	for k, v := range headers {
		req.Header.Set(k, v)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	var response Response[T]
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

func nonGet[T any](url string, method string, body interface{}, headers map[string]string, traceID string, traces []micro.Trace) (*Response[T], error) {
	if headers == nil {
		headers = make(map[string]string)
	}
	headers["Micro-TraceID"] = traceID
	tracesStr, _ := json.Marshal(traces)
	headers["Micro-Traces"] = string(tracesStr)
	b, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(method, url, bytes.NewReader(b))
	if err != nil {
		return nil, err
	}

	for k, v := range headers {
		req.Header.Set(k, v)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	var response Response[T]
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

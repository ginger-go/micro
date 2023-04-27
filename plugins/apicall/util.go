package apicall

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/ginger-go/micro"
)

func get[T any](url string, params map[string]string, headers map[string]string, traceID string, traces []micro.Trace) (*Response[T], error) {
	if headers == nil {
		headers = make(map[string]string)
	}
	headers[micro.MICRO_HEADER_TRACE_ID] = traceID
	tracesStr, _ := json.Marshal(traces)
	headers[micro.MICRO_HEADER_TRACES] = string(tracesStr)
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
	headers[micro.MICRO_HEADER_TRACE_ID] = traceID
	tracesStr, _ := json.Marshal(traces)
	headers[micro.MICRO_HEADER_TRACES] = string(tracesStr)
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

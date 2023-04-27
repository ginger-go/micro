package apicall

import "github.com/ginger-go/micro"

// GET make a GET request with micro trace standard
func GET[T any](url string, params map[string]string, headers map[string]string, traceID string, traces []micro.Trace) (*Response[T], error) {
	return get[T](url, params, headers, traceID, traces)
}

// POST make a POST request with micro trace standard
func POST[T any](url string, body interface{}, headers map[string]string, traceID string, traces []micro.Trace) (*Response[T], error) {
	return nonGet[T](url, "POST", body, headers, traceID, traces)
}

// PUT make a PUT request with micro trace standard
func PUT[T any](url string, body interface{}, headers map[string]string, traceID string, traces []micro.Trace) (*Response[T], error) {
	return nonGet[T](url, "PUT", body, headers, traceID, traces)
}

// DELETE make a DELETE request with micro trace standard
func DELETE[T any](url string, body interface{}, headers map[string]string, traceID string, traces []micro.Trace) (*Response[T], error) {
	return nonGet[T](url, "DELETE", body, headers, traceID, traces)
}

package incidentio

import (
	"context"
	"encoding/json"
	"net/http"
	"net/url"
)

// send builds and executes a request, decoding the JSON response into v when v is not nil.
func (c *Client) send(ctx context.Context, method, path string, query url.Values, body, v any) (*http.Response, error) {
	req, err := c.NewRequest(method, path, body)
	if err != nil {
		return nil, err
	}
	if len(query) > 0 {
		req.URL.RawQuery = query.Encode()
	}
	return c.Do(ctx, req, v)
}

// getKey executes a request and decodes the value stored under key in the
// response envelope, for example {"incident": {...}}. A missing key, or an
// empty response, yields the zero value.
func getKey[T any](ctx context.Context, c *Client, method, path string, query url.Values, body any, key string) (T, *http.Response, error) {
	var zero T

	var envelope map[string]json.RawMessage
	resp, err := c.send(ctx, method, path, query, body, &envelope)
	if err != nil {
		return zero, resp, err
	}

	raw, ok := envelope[key]
	if !ok {
		return zero, resp, nil
	}

	var out T
	if err := json.Unmarshal(raw, &out); err != nil {
		return zero, resp, err
	}
	return out, resp, nil
}

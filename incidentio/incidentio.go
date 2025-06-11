package incidentio

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

const (
	defaultBaseURL = "https://api.incident.io/"
	userAgent      = "incidentio-go-client/1.0.0"
)

// Client manages communication with the Incident.io API.
type Client struct {
	client    *http.Client
	BaseURL   *url.URL
	UserAgent string
	apiKey    string

	// Services used for talking to different parts of the Incident.io API.
	Incidents     *IncidentsService
	Severities    *SeveritiesService
	IncidentTypes *IncidentTypesService
	IncidentRoles *IncidentRolesService
	CustomFields  *CustomFieldsService
	Actions       *ActionsService
	Workflows     *WorkflowsService
	Schedules     *SchedulesService
	Users         *UsersService
	Webhooks      *WebhooksService
}

// ClientOption allows for functional options to configure the client.
type ClientOption func(*Client)

// WithHTTPClient sets a custom HTTP client.
func WithHTTPClient(httpClient *http.Client) ClientOption {
	return func(c *Client) {
		c.client = httpClient
	}
}

// WithBaseURL sets a custom base URL.
func WithBaseURL(baseURL string) ClientOption {
	return func(c *Client) {
		u, _ := url.Parse(baseURL)
		c.BaseURL = u
	}
}

// NewClient returns a new Incident.io API client.
func NewClient(apiKey string, opts ...ClientOption) *Client {
	baseURL, _ := url.Parse(defaultBaseURL)

	c := &Client{
		client:    &http.Client{Timeout: 30 * time.Second},
		BaseURL:   baseURL,
		UserAgent: userAgent,
		apiKey:    apiKey,
	}

	// Apply options
	for _, opt := range opts {
		opt(c)
	}

	// Initialize services
	c.Incidents = &IncidentsService{client: c}
	c.Severities = &SeveritiesService{client: c}
	c.IncidentTypes = &IncidentTypesService{client: c}
	c.IncidentRoles = &IncidentRolesService{client: c}
	c.CustomFields = &CustomFieldsService{client: c}
	c.Actions = &ActionsService{client: c}
	c.Workflows = &WorkflowsService{client: c}
	c.Schedules = &SchedulesService{client: c}
	c.Users = &UsersService{client: c}
	c.Webhooks = &WebhooksService{client: c}

	return c
}

// NewRequest creates an API request.
func (c *Client) NewRequest(method, urlStr string, body interface{}) (*http.Request, error) {
	u, err := c.BaseURL.Parse(urlStr)
	if err != nil {
		return nil, err
	}

	var buf io.ReadWriter
	if body != nil {
		buf = &bytes.Buffer{}
		enc := json.NewEncoder(buf)
		enc.SetEscapeHTML(false)
		err := enc.Encode(body)
		if err != nil {
			return nil, err
		}
	}

	req, err := http.NewRequest(method, u.String(), buf)
	if err != nil {
		return nil, err
	}

	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", c.UserAgent)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.apiKey))

	return req, nil
}

// Do sends an API request and returns the API response.
func (c *Client) Do(ctx context.Context, req *http.Request, v interface{}) (*http.Response, error) {
	req = req.WithContext(ctx)

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	err = CheckResponse(resp)
	if err != nil {
		return resp, err
	}

	if v != nil && resp.StatusCode != http.StatusNoContent {
		if w, ok := v.(io.Writer); ok {
			io.Copy(w, resp.Body)
		} else {
			err = json.NewDecoder(resp.Body).Decode(v)
		}
	}

	return resp, err
}

// CheckResponse checks the API response for errors.
func CheckResponse(r *http.Response) error {
	if c := r.StatusCode; 200 <= c && c <= 299 {
		return nil
	}

	errorResponse := &ErrorResponse{Response: r}
	data, err := io.ReadAll(r.Body)
	if err == nil && data != nil {
		json.Unmarshal(data, errorResponse)
	}

	return errorResponse
}

// ErrorResponse represents an error response from the API.
type ErrorResponse struct {
	Response *http.Response
	Type     string `json:"type"`
	Status   int    `json:"status"`
	Detail   string `json:"detail"`
	Errors   []struct {
		Code   string `json:"code"`
		Detail string `json:"detail"`
		Source struct {
			Pointer string `json:"pointer"`
		} `json:"source"`
	} `json:"errors"`
}

func (r *ErrorResponse) Error() string {
	return fmt.Sprintf("%v %v: %d %v",
		r.Response.Request.Method, r.Response.Request.URL,
		r.Response.StatusCode, r.Detail)
}

// ListOptions specifies the optional parameters to various List methods.
type ListOptions struct {
	PageSize int    `url:"page_size,omitempty"`
	After    string `url:"after,omitempty"`
}

// Common types used across the API

type ExternalResource struct {
	ExternalID  string `json:"external_id"`
	DisplayName string `json:"display_name"`
}

type Timestamp struct {
	time.Time
}

func (t *Timestamp) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	parsedTime, err := time.Parse(time.RFC3339, s)
	if err != nil {
		return err
	}
	t.Time = parsedTime
	return nil
}

func (t Timestamp) MarshalJSON() ([]byte, error) {
	return json.Marshal(t.Format(time.RFC3339))
}

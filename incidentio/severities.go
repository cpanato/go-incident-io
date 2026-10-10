package incidentio

import (
	"context"
	"net/http"
)

// SeveritiesService handles communication with the severity related methods.
type SeveritiesService struct {
	client *Client
}

// List returns a list of severities.
func (s *SeveritiesService) List(ctx context.Context) ([]*Severity, *http.Response, error) {
	return getKey[[]*Severity](ctx, s.client, "GET", "v1/severities", nil, nil, "severities")
}

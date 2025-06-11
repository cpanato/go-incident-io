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
	u := "v1/severities"

	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	var result struct {
		Severities []*Severity `json:"severities"`
	}
	resp, err := s.client.Do(ctx, req, &result)
	if err != nil {
		return nil, resp, err
	}

	return result.Severities, resp, nil
}

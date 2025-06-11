package incidentio

import (
	"context"
	"net/http"
)

// IncidentTypesService handles communication with the incident type related methods.
type IncidentTypesService struct {
	client *Client
}

// IncidentType represents an incident type in Incident.io.
type IncidentType struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	CreatedAt   Timestamp `json:"created_at"`
	UpdatedAt   Timestamp `json:"updated_at"`
}

// List returns a list of incident types.
func (s *IncidentTypesService) List(ctx context.Context) ([]*IncidentType, *http.Response, error) {
	u := "v1/incident_types"

	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	var result struct {
		IncidentTypes []*IncidentType `json:"incident_types"`
	}
	resp, err := s.client.Do(ctx, req, &result)
	if err != nil {
		return nil, resp, err
	}

	return result.IncidentTypes, resp, nil
}

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
	return getKey[[]*IncidentType](ctx, s.client, "GET", "v1/incident_types", nil, nil, "incident_types")
}

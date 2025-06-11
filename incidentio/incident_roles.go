package incidentio

import (
	"context"
	"net/http"
)

type CreateRoleAssignment struct {
	IncidentRoleID string `json:"incident_role_id"`
	UserID         string `json:"user_id"`
}

// IncidentRoleAssignment represents the assignment of a role to a user in an incident.
type IncidentRoleAssignment struct {
	Role     *IncidentRole `json:"role"`
	Assignee *User         `json:"assignee"`
}

type Severity struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Rank        int       `json:"rank"`
	CreatedAt   Timestamp `json:"created_at"`
	UpdatedAt   Timestamp `json:"updated_at"`
}

type IncidentRole struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Required    bool      `json:"required"`
	CreatedAt   Timestamp `json:"created_at"`
	UpdatedAt   Timestamp `json:"updated_at"`
}

// IncidentRolesService handles communication with the incident role related methods.
type IncidentRolesService struct {
	client *Client
}

// List returns a list of incident roles.
func (s *IncidentRolesService) List(ctx context.Context) ([]*IncidentRole, *http.Response, error) {
	u := "v2/incident_roles"

	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	var result struct {
		IncidentRoles []*IncidentRole `json:"incident_roles"`
	}
	resp, err := s.client.Do(ctx, req, &result)
	if err != nil {
		return nil, resp, err
	}

	return result.IncidentRoles, resp, nil
}

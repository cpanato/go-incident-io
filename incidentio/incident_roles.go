package incidentio

import (
	"context"
	"net/http"
)

// CreateRoleAssignment represents the payload for creating a role assignment in an incident.
type CreateRoleAssignment struct {
	IncidentRoleID string        `json:"incident_role_id"`
	Assignee       UserReference `json:"assignee"`
}

// IncidentRoleAssignment represents the assignment of a role to a user in an incident.
type IncidentRoleAssignment struct {
	Role     *IncidentRole `json:"role"`
	Assignee *User         `json:"assignee"`
}

// Severity represents the severity of an incident.
type Severity struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Rank        int       `json:"rank"`
	CreatedAt   Timestamp `json:"created_at"`
	UpdatedAt   Timestamp `json:"updated_at"`
}

// IncidentRole represents a role that can be assigned to users in an incident.
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
	return getKey[[]*IncidentRole](ctx, s.client, "GET", "v2/incident_roles", nil, nil, "incident_roles")
}

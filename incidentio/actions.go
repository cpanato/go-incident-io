package incidentio

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
)

// ActionsService handles communication with the actions related methods.
type ActionsService struct {
	client *Client
}

// Action represents an action in Incident.io.
type Action struct {
	ID          string     `json:"id"`
	IncidentID  string     `json:"incident_id"`
	Description string     `json:"description"`
	Status      string     `json:"status"`
	Assignee    *User      `json:"assignee,omitempty"`
	Creator     *Actor     `json:"creator,omitempty"`
	CompletedAt *Timestamp `json:"completed_at,omitempty"`
	CreatedAt   Timestamp  `json:"created_at"`
	UpdatedAt   Timestamp  `json:"updated_at"`
}

// ActionListOptions represents the options for listing actions.
type ActionListOptions struct {
	ListOptions
	IncidentID string
	// IncidentMode limits actions to incidents of this mode. By default only
	// standard and retrospective incidents are included.
	IncidentMode string
	CreatedAt    Filter
	UpdatedAt    Filter
}

// CreateActionOptions represents the options for creating an action.
type CreateActionOptions struct {
	IncidentID  string `json:"incident_id"`
	Description string `json:"description"`
	AssigneeID  string `json:"assignee_id,omitempty"`
}

// UpdateActionOptions represents the options for updating an action.
type UpdateActionOptions struct {
	Description string `json:"description"`
	Status      string `json:"status"`
	AssigneeID  string `json:"assignee_id,omitempty"`
}

// List returns a list of actions.
func (s *ActionsService) List(ctx context.Context, opts *ActionListOptions) ([]*Action, *http.Response, error) {
	q := url.Values{}
	if opts != nil {
		opts.apply(q)
		if opts.IncidentID != "" {
			q.Set("incident_id", opts.IncidentID)
		}
		if opts.IncidentMode != "" {
			q.Set("incident_mode", opts.IncidentMode)
		}
		addFilter(q, "created_at", opts.CreatedAt)
		addFilter(q, "updated_at", opts.UpdatedAt)
	}

	return getKey[[]*Action](ctx, s.client, "GET", "v3/actions", q, nil, "actions")
}

// Get returns a single action.
func (s *ActionsService) Get(ctx context.Context, id string) (*Action, *http.Response, error) {
	return getKey[*Action](ctx, s.client, "GET", fmt.Sprintf("v3/actions/%s", id), nil, nil, "action")
}

// Create creates a new action.
func (s *ActionsService) Create(ctx context.Context, opts *CreateActionOptions) (*Action, *http.Response, error) {
	return getKey[*Action](ctx, s.client, "POST", "v3/actions", nil, opts, "action")
}

// Update updates an action.
func (s *ActionsService) Update(ctx context.Context, id string, opts *UpdateActionOptions) (*Action, *http.Response, error) {
	return getKey[*Action](ctx, s.client, "PUT", fmt.Sprintf("v3/actions/%s", id), nil, opts, "action")
}

// Delete deletes an action.
func (s *ActionsService) Delete(ctx context.Context, id string) (*http.Response, error) {
	return s.client.send(ctx, "DELETE", fmt.Sprintf("v3/actions/%s", id), nil, nil, nil)
}

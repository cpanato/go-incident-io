package incidentio

import (
	"context"
	"fmt"
	"net/http"
)

// IncidentsService handles communication with the incident related methods.
type IncidentsService struct {
	client *Client
}

// Incident represents an incident in Incident.io.
type Incident struct {
	ID                      string                   `json:"id"`
	Name                    string                   `json:"name"`
	Summary                 string                   `json:"summary,omitempty"`
	Type                    string                   `json:"type"`
	Status                  string                   `json:"status"`
	Severity                *Severity                `json:"severity,omitempty"`
	IncidentRoleAssignments []IncidentRoleAssignment `json:"incident_role_assignments,omitempty"`
	CustomFieldValues       map[string]interface{}   `json:"custom_field_values,omitempty"`
	CreatedAt               Timestamp                `json:"created_at"`
	UpdatedAt               Timestamp                `json:"updated_at"`
	ReportedAt              *Timestamp               `json:"reported_at,omitempty"`
	ClosedAt                *Timestamp               `json:"closed_at,omitempty"`
	LastActivityAt          *Timestamp               `json:"last_activity_at,omitempty"`
	Mode                    string                   `json:"mode"`
	Visibility              string                   `json:"visibility"`
	SlackChannelID          string                   `json:"slack_channel_id,omitempty"`
	SlackChannelName        string                   `json:"slack_channel_name,omitempty"`
	Creator                 *User                    `json:"creator,omitempty"`
	ExternalIssueReference  *ExternalResource        `json:"external_issue_reference,omitempty"`
	PostmortemDocumentURL   string                   `json:"postmortem_document_url,omitempty"`
	SlackThreadURL          string                   `json:"slack_thread_url,omitempty"`
	CallURL                 string                   `json:"call_url,omitempty"`
}

// CreateIncidentOptions represents the options for creating an incident.
type CreateIncidentOptions struct {
	Name                     string                 `json:"name"`
	Summary                  string                 `json:"summary,omitempty"`
	IncidentTypeID           string                 `json:"incident_type_id"`
	SeverityID               string                 `json:"severity_id,omitempty"`
	IncidentRoleAssignments  []CreateRoleAssignment `json:"incident_role_assignments,omitempty"`
	CustomFieldValues        map[string]interface{} `json:"custom_field_values,omitempty"`
	Mode                     string                 `json:"mode,omitempty"`
	Visibility               string                 `json:"visibility,omitempty"`
	SlackChannelNameOverride string                 `json:"slack_channel_name_override,omitempty"`
}

// List returns a list of incidents.
func (s *IncidentsService) List(ctx context.Context, _ *ListOptions) ([]*Incident, *http.Response, error) {
	u := "v2/incidents"

	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	var result struct {
		Incidents []*Incident `json:"incidents"`
	}
	resp, err := s.client.Do(ctx, req, &result)
	if err != nil {
		return nil, resp, err
	}

	return result.Incidents, resp, nil
}

// Get returns a single incident.
func (s *IncidentsService) Get(ctx context.Context, id string) (*Incident, *http.Response, error) {
	u := fmt.Sprintf("v2/incidents/%s", id)

	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	var result struct {
		Incident *Incident `json:"incident"`
	}
	resp, err := s.client.Do(ctx, req, &result)
	if err != nil {
		return nil, resp, err
	}

	return result.Incident, resp, nil
}

// Create creates a new incident.
func (s *IncidentsService) Create(ctx context.Context, opts *CreateIncidentOptions) (*Incident, *http.Response, error) {
	u := "v2/incidents"

	req, err := s.client.NewRequest("POST", u, opts)
	if err != nil {
		return nil, nil, err
	}

	var result struct {
		Incident *Incident `json:"incident"`
	}
	resp, err := s.client.Do(ctx, req, &result)
	if err != nil {
		return nil, resp, err
	}

	return result.Incident, resp, nil
}

// UpdateIncidentOptions represents the options for updating an incident.
type UpdateIncidentOptions struct {
	Name                    *string                `json:"name,omitempty"`
	Summary                 *string                `json:"summary,omitempty"`
	Status                  *string                `json:"status,omitempty"`
	SeverityID              *string                `json:"severity_id,omitempty"`
	IncidentRoleAssignments []CreateRoleAssignment `json:"incident_role_assignments,omitempty"`
	CustomFieldValues       map[string]interface{} `json:"custom_field_values,omitempty"`
}

// Update updates an incident.
func (s *IncidentsService) Update(ctx context.Context, id string, opts *UpdateIncidentOptions) (*Incident, *http.Response, error) {
	u := fmt.Sprintf("v2/incidents/%s", id)

	req, err := s.client.NewRequest("PUT", u, opts)
	if err != nil {
		return nil, nil, err
	}

	var result struct {
		Incident *Incident `json:"incident"`
	}
	resp, err := s.client.Do(ctx, req, &result)
	if err != nil {
		return nil, resp, err
	}

	return result.Incident, resp, nil
}

// Delete deletes an incident.
func (s *IncidentsService) Delete(ctx context.Context, id string) (*http.Response, error) {
	u := fmt.Sprintf("v2/incidents/%s", id)

	req, err := s.client.NewRequest("DELETE", u, nil)
	if err != nil {
		return nil, err
	}

	return s.client.Do(ctx, req, nil)
}

package incidentio

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"net/http"
	"net/url"
)

// IncidentsService handles communication with the incident related methods.
type IncidentsService struct {
	client *Client
}

// IncidentStatus represents the status of an incident.
type IncidentStatus struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Category    string    `json:"category"`
	Rank        int       `json:"rank"`
	CreatedAt   Timestamp `json:"created_at"`
	UpdatedAt   Timestamp `json:"updated_at"`
}

// Actor represents who or what performed an action. Only one field is set.
type Actor struct {
	User     *User          `json:"user,omitempty"`
	APIKey   *APIKeyActor   `json:"api_key,omitempty"`
	Workflow *WorkflowActor `json:"workflow,omitempty"`
	Alert    *AlertActor    `json:"alert,omitempty"`
}

// APIKeyActor represents an API key acting on an incident.
type APIKeyActor struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// WorkflowActor represents a workflow acting on an incident.
type WorkflowActor struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// AlertActor represents an alert acting on an incident.
type AlertActor struct {
	ID    string `json:"id"`
	Title string `json:"title"`
}

// CustomFieldValue is a single value of a custom field entry.
type CustomFieldValue struct {
	ValueText         string             `json:"value_text,omitempty"`
	ValueLink         string             `json:"value_link,omitempty"`
	ValueNumeric      string             `json:"value_numeric,omitempty"`
	ValueOption       *CustomFieldOption `json:"value_option,omitempty"`
	ValueCatalogEntry *CatalogEntryRef   `json:"value_catalog_entry,omitempty"`
}

// CustomFieldOption is an option of a select-type custom field.
type CustomFieldOption struct {
	ID            string `json:"id"`
	CustomFieldID string `json:"custom_field_id"`
	Value         string `json:"value"`
	SortKey       int    `json:"sort_key"`
}

// CatalogEntryRef is a reference to a catalog entry.
type CatalogEntryRef struct {
	ID         string   `json:"id"`
	Name       string   `json:"name"`
	ExternalID string   `json:"external_id,omitempty"`
	Aliases    []string `json:"aliases,omitempty"`
}

// CustomFieldEntry is the set of values an incident has for a custom field.
type CustomFieldEntry struct {
	CustomField *CustomField       `json:"custom_field"`
	Values      []CustomFieldValue `json:"values"`
}

// ExternalIssueReference links an incident to an issue in an external tracker.
type ExternalIssueReference struct {
	IssueName      string `json:"issue_name"`
	IssuePermalink string `json:"issue_permalink"`
	Provider       string `json:"provider"`
}

// Incident represents an incident in Incident.io.
type Incident struct {
	ID                      string                   `json:"id"`
	Reference               string                   `json:"reference"`
	Name                    string                   `json:"name"`
	Summary                 string                   `json:"summary,omitempty"`
	Permalink               string                   `json:"permalink,omitempty"`
	IncidentStatus          *IncidentStatus          `json:"incident_status,omitempty"`
	IncidentType            *IncidentType            `json:"incident_type,omitempty"`
	Severity                *Severity                `json:"severity,omitempty"`
	IncidentRoleAssignments []IncidentRoleAssignment `json:"incident_role_assignments,omitempty"`
	CustomFieldEntries      []CustomFieldEntry       `json:"custom_field_entries,omitempty"`
	TeamIDs                 []string                 `json:"team_ids,omitempty"`
	CreatedAt               Timestamp                `json:"created_at"`
	UpdatedAt               Timestamp                `json:"updated_at"`
	LastActivityAt          Timestamp                `json:"last_activity_at"`
	Mode                    string                   `json:"mode"`
	Visibility              string                   `json:"visibility"`
	SlackTeamID             string                   `json:"slack_team_id,omitempty"`
	SlackChannelID          string                   `json:"slack_channel_id,omitempty"`
	SlackChannelName        string                   `json:"slack_channel_name,omitempty"`
	SlackChannelURL         string                   `json:"slack_channel_url,omitempty"`
	Creator                 *Actor                   `json:"creator,omitempty"`
	ExternalIssueReference  *ExternalIssueReference  `json:"external_issue_reference,omitempty"`
	PostmortemDocumentURL   string                   `json:"postmortem_document_url,omitempty"`
	CallURL                 string                   `json:"call_url,omitempty"`
}

// CustomFieldValuePayload is a single value to set on a custom field.
// Exactly one of the value fields should be set.
type CustomFieldValuePayload struct {
	ValueCatalogEntryID string `json:"value_catalog_entry_id,omitempty"`
	ValueLink           string `json:"value_link,omitempty"`
	ValueNumeric        string `json:"value_numeric,omitempty"`
	ValueOptionID       string `json:"value_option_id,omitempty"`
	ValueText           string `json:"value_text,omitempty"`
}

// CustomFieldEntryPayload sets the values of a custom field on an incident.
type CustomFieldEntryPayload struct {
	CustomFieldID string                    `json:"custom_field_id"`
	Values        []CustomFieldValuePayload `json:"values"`
}

// IncidentTimestampValuePayload sets the value of an incident timestamp.
type IncidentTimestampValuePayload struct {
	IncidentTimestampID string `json:"incident_timestamp_id"`
	Value               string `json:"value,omitempty"`
}

// RetrospectiveIncidentOptions are required when creating a retrospective incident.
type RetrospectiveIncidentOptions struct {
	ExternalID            int    `json:"external_id,omitempty"`
	PostmortemDocumentURL string `json:"postmortem_document_url,omitempty"`
	SlackChannelID        string `json:"slack_channel_id,omitempty"`
}

// CreateIncidentOptions represents the options for creating an incident.
type CreateIncidentOptions struct {
	// IdempotencyKey makes creation safe to retry. A random key is generated when empty.
	IdempotencyKey               string                          `json:"idempotency_key"`
	Visibility                   string                          `json:"visibility"`
	Name                         string                          `json:"name,omitempty"`
	Summary                      string                          `json:"summary,omitempty"`
	IncidentTypeID               string                          `json:"incident_type_id,omitempty"`
	IncidentStatusID             string                          `json:"incident_status_id,omitempty"`
	SeverityID                   string                          `json:"severity_id,omitempty"`
	Mode                         string                          `json:"mode,omitempty"`
	SlackChannelNameOverride     string                          `json:"slack_channel_name_override,omitempty"`
	SlackTeamID                  string                          `json:"slack_team_id,omitempty"`
	IncidentRoleAssignments      []CreateRoleAssignment          `json:"incident_role_assignments,omitempty"`
	CustomFieldEntries           []CustomFieldEntryPayload       `json:"custom_field_entries,omitempty"`
	IncidentTimestampValues      []IncidentTimestampValuePayload `json:"incident_timestamp_values,omitempty"`
	RetrospectiveIncidentOptions *RetrospectiveIncidentOptions   `json:"retrospective_incident_options,omitempty"`
}

// Filter maps a filter operator (for example "one_of" or "gte") to its values.
type Filter map[string][]string

// IncidentListOptions represents the options for listing incidents.
type IncidentListOptions struct {
	ListOptions
	// SortBy is "created_at_newest_first" (default) or "created_at_oldest_first".
	SortBy string
	// FilterMode is "all" (default) or "any".
	FilterMode string

	Status         Filter
	StatusCategory Filter
	CreatedAt      Filter
	UpdatedAt      Filter
	Severity       Filter
	IncidentType   Filter
	Mode           Filter

	// IncidentRole and CustomField filter by role or custom field ID, then operator.
	IncidentRole map[string]Filter
	CustomField  map[string]Filter
}

func (o *IncidentListOptions) apply(q url.Values) {
	if o == nil {
		return
	}
	o.ListOptions.apply(q)
	if o.SortBy != "" {
		q.Set("sort_by", o.SortBy)
	}
	if o.FilterMode != "" {
		q.Set("filter_mode", o.FilterMode)
	}

	for name, f := range map[string]Filter{
		"status":          o.Status,
		"status_category": o.StatusCategory,
		"created_at":      o.CreatedAt,
		"updated_at":      o.UpdatedAt,
		"severity":        o.Severity,
		"incident_type":   o.IncidentType,
		"mode":            o.Mode,
	} {
		addFilter(q, name, f)
	}
	for name, byID := range map[string]map[string]Filter{
		"incident_role": o.IncidentRole,
		"custom_field":  o.CustomField,
	} {
		for id, f := range byID {
			addFilter(q, fmt.Sprintf("%s[%s]", name, id), f)
		}
	}
}

func addFilter(q url.Values, key string, f Filter) {
	for op, values := range f {
		k := fmt.Sprintf("%s[%s]", key, op)
		for _, v := range values {
			q.Add(k, v)
		}
	}
}

// List returns a list of incidents.
func (s *IncidentsService) List(ctx context.Context, opts *IncidentListOptions) ([]*Incident, *http.Response, error) {
	u := "v2/incidents"

	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	if opts != nil {
		q := req.URL.Query()
		opts.apply(q)
		req.URL.RawQuery = q.Encode()
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

	body := opts
	if opts != nil && opts.IdempotencyKey == "" {
		key, err := newIdempotencyKey()
		if err != nil {
			return nil, nil, err
		}
		cp := *opts
		cp.IdempotencyKey = key
		body = &cp
	}

	req, err := s.client.NewRequest("POST", u, body)
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

func newIdempotencyKey() (string, error) {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

// UpdateIncidentOptions represents the options for editing an incident.
// Only fields that are set are sent.
type UpdateIncidentOptions struct {
	Name                     *string                         `json:"name,omitempty"`
	Summary                  *string                         `json:"summary,omitempty"`
	IncidentStatusID         *string                         `json:"incident_status_id,omitempty"`
	SeverityID               *string                         `json:"severity_id,omitempty"`
	CallURL                  *string                         `json:"call_url,omitempty"`
	SlackChannelNameOverride *string                         `json:"slack_channel_name_override,omitempty"`
	IncidentRoleAssignments  []CreateRoleAssignment          `json:"incident_role_assignments,omitempty"`
	CustomFieldEntries       []CustomFieldEntryPayload       `json:"custom_field_entries,omitempty"`
	IncidentTimestampValues  []IncidentTimestampValuePayload `json:"incident_timestamp_values,omitempty"`

	// NotifyIncidentChannel announces the edit in the incident channel.
	NotifyIncidentChannel bool `json:"-"`
}

// Update edits an incident.
func (s *IncidentsService) Update(ctx context.Context, id string, opts *UpdateIncidentOptions) (*Incident, *http.Response, error) {
	u := fmt.Sprintf("v2/incidents/%s/actions/edit", id)

	if opts == nil {
		opts = &UpdateIncidentOptions{}
	}
	body := struct {
		Incident              *UpdateIncidentOptions `json:"incident"`
		NotifyIncidentChannel bool                   `json:"notify_incident_channel"`
	}{opts, opts.NotifyIncidentChannel}

	req, err := s.client.NewRequest("POST", u, body)
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

// PostmortemDocument represents a postmortem document attached to an incident.
type PostmortemDocument struct {
	ID          string    `json:"id"`
	IncidentID  string    `json:"incident_id"`
	Title       string    `json:"title"`
	Status      string    `json:"status"`
	DocumentURL string    `json:"document_url"`
	Type        string    `json:"type"`
	CreatedAt   Timestamp `json:"created_at"`
	UpdatedAt   Timestamp `json:"updated_at"`
}

// ImportPostmortemDocument imports a postmortem document into an incident.
func (s *IncidentsService) ImportPostmortemDocument(ctx context.Context, id, title, content string) (*PostmortemDocument, *http.Response, error) {
	u := fmt.Sprintf("v2/incidents/%s/actions/import_postmortem_document", id)

	body := struct {
		Title   string `json:"title"`
		Content string `json:"content"`
	}{title, content}

	req, err := s.client.NewRequest("POST", u, body)
	if err != nil {
		return nil, nil, err
	}

	var result struct {
		PostmortemDocument *PostmortemDocument `json:"postmortem_document"`
	}
	resp, err := s.client.Do(ctx, req, &result)
	if err != nil {
		return nil, resp, err
	}

	return result.PostmortemDocument, resp, nil
}

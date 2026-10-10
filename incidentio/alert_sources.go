package incidentio

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

// AlertSourcesService handles communication with the alert sources related methods.
type AlertSourcesService struct {
	client *Client
}

// AlertSource represents a source of alerts, such as Datadog or a custom HTTP integration.
//
// Deeply nested, source-specific structures (template, condition groups and
// per-source options) are exposed as raw JSON so that callers can use the full
// API surface; see the Incident.io API reference for their shape.
type AlertSource struct {
	ID                        string   `json:"id"`
	Name                      string   `json:"name"`
	SourceType                string   `json:"source_type"`
	SecretToken               string   `json:"secret_token,omitempty"`
	AlertEventsURL            string   `json:"alert_events_url,omitempty"`
	AutoResolveIncidentAlerts bool     `json:"auto_resolve_incident_alerts,omitempty"`
	AutoResolveTimeoutMinutes int      `json:"auto_resolve_timeout_minutes,omitempty"`
	FixedTeamID               string   `json:"fixed_team_id,omitempty"`
	OwningTeamIDs             []string `json:"owning_team_ids,omitempty"`

	Template              json.RawMessage `json:"template,omitempty"`
	FilterConditionGroups json.RawMessage `json:"filter_condition_groups,omitempty"`
	RateLimitSharding     json.RawMessage `json:"rate_limit_sharding,omitempty"`
	AzureDevopsOptions    json.RawMessage `json:"azure_devops_options,omitempty"`
	EmailOptions          json.RawMessage `json:"email_options,omitempty"`
	HeartbeatOptions      json.RawMessage `json:"heartbeat_options,omitempty"`
	HTTPCustomOptions     json.RawMessage `json:"http_custom_options,omitempty"`
	JiraOptions           json.RawMessage `json:"jira_options,omitempty"`
}

// CreateAlertSourceOptions represents the options for creating an alert source.
type CreateAlertSourceOptions struct {
	Name       string          `json:"name"`
	SourceType string          `json:"source_type"`
	Template   json.RawMessage `json:"template"`

	AutoResolveIncidentAlerts bool     `json:"auto_resolve_incident_alerts,omitempty"`
	AutoResolveTimeoutMinutes int      `json:"auto_resolve_timeout_minutes,omitempty"`
	FixedTeamID               string   `json:"fixed_team_id,omitempty"`
	OwningTeamIDs             []string `json:"owning_team_ids,omitempty"`

	FilterConditionGroups json.RawMessage `json:"filter_condition_groups,omitempty"`
	RateLimitSharding     json.RawMessage `json:"rate_limit_sharding,omitempty"`
	AzureDevopsOptions    json.RawMessage `json:"azure_devops_options,omitempty"`
	EmailOptions          json.RawMessage `json:"email_options,omitempty"`
	HeartbeatOptions      json.RawMessage `json:"heartbeat_options,omitempty"`
	HTTPCustomOptions     json.RawMessage `json:"http_custom_options,omitempty"`
	JiraOptions           json.RawMessage `json:"jira_options,omitempty"`
}

// UpdateAlertSourceOptions represents the options for updating an alert source.
// The API replaces the whole source, so Name and Template are always required.
type UpdateAlertSourceOptions struct {
	Name     string          `json:"name"`
	Template json.RawMessage `json:"template"`

	Disabled                  *bool    `json:"disabled,omitempty"`
	AutoResolveIncidentAlerts bool     `json:"auto_resolve_incident_alerts,omitempty"`
	AutoResolveTimeoutMinutes int      `json:"auto_resolve_timeout_minutes,omitempty"`
	FixedTeamID               string   `json:"fixed_team_id,omitempty"`
	OwningTeamIDs             []string `json:"owning_team_ids,omitempty"`

	FilterConditionGroups json.RawMessage `json:"filter_condition_groups,omitempty"`
	RateLimitSharding     json.RawMessage `json:"rate_limit_sharding,omitempty"`
	AzureDevopsOptions    json.RawMessage `json:"azure_devops_options,omitempty"`
	EmailOptions          json.RawMessage `json:"email_options,omitempty"`
	HeartbeatOptions      json.RawMessage `json:"heartbeat_options,omitempty"`
	HTTPCustomOptions     json.RawMessage `json:"http_custom_options,omitempty"`
	JiraOptions           json.RawMessage `json:"jira_options,omitempty"`
}

// List returns all alert sources.
func (s *AlertSourcesService) List(ctx context.Context) ([]*AlertSource, *http.Response, error) {
	return getKey[[]*AlertSource](ctx, s.client, "GET", "v2/alert_sources", nil, nil, "alert_sources")
}

// Get returns a single alert source.
func (s *AlertSourcesService) Get(ctx context.Context, id string) (*AlertSource, *http.Response, error) {
	return getKey[*AlertSource](ctx, s.client, "GET", fmt.Sprintf("v2/alert_sources/%s", id), nil, nil, "alert_source")
}

// Create creates an alert source.
func (s *AlertSourcesService) Create(ctx context.Context, opts *CreateAlertSourceOptions) (*AlertSource, *http.Response, error) {
	return getKey[*AlertSource](ctx, s.client, "POST", "v2/alert_sources", nil, opts, "alert_source")
}

// Update replaces an alert source.
func (s *AlertSourcesService) Update(ctx context.Context, id string, opts *UpdateAlertSourceOptions) (*AlertSource, *http.Response, error) {
	return getKey[*AlertSource](ctx, s.client, "PUT", fmt.Sprintf("v2/alert_sources/%s", id), nil, opts, "alert_source")
}

// Delete deletes an alert source.
func (s *AlertSourcesService) Delete(ctx context.Context, id string) (*http.Response, error) {
	return s.client.send(ctx, "DELETE", fmt.Sprintf("v2/alert_sources/%s", id), nil, nil, nil)
}

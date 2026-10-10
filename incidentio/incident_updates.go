package incidentio

import (
	"context"
	"net/http"
	"net/url"
)

// IncidentUpdatesService handles communication with the incident updates related methods.
type IncidentUpdatesService struct {
	client *Client
}

// IncidentUpdate represents a status or severity update posted to an incident.
type IncidentUpdate struct {
	ID                   string          `json:"id"`
	IncidentID           string          `json:"incident_id"`
	Message              string          `json:"message,omitempty"`
	NewIncidentStatus    *IncidentStatus `json:"new_incident_status,omitempty"`
	NewSeverity          *Severity       `json:"new_severity,omitempty"`
	MergedIntoIncidentID string          `json:"merged_into_incident_id,omitempty"`
	Updater              *Actor          `json:"updater,omitempty"`
	CreatedAt            Timestamp       `json:"created_at"`
}

// IncidentUpdateListOptions represents the options for listing incident updates.
type IncidentUpdateListOptions struct {
	ListOptions
	IncidentID string
}

// CreateIncidentUpdateOptions represents the options for creating an incident update.
type CreateIncidentUpdateOptions struct {
	IncidentID string `json:"incident_id"`
	// IdempotencyKey makes creation safe to retry. A random key is generated when empty.
	IdempotencyKey     string `json:"idempotency_key"`
	Message            string `json:"message,omitempty"`
	ToIncidentStatusID string `json:"to_incident_status_id,omitempty"`
	ToSeverityID       string `json:"to_severity_id,omitempty"`
}

// List returns incident updates.
func (s *IncidentUpdatesService) List(ctx context.Context, opts *IncidentUpdateListOptions) ([]*IncidentUpdate, *http.Response, error) {
	q := url.Values{}
	if opts != nil {
		opts.apply(q)
		if opts.IncidentID != "" {
			q.Set("incident_id", opts.IncidentID)
		}
	}

	return getKey[[]*IncidentUpdate](ctx, s.client, "GET", "v2/incident_updates", q, nil, "incident_updates")
}

// Create posts an update to an incident.
func (s *IncidentUpdatesService) Create(ctx context.Context, opts *CreateIncidentUpdateOptions) (*IncidentUpdate, *http.Response, error) {
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

	return getKey[*IncidentUpdate](ctx, s.client, "POST", "v2/incident_updates", nil, body, "incident_update")
}

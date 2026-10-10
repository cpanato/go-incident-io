package incidentio

import (
	"context"
	"net/http"
	"net/url"
)

// IncidentTimelineItemsService handles communication with the incident timeline related methods.
type IncidentTimelineItemsService struct {
	client *Client
}

// IncidentTimelineItem represents an entry on an incident's timeline.
type IncidentTimelineItem struct {
	ID            string    `json:"id"`
	IncidentID    string    `json:"incident_id"`
	ActivityLogID string    `json:"activity_log_id,omitempty"`
	Title         string    `json:"title"`
	Description   string    `json:"description,omitempty"`
	Timestamp     Timestamp `json:"timestamp"`
	Creator       *Actor    `json:"creator,omitempty"`
	CreatedAt     Timestamp `json:"created_at"`
	UpdatedAt     Timestamp `json:"updated_at"`
}

// IncidentTimelineItemListOptions represents the options for listing timeline items.
type IncidentTimelineItemListOptions struct {
	ListOptions
	IncidentID string
}

// CreateIncidentTimelineItemOptions represents the options for creating a timeline item.
type CreateIncidentTimelineItemOptions struct {
	IncidentID string    `json:"incident_id"`
	Title      string    `json:"title"`
	Timestamp  Timestamp `json:"timestamp"`
	// IdempotencyKey makes creation safe to retry. A random key is generated when empty.
	IdempotencyKey string `json:"idempotency_key"`
	Description    string `json:"description,omitempty"`
}

// List returns timeline items.
func (s *IncidentTimelineItemsService) List(ctx context.Context, opts *IncidentTimelineItemListOptions) ([]*IncidentTimelineItem, *http.Response, error) {
	q := url.Values{}
	if opts != nil {
		opts.apply(q)
		if opts.IncidentID != "" {
			q.Set("incident_id", opts.IncidentID)
		}
	}

	return getKey[[]*IncidentTimelineItem](ctx, s.client, "GET", "v2/incident_timeline_items", q, nil, "incident_timeline_items")
}

// Create adds an item to an incident's timeline.
func (s *IncidentTimelineItemsService) Create(ctx context.Context, opts *CreateIncidentTimelineItemOptions) (*IncidentTimelineItem, *http.Response, error) {
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

	return getKey[*IncidentTimelineItem](ctx, s.client, "POST", "v2/incident_timeline_items", nil, body, "incident_timeline_item")
}

package incidentio

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

// AlertsService handles communication with the alerts related methods.
type AlertsService struct {
	client *Client
}

// AlertTag is a tag attached to an alert.
type AlertTag struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// AlertAttribute describes an alert attribute definition.
type AlertAttribute struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Type     string `json:"type"`
	Array    bool   `json:"array"`
	Required bool   `json:"required"`
	Emoji    string `json:"emoji,omitempty"`
}

// AlertAttributeValue is a single value of an alert attribute.
type AlertAttributeValue struct {
	Label        string          `json:"label,omitempty"`
	Literal      string          `json:"literal,omitempty"`
	CatalogEntry json.RawMessage `json:"catalog_entry,omitempty"`
}

// AlertAttributeEntry is an attribute along with its value(s) for an alert.
type AlertAttributeEntry struct {
	Attribute  AlertAttribute         `json:"attribute"`
	Value      *AlertAttributeValue   `json:"value,omitempty"`
	ArrayValue []*AlertAttributeValue `json:"array_value,omitempty"`
}

// Alert represents an alert received from an alert source.
type Alert struct {
	ID               string                 `json:"id"`
	AlertSourceID    string                 `json:"alert_source_id"`
	DeduplicationKey string                 `json:"deduplication_key"`
	Status           string                 `json:"status"`
	Title            string                 `json:"title"`
	Description      string                 `json:"description,omitempty"`
	SourceURL        string                 `json:"source_url,omitempty"`
	AlertGroupIDs    []string               `json:"alert_group_ids,omitempty"`
	Attributes       []*AlertAttributeEntry `json:"attributes"`
	Tags             []*AlertTag            `json:"tags,omitempty"`
	CreatedAt        Timestamp              `json:"created_at"`
	UpdatedAt        Timestamp              `json:"updated_at"`
	ResolvedAt       *Timestamp             `json:"resolved_at,omitempty"`
}

// AlertListOptions represents the options for listing alerts.
type AlertListOptions struct {
	ListOptions
	DeduplicationKey Filter
	Status           Filter
	AlertSource      Filter
	AlertGroupID     Filter
	CreatedAt        Filter
	UpdatedAt        Filter
	Tags             Filter
	HasNotes         Filter
	// IncludeMaintenanceWindow is a filter such as Filter{"is": {"true"}}.
	IncludeMaintenanceWindow Filter
}

// List returns alerts.
func (s *AlertsService) List(ctx context.Context, opts *AlertListOptions) ([]*Alert, *http.Response, error) {
	q := url.Values{}
	if opts != nil {
		opts.apply(q)
		for name, f := range map[string]Filter{
			"deduplication_key":          opts.DeduplicationKey,
			"status":                     opts.Status,
			"alert_source":               opts.AlertSource,
			"alert_group_id":             opts.AlertGroupID,
			"created_at":                 opts.CreatedAt,
			"updated_at":                 opts.UpdatedAt,
			"tags":                       opts.Tags,
			"has_notes":                  opts.HasNotes,
			"include_maintenance_window": opts.IncludeMaintenanceWindow,
		} {
			addFilter(q, name, f)
		}
	}

	return getKey[[]*Alert](ctx, s.client, "GET", "v2/alerts", q, nil, "alerts")
}

// Get returns a single alert.
func (s *AlertsService) Get(ctx context.Context, id string) (*Alert, *http.Response, error) {
	return getKey[*Alert](ctx, s.client, "GET", fmt.Sprintf("v2/alerts/%s", id), nil, nil, "alert")
}

// Resolve marks an alert as resolved.
func (s *AlertsService) Resolve(ctx context.Context, id string) (*Alert, *http.Response, error) {
	return getKey[*Alert](ctx, s.client, "POST", fmt.Sprintf("v2/alerts/%s/actions/resolve", id), nil, struct{}{}, "alert")
}

type alertTagsPayload struct {
	Tags []string `json:"tags"`
}

// AddTags adds tags to an alert.
func (s *AlertsService) AddTags(ctx context.Context, id string, tags []string) (*Alert, *http.Response, error) {
	return s.tagsAction(ctx, id, "add_tags", tags)
}

// RemoveTags removes tags from an alert.
func (s *AlertsService) RemoveTags(ctx context.Context, id string, tags []string) (*Alert, *http.Response, error) {
	return s.tagsAction(ctx, id, "remove_tags", tags)
}

// SetTags replaces all tags on an alert.
func (s *AlertsService) SetTags(ctx context.Context, id string, tags []string) (*Alert, *http.Response, error) {
	return s.tagsAction(ctx, id, "set_tags", tags)
}

func (s *AlertsService) tagsAction(ctx context.Context, id, action string, tags []string) (*Alert, *http.Response, error) {
	if tags == nil {
		tags = []string{}
	}

	return getKey[*Alert](ctx, s.client, "POST", fmt.Sprintf("v2/alerts/%s/actions/%s", id, action), nil, alertTagsPayload{Tags: tags}, "alert")
}

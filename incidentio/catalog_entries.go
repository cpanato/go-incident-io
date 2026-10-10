package incidentio

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
)

// CatalogEntriesService handles communication with the catalog entries related methods.
type CatalogEntriesService struct {
	client *Client
}

// CatalogBindingValue is a single value of a catalog attribute.
// Label is only populated in responses.
type CatalogBindingValue struct {
	Literal string `json:"literal,omitempty"`
	Label   string `json:"label,omitempty"`
}

// CatalogAttributeBinding holds the value (or array of values) of a catalog attribute.
type CatalogAttributeBinding struct {
	Value      *CatalogBindingValue   `json:"value,omitempty"`
	ArrayValue []*CatalogBindingValue `json:"array_value,omitempty"`
}

// CatalogEntry represents an entry of a catalog type.
type CatalogEntry struct {
	ID              string                              `json:"id"`
	CatalogTypeID   string                              `json:"catalog_type_id"`
	Name            string                              `json:"name"`
	ExternalID      string                              `json:"external_id,omitempty"`
	Aliases         []string                            `json:"aliases"`
	Rank            int                                 `json:"rank"`
	AttributeValues map[string]*CatalogAttributeBinding `json:"attribute_values"`
	ArchivedAt      *Timestamp                          `json:"archived_at,omitempty"`
	CreatedAt       Timestamp                           `json:"created_at"`
	UpdatedAt       Timestamp                           `json:"updated_at"`
}

// CatalogEntryListOptions represents the options for listing catalog entries.
type CatalogEntryListOptions struct {
	ListOptions
	// CatalogTypeID is required by the API.
	CatalogTypeID string
	// Identifier matches on ID, external ID and aliases.
	Identifier string
}

// CreateCatalogEntryOptions represents the options for creating a catalog entry.
// AttributeValues is keyed by attribute ID.
type CreateCatalogEntryOptions struct {
	CatalogTypeID   string                              `json:"catalog_type_id"`
	Name            string                              `json:"name"`
	AttributeValues map[string]*CatalogAttributeBinding `json:"attribute_values"`
	Aliases         []string                            `json:"aliases,omitempty"`
	ExternalID      string                              `json:"external_id,omitempty"`
	Rank            *int                                `json:"rank,omitempty"`
}

// UpdateCatalogEntryOptions represents the options for updating a catalog entry.
type UpdateCatalogEntryOptions struct {
	Name            string                              `json:"name"`
	AttributeValues map[string]*CatalogAttributeBinding `json:"attribute_values"`
	Aliases         []string                            `json:"aliases,omitempty"`
	ExternalID      string                              `json:"external_id,omitempty"`
	Rank            *int                                `json:"rank,omitempty"`
	// UpdateAttributes limits the update to the listed attribute IDs.
	UpdateAttributes []string `json:"update_attributes,omitempty"`
}

// PartialCatalogEntry is one entry in a bulk update.
type PartialCatalogEntry struct {
	EntryID         string                              `json:"entry_id"`
	AttributeValues map[string]*CatalogAttributeBinding `json:"attribute_values"`
	Name            string                              `json:"name,omitempty"`
	Aliases         []string                            `json:"aliases,omitempty"`
	ExternalID      string                              `json:"external_id,omitempty"`
	Rank            *int                                `json:"rank,omitempty"`
}

// BulkUpdateCatalogEntriesOptions represents the options for bulk updating entries of a type.
type BulkUpdateCatalogEntriesOptions struct {
	CatalogTypeID    string                 `json:"catalog_type_id"`
	Entries          []*PartialCatalogEntry `json:"entries"`
	UpdateAttributes []string               `json:"update_attributes,omitempty"`
}

// List returns the entries of a catalog type.
func (s *CatalogEntriesService) List(ctx context.Context, opts *CatalogEntryListOptions) ([]*CatalogEntry, *http.Response, error) {
	q := url.Values{}
	if opts != nil {
		opts.apply(q)
		if opts.CatalogTypeID != "" {
			q.Set("catalog_type_id", opts.CatalogTypeID)
		}
		if opts.Identifier != "" {
			q.Set("identifier", opts.Identifier)
		}
	}

	return getKey[[]*CatalogEntry](ctx, s.client, "GET", "v3/catalog_entries", q, nil, "catalog_entries")
}

// Get returns a single catalog entry.
func (s *CatalogEntriesService) Get(ctx context.Context, id string) (*CatalogEntry, *http.Response, error) {
	return getKey[*CatalogEntry](ctx, s.client, "GET", fmt.Sprintf("v3/catalog_entries/%s", id), nil, nil, "catalog_entry")
}

// Create creates a catalog entry.
func (s *CatalogEntriesService) Create(ctx context.Context, opts *CreateCatalogEntryOptions) (*CatalogEntry, *http.Response, error) {
	return getKey[*CatalogEntry](ctx, s.client, "POST", "v3/catalog_entries", nil, opts, "catalog_entry")
}

// Update updates a catalog entry.
func (s *CatalogEntriesService) Update(ctx context.Context, id string, opts *UpdateCatalogEntryOptions) (*CatalogEntry, *http.Response, error) {
	return getKey[*CatalogEntry](ctx, s.client, "PUT", fmt.Sprintf("v3/catalog_entries/%s", id), nil, opts, "catalog_entry")
}

// Delete deletes a catalog entry.
func (s *CatalogEntriesService) Delete(ctx context.Context, id string) (*http.Response, error) {
	return s.client.send(ctx, "DELETE", fmt.Sprintf("v3/catalog_entries/%s", id), nil, nil, nil)
}

// BulkUpdate updates many entries of a catalog type in one request.
func (s *CatalogEntriesService) BulkUpdate(ctx context.Context, opts *BulkUpdateCatalogEntriesOptions) (*http.Response, error) {
	return s.client.send(ctx, "POST", "v3/catalog_entries/actions/bulk_update", nil, opts, nil)
}

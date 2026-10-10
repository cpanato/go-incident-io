package incidentio

import (
	"context"
	"fmt"
	"net/http"
)

// CatalogTypesService handles communication with the catalog types related methods.
type CatalogTypesService struct {
	client *Client
}

// CatalogTypeAttributePathItem is a step in the path of a path-mode attribute.
type CatalogTypeAttributePathItem struct {
	AttributeID   string `json:"attribute_id,omitempty"`
	AttributeName string `json:"attribute_name,omitempty"`
}

// CatalogTypeAttribute is an attribute in a catalog type's schema.
type CatalogTypeAttribute struct {
	ID                string                         `json:"id,omitempty"`
	Name              string                         `json:"name"`
	Type              string                         `json:"type"`
	Array             bool                           `json:"array"`
	Mode              string                         `json:"mode,omitempty"`
	BacklinkAttribute string                         `json:"backlink_attribute,omitempty"`
	Path              []CatalogTypeAttributePathItem `json:"path,omitempty"`
}

// CatalogTypeSchema is the versioned schema of a catalog type.
type CatalogTypeSchema struct {
	Version    int                     `json:"version"`
	Attributes []*CatalogTypeAttribute `json:"attributes"`
}

// CatalogType represents a catalog type.
type CatalogType struct {
	ID                   string             `json:"id"`
	Name                 string             `json:"name"`
	Description          string             `json:"description"`
	TypeName             string             `json:"type_name,omitempty"`
	EngineResourceType   string             `json:"engine_resource_type,omitempty"`
	RegistryType         string             `json:"registry_type,omitempty"`
	Ranked               bool               `json:"ranked"`
	Icon                 string             `json:"icon,omitempty"`
	Color                string             `json:"color,omitempty"`
	Categories           []string           `json:"categories,omitempty"`
	IsEditable           bool               `json:"is_editable"`
	IsTeamType           bool               `json:"is_team_type,omitempty"`
	UseNameAsIdentifier  bool               `json:"use_name_as_identifier"`
	EstimatedCount       int                `json:"estimated_count,omitempty"`
	OwningTeamIDs        []string           `json:"owning_team_ids,omitempty"`
	RequiredIntegrations []string           `json:"required_integrations,omitempty"`
	SourceRepoURL        string             `json:"source_repo_url,omitempty"`
	Annotations          map[string]string  `json:"annotations,omitempty"`
	Schema               *CatalogTypeSchema `json:"schema,omitempty"`
	LastSyncedAt         *Timestamp         `json:"last_synced_at,omitempty"`
	CreatedAt            Timestamp          `json:"created_at"`
	UpdatedAt            Timestamp          `json:"updated_at"`
}

// CreateCatalogTypeOptions represents the options for creating a catalog type.
type CreateCatalogTypeOptions struct {
	Name                string            `json:"name"`
	Description         string            `json:"description"`
	TypeName            string            `json:"type_name,omitempty"`
	Ranked              *bool             `json:"ranked,omitempty"`
	Icon                string            `json:"icon,omitempty"`
	Color               string            `json:"color,omitempty"`
	Categories          []string          `json:"categories,omitempty"`
	UseNameAsIdentifier *bool             `json:"use_name_as_identifier,omitempty"`
	OwningTeamIDs       []string          `json:"owning_team_ids,omitempty"`
	SourceRepoURL       string            `json:"source_repo_url,omitempty"`
	Annotations         map[string]string `json:"annotations,omitempty"`
}

// UpdateCatalogTypeOptions represents the options for updating a catalog type.
type UpdateCatalogTypeOptions struct {
	Name                string            `json:"name"`
	Description         string            `json:"description"`
	Ranked              *bool             `json:"ranked,omitempty"`
	Icon                string            `json:"icon,omitempty"`
	Color               string            `json:"color,omitempty"`
	Categories          []string          `json:"categories,omitempty"`
	UseNameAsIdentifier *bool             `json:"use_name_as_identifier,omitempty"`
	OwningTeamIDs       []string          `json:"owning_team_ids,omitempty"`
	SourceRepoURL       string            `json:"source_repo_url,omitempty"`
	Annotations         map[string]string `json:"annotations,omitempty"`
}

// UpdateCatalogTypeSchemaOptions represents the options for replacing a catalog type's schema.
// Version must match the current schema version.
type UpdateCatalogTypeSchemaOptions struct {
	Version    int                     `json:"version"`
	Attributes []*CatalogTypeAttribute `json:"attributes"`
}

// List returns all catalog types.
func (s *CatalogTypesService) List(ctx context.Context) ([]*CatalogType, *http.Response, error) {
	return getKey[[]*CatalogType](ctx, s.client, "GET", "v3/catalog_types", nil, nil, "catalog_types")
}

// Get returns a single catalog type.
func (s *CatalogTypesService) Get(ctx context.Context, id string) (*CatalogType, *http.Response, error) {
	return getKey[*CatalogType](ctx, s.client, "GET", fmt.Sprintf("v3/catalog_types/%s", id), nil, nil, "catalog_type")
}

// Create creates a catalog type.
func (s *CatalogTypesService) Create(ctx context.Context, opts *CreateCatalogTypeOptions) (*CatalogType, *http.Response, error) {
	return getKey[*CatalogType](ctx, s.client, "POST", "v3/catalog_types", nil, opts, "catalog_type")
}

// Update updates a catalog type.
func (s *CatalogTypesService) Update(ctx context.Context, id string, opts *UpdateCatalogTypeOptions) (*CatalogType, *http.Response, error) {
	return getKey[*CatalogType](ctx, s.client, "PUT", fmt.Sprintf("v3/catalog_types/%s", id), nil, opts, "catalog_type")
}

// UpdateSchema replaces the schema of a catalog type.
func (s *CatalogTypesService) UpdateSchema(ctx context.Context, id string, opts *UpdateCatalogTypeSchemaOptions) (*CatalogType, *http.Response, error) {
	return getKey[*CatalogType](ctx, s.client, "POST", fmt.Sprintf("v3/catalog_types/%s/actions/update_schema", id), nil, opts, "catalog_type")
}

// Delete deletes a catalog type.
func (s *CatalogTypesService) Delete(ctx context.Context, id string) (*http.Response, error) {
	return s.client.send(ctx, "DELETE", fmt.Sprintf("v3/catalog_types/%s", id), nil, nil, nil)
}

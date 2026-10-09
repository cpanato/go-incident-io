package incidentio

import (
	"context"
	"fmt"
	"net/http"
)

// CustomFieldsService handles communication with the custom fields related methods.
type CustomFieldsService struct {
	client *Client
}

// CustomField represents a custom field in Incident.io.
type CustomField struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	FieldType   string `json:"field_type"`
	Options     []struct {
		ID    string `json:"id"`
		Value string `json:"value"`
		Label string `json:"label"`
	} `json:"options,omitempty"`
	Required                   bool      `json:"required"`
	CatalogTypeID              string    `json:"catalog_type_id,omitempty"`
	GroupByCatalogAttributeID  string    `json:"group_by_catalog_attribute_id,omitempty"`
	HelptextCatalogAttributeID string    `json:"helptext_catalog_attribute_id,omitempty"`
	CreatedAt                  Timestamp `json:"created_at"`
	UpdatedAt                  Timestamp `json:"updated_at"`
}

// CreateCustomFieldOptions represents the options for creating a custom field.
type CreateCustomFieldOptions struct {
	Name                       string `json:"name"`
	Description                string `json:"description"`
	FieldType                  string `json:"field_type"`
	CatalogTypeID              string `json:"catalog_type_id,omitempty"`
	GroupByCatalogAttributeID  string `json:"group_by_catalog_attribute_id,omitempty"`
	HelptextCatalogAttributeID string `json:"helptext_catalog_attribute_id,omitempty"`
}

// UpdateCustomFieldOptions represents the options for updating a custom field.
type UpdateCustomFieldOptions struct {
	Name                       string `json:"name"`
	Description                string `json:"description"`
	GroupByCatalogAttributeID  string `json:"group_by_catalog_attribute_id,omitempty"`
	HelptextCatalogAttributeID string `json:"helptext_catalog_attribute_id,omitempty"`
}

// List returns a list of custom fields.
func (s *CustomFieldsService) List(ctx context.Context) ([]*CustomField, *http.Response, error) {
	u := "v2/custom_fields"

	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	var result struct {
		CustomFields []*CustomField `json:"custom_fields"`
	}

	resp, err := s.client.Do(ctx, req, &result)
	if err != nil {
		return nil, resp, err
	}

	return result.CustomFields, resp, nil
}

// Get returns a single custom field.
func (s *CustomFieldsService) Get(ctx context.Context, id string) (*CustomField, *http.Response, error) {
	req, err := s.client.NewRequest("GET", fmt.Sprintf("v2/custom_fields/%s", id), nil)
	if err != nil {
		return nil, nil, err
	}

	var result struct {
		CustomField *CustomField `json:"custom_field"`
	}
	resp, err := s.client.Do(ctx, req, &result)
	if err != nil {
		return nil, resp, err
	}

	return result.CustomField, resp, nil
}

// Create creates a new custom field.
func (s *CustomFieldsService) Create(ctx context.Context, opts *CreateCustomFieldOptions) (*CustomField, *http.Response, error) {
	req, err := s.client.NewRequest("POST", "v2/custom_fields", opts)
	if err != nil {
		return nil, nil, err
	}

	var result struct {
		CustomField *CustomField `json:"custom_field"`
	}
	resp, err := s.client.Do(ctx, req, &result)
	if err != nil {
		return nil, resp, err
	}

	return result.CustomField, resp, nil
}

// Update updates a custom field.
func (s *CustomFieldsService) Update(ctx context.Context, id string, opts *UpdateCustomFieldOptions) (*CustomField, *http.Response, error) {
	req, err := s.client.NewRequest("PUT", fmt.Sprintf("v2/custom_fields/%s", id), opts)
	if err != nil {
		return nil, nil, err
	}

	var result struct {
		CustomField *CustomField `json:"custom_field"`
	}
	resp, err := s.client.Do(ctx, req, &result)
	if err != nil {
		return nil, resp, err
	}

	return result.CustomField, resp, nil
}

// Delete deletes a custom field.
func (s *CustomFieldsService) Delete(ctx context.Context, id string) (*http.Response, error) {
	req, err := s.client.NewRequest("DELETE", fmt.Sprintf("v2/custom_fields/%s", id), nil)
	if err != nil {
		return nil, err
	}

	return s.client.Do(ctx, req, nil)
}

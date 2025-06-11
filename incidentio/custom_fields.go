package incidentio

import (
	"context"
	"net/http"
)

// CustomFieldsService handles communication with the custom fields related methods.
type CustomFieldsService struct {
	client *Client
}

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
	Required  bool      `json:"required"`
	CreatedAt Timestamp `json:"created_at"`
	UpdatedAt Timestamp `json:"updated_at"`
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

package incidentio

// WebhooksService handles communication with the webhooks related methods.
type WebhooksService struct {
	client *Client
}

// Webhook represents a webhook in Incident.io.
type Webhook struct {
	ID               string    `json:"id"`
	Name             string    `json:"name"`
	Endpoint         string    `json:"endpoint"`
	PrivateIncidents bool      `json:"private_incidents"`
	CreatedAt        Timestamp `json:"created_at"`
	UpdatedAt        Timestamp `json:"updated_at"`
}

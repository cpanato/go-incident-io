package incidentio

// ActionsService handles communication with the actions related methods.
type ActionsService struct {
	client *Client
}

type Action struct {
	ID          string     `json:"id"`
	IncidentID  string     `json:"incident_id"`
	Description string     `json:"description"`
	Status      string     `json:"status"`
	Assignee    *User      `json:"assignee,omitempty"`
	CompletedAt *Timestamp `json:"completed_at,omitempty"`
	CreatedAt   Timestamp  `json:"created_at"`
	UpdatedAt   Timestamp  `json:"updated_at"`
}

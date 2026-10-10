package incidentio

import (
	"context"
	"net/http"
	"net/url"
)

// IncidentParticipantsService handles communication with the incident participants related methods.
type IncidentParticipantsService struct {
	client *Client
}

// IncidentParticipant is a user taking part in an incident.
type IncidentParticipant struct {
	User *User `json:"user"`
	// ParticipantType is one of "observer", "collaborator" or "responder".
	ParticipantType string `json:"participant_type"`
}

// IncidentParticipants groups the participants of an incident.
type IncidentParticipants struct {
	Active  []*IncidentParticipant `json:"active"`
	Passive []*IncidentParticipant `json:"passive"`
}

// List returns the participants of an incident.
func (s *IncidentParticipantsService) List(ctx context.Context, incidentID string) (*IncidentParticipants, *http.Response, error) {
	q := url.Values{"incident_id": {incidentID}}

	return getKey[*IncidentParticipants](ctx, s.client, "GET", "v2/incident_participants", q, nil, "incident_participants")
}

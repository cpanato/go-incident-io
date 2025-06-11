package incidentio

import (
	"context"
	"net/http"
)

// UsersService handles communication with the users related methods.
type UsersService struct {
	client *Client
}

// User represents a user in Incident.io.
type User struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
	Role  string `json:"role"`
	Slack *struct {
		UserID string `json:"user_id"`
		TeamID string `json:"team_id"`
	} `json:"slack,omitempty"`
}

// List returns a list of users.
func (s *UsersService) List(ctx context.Context, _ *ListOptions) ([]*User, *http.Response, error) {
	u := "v2/users"

	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	var result struct {
		Users []*User `json:"users"`
	}
	resp, err := s.client.Do(ctx, req, &result)
	if err != nil {
		return nil, resp, err
	}

	return result.Users, resp, nil
}

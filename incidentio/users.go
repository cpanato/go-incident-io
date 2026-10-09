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
	// SlackUserID is the flat Slack user ID returned by the API.
	SlackUserID string `json:"slack_user_id,omitempty"`
	Slack       *struct {
		UserID string `json:"user_id"`
		TeamID string `json:"team_id"`
	} `json:"slack,omitempty"`
}

// UserListOptions represents the options for listing users.
type UserListOptions struct {
	ListOptions
	Email           string
	SlackUserID     string
	IncludeInactive bool
}

// List returns a list of users.
func (s *UsersService) List(ctx context.Context, opts *UserListOptions) ([]*User, *http.Response, error) {
	u := "v2/users"

	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	if opts != nil {
		q := req.URL.Query()
		opts.apply(q)
		if opts.Email != "" {
			q.Set("email", opts.Email)
		}
		if opts.SlackUserID != "" {
			q.Set("slack_user_id", opts.SlackUserID)
		}
		if opts.IncludeInactive {
			q.Set("include_inactive", "true")
		}
		req.URL.RawQuery = q.Encode()
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

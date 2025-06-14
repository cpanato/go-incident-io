package incidentio

import (
	"context"
	"fmt"
	"net/http"
	"reflect"
	"testing"
)

func TestUsersService_List(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/v2/users", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testHeader(t, r, "Authorization", "Bearer test-key")

		response := `{
			"users": [
				{
					"id": "01FCNDV6P870EA6S7TK1DSYDG0",
					"name": "John Doe",
					"email": "john@example.com",
					"role": "owner",
					"slack": {
						"user_id": "U01FCNDV6P8",
						"team_id": "T01FCNDV6P8"
					}
				},
				{
					"id": "01FCNDV6P870EA6S7TK1DSYDG1",
					"name": "Jane Smith",
					"email": "jane@example.com",
					"role": "admin"
				}
			]
		}`

		_, _ = fmt.Fprint(w, response)
	})

	ctx := context.Background()
	users, _, err := client.Users.List(ctx, nil)
	if err != nil {
		t.Errorf("Users.List returned error: %v", err)
	}

	expected := []*User{
		{
			ID:    "01FCNDV6P870EA6S7TK1DSYDG0",
			Name:  "John Doe",
			Email: "john@example.com",
			Role:  "owner",
			Slack: &struct {
				UserID string `json:"user_id"`
				TeamID string `json:"team_id"`
			}{
				UserID: "U01FCNDV6P8",
				TeamID: "T01FCNDV6P8",
			},
		},
		{
			ID:    "01FCNDV6P870EA6S7TK1DSYDG1",
			Name:  "Jane Smith",
			Email: "jane@example.com",
			Role:  "admin",
		},
	}

	if !reflect.DeepEqual(users, expected) {
		t.Errorf("Users.List returned %+v, want %+v", users, expected)
	}
}

func TestUsersService_List_EmptyResponse(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/v2/users", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testHeader(t, r, "Authorization", "Bearer test-key")

		response := `{
			"users": []
		}`

		_, _ = fmt.Fprint(w, response)
	})

	ctx := context.Background()
	users, _, err := client.Users.List(ctx, nil)
	if err != nil {
		t.Errorf("Users.List returned error: %v", err)
	}

	if len(users) != 0 {
		t.Errorf("Users.List returned %d users, want 0", len(users))
	}
}

func TestUsersService_List_HTTPError(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/v2/users", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		w.WriteHeader(http.StatusInternalServerError)
		response := `{
			"type": "internal_error",
			"status": 500,
			"detail": "Internal server error"
		}`
		_, _ = fmt.Fprint(w, response)
	})

	ctx := context.Background()
	users, resp, err := client.Users.List(ctx, nil)

	if err == nil {
		t.Error("Expected error, got nil")
	}

	if users != nil {
		t.Errorf("Users.List returned users %+v, want nil", users)
	}

	if resp == nil {
		t.Error("Expected response, got nil")
	} else if resp.StatusCode != http.StatusInternalServerError {
		t.Errorf("Expected status %d, got %d", http.StatusInternalServerError, resp.StatusCode)
	}
}

func TestUsersService_List_InvalidJSON(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/v2/users", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		response := `{"invalid": json}`
		_, _ = fmt.Fprint(w, response)
	})

	ctx := context.Background()
	users, _, err := client.Users.List(ctx, nil)

	if err == nil {
		t.Error("Expected error, got nil")
	}

	if users != nil {
		t.Errorf("Users.List returned users %+v, want nil", users)
	}
}

func TestUsersService_List_WithoutSlack(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/v2/users", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testHeader(t, r, "Authorization", "Bearer test-key")

		response := `{
			"users": [
				{
					"id": "01FCNDV6P870EA6S7TK1DSYDG0",
					"name": "John Doe",
					"email": "john@example.com",
					"role": "owner"
				}
			]
		}`

		_, _ = fmt.Fprint(w, response)
	})

	ctx := context.Background()
	users, _, err := client.Users.List(ctx, nil)
	if err != nil {
		t.Errorf("Users.List returned error: %v", err)
	}

	if len(users) != 1 {
		t.Errorf("Users.List returned %d users, want 1", len(users))
		return
	}

	user := users[0]
	if user.Slack != nil {
		t.Errorf("Users.List returned user with Slack data %+v, want nil", user.Slack)
	}

	expected := &User{
		ID:    "01FCNDV6P870EA6S7TK1DSYDG0",
		Name:  "John Doe",
		Email: "john@example.com",
		Role:  "owner",
	}

	if !reflect.DeepEqual(user, expected) {
		t.Errorf("Users.List returned %+v, want %+v", user, expected)
	}
}

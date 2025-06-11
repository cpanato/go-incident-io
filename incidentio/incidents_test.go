package incidentio

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"testing"
	"time"
)

// Test helpers

func setup() (client *Client, mux *http.ServeMux, serverURL string, teardown func()) { //nolint: unparam
	mux = http.NewServeMux()
	server := httptest.NewServer(mux)

	client = NewClient("test-key")
	url, _ := url.Parse(server.URL + "/v2")
	client.BaseURL = url

	return client, mux, server.URL, server.Close
}

func testMethod(t *testing.T, r *http.Request, want string) {
	t.Helper()
	if got := r.Method; got != want {
		t.Errorf("Request method: %v, want %v", got, want)
	}
}

func testHeader(t *testing.T, r *http.Request, header string, want string) {
	t.Helper()
	if got := r.Header.Get(header); got != want {
		t.Errorf("Header.Get(%q) returned %q, want %q", header, got, want)
	}
}

func parseTime(s string) time.Time {
	t, _ := time.Parse(time.RFC3339, s)
	return t
}

func TestIncidentsService_List(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/v2/incidents", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testHeader(t, r, "Authorization", "Bearer test-key")

		response := `{
			"incidents": [
				{
					"id": "01FDAG4SAP5TYPT98WGR2N7W91",
					"name": "Database Connection Issues",
					"summary": "Users experiencing connection timeouts",
					"type": "incident",
					"status": "triage",
					"severity": {
						"id": "01FH5TZRWMNAFB0DZ23FD1V96N",
						"name": "Minor",
						"description": "Minor impact",
						"rank": 1,
						"created_at": "2021-08-17T13:28:57.801578Z",
						"updated_at": "2021-08-17T13:28:57.801578Z"
					},
					"created_at": "2021-08-17T13:28:57.801578Z",
					"updated_at": "2021-08-17T13:28:57.801578Z",
					"mode": "real",
					"visibility": "public"
				}
			]
		}`

		_, _ = fmt.Fprint(w, response)
	})

	ctx := context.Background()
	incidents, _, err := client.Incidents.List(ctx, nil)
	if err != nil {
		t.Errorf("Incidents.List returned error: %v", err)
	}

	expected := []*Incident{
		{
			ID:      "01FDAG4SAP5TYPT98WGR2N7W91",
			Name:    "Database Connection Issues",
			Summary: "Users experiencing connection timeouts",
			Type:    "incident",
			Status:  "triage",
			Severity: &Severity{
				ID:          "01FH5TZRWMNAFB0DZ23FD1V96N",
				Name:        "Minor",
				Description: "Minor impact",
				Rank:        1,
				CreatedAt:   Timestamp{parseTime("2021-08-17T13:28:57.801578Z")},
				UpdatedAt:   Timestamp{parseTime("2021-08-17T13:28:57.801578Z")},
			},
			CreatedAt:  Timestamp{parseTime("2021-08-17T13:28:57.801578Z")},
			UpdatedAt:  Timestamp{parseTime("2021-08-17T13:28:57.801578Z")},
			Mode:       "real",
			Visibility: "public",
		},
	}

	if !reflect.DeepEqual(incidents, expected) {
		t.Errorf("Incidents.List returned %+v, want %+v", incidents, expected)
	}
}

func TestIncidentsService_Get(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	incidentID := "01FDAG4SAP5TYPT98WGR2N7W91"

	mux.HandleFunc(fmt.Sprintf("/v2/incidents/%s", incidentID), func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testHeader(t, r, "Authorization", "Bearer test-key")

		response := `{
			"incident": {
				"id": "01FDAG4SAP5TYPT98WGR2N7W91",
				"name": "Database Connection Issues",
				"summary": "Users experiencing connection timeouts",
				"type": "incident",
				"status": "triage",
				"created_at": "2021-08-17T13:28:57.801578Z",
				"updated_at": "2021-08-17T13:28:57.801578Z",
				"mode": "real",
				"visibility": "public",
				"slack_channel_id": "C02AW36C1M5",
				"slack_channel_name": "inc-database-issues",
				"creator": {
					"id": "01FCNDV6P870EA6S7TK1DSYDG0",
					"name": "John Doe",
					"email": "john@example.com",
					"role": "owner"
				}
			}
		}`

		_, _ = fmt.Fprint(w, response)
	})

	ctx := context.Background()
	incident, _, err := client.Incidents.Get(ctx, incidentID)
	if err != nil {
		t.Errorf("Incidents.Get returned error: %v", err)
	}

	expected := &Incident{
		ID:               "01FDAG4SAP5TYPT98WGR2N7W91",
		Name:             "Database Connection Issues",
		Summary:          "Users experiencing connection timeouts",
		Type:             "incident",
		Status:           "triage",
		CreatedAt:        Timestamp{parseTime("2021-08-17T13:28:57.801578Z")},
		UpdatedAt:        Timestamp{parseTime("2021-08-17T13:28:57.801578Z")},
		Mode:             "real",
		Visibility:       "public",
		SlackChannelID:   "C02AW36C1M5",
		SlackChannelName: "inc-database-issues",
		Creator: &User{
			ID:    "01FCNDV6P870EA6S7TK1DSYDG0",
			Name:  "John Doe",
			Email: "john@example.com",
			Role:  "owner",
		},
	}

	if !reflect.DeepEqual(incident, expected) {
		t.Errorf("Incidents.Get returned %+v, want %+v", incident, expected)
	}
}

func TestIncidentsService_Create(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	input := &CreateIncidentOptions{
		Name:           "New Database Issue",
		Summary:        "Connection pool exhausted",
		IncidentTypeID: "01FH5TZRWMNAFB0DZ23FD1V96N",
		SeverityID:     "01FH5TZRWMNAFB0DZ23FD1V96N",
		Mode:           "real",
		Visibility:     "public",
		IncidentRoleAssignments: []CreateRoleAssignment{
			{
				IncidentRoleID: "01FH5TZRWMNAFB0DZ23FD1V96N",
				UserID:         "01FCNDV6P870EA6S7TK1DSYDG0",
			},
		},
		CustomFieldValues: map[string]interface{}{
			"01FH5TZRWMNAFB0DZ23FD1V96N": "high-priority",
		},
	}

	mux.HandleFunc("/v2/incidents", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
		testHeader(t, r, "Content-Type", "application/json")
		testHeader(t, r, "Authorization", "Bearer test-key")

		var received CreateIncidentOptions
		err := json.NewDecoder(r.Body).Decode(&received)
		if err != nil {
			t.Errorf("error decoding request body: %v", err)
		}

		if !reflect.DeepEqual(received, *input) {
			t.Errorf("Request body = %+v, want %+v", received, *input)
		}

		response := `{
			"incident": {
				"id": "01FDAG4SAP5TYPT98WGR2N7W91",
				"name": "New Database Issue",
				"summary": "Connection pool exhausted",
				"type": "incident",
				"status": "triage",
				"created_at": "2021-08-17T13:28:57.801578Z",
				"updated_at": "2021-08-17T13:28:57.801578Z",
				"mode": "real",
				"visibility": "public",
				"incident_role_assignments": [
					{
						"role": {
							"id": "01FH5TZRWMNAFB0DZ23FD1V96N",
							"name": "Incident Lead",
							"description": "Lead responder",
							"required": true,
							"created_at": "2021-08-17T13:28:57.801578Z",
							"updated_at": "2021-08-17T13:28:57.801578Z"
						},
						"assignee": {
							"id": "01FCNDV6P870EA6S7TK1DSYDG0",
							"name": "John Doe",
							"email": "john@example.com",
							"role": "owner"
						}
					}
				],
				"custom_field_values": {
					"01FH5TZRWMNAFB0DZ23FD1V96N": "high-priority"
				}
			}
		}`

		w.WriteHeader(http.StatusCreated)
		_, _ = fmt.Fprint(w, response)
	})

	ctx := context.Background()
	incident, resp, err := client.Incidents.Create(ctx, input)
	if err != nil {
		t.Errorf("Incidents.Create returned error: %v", err)
	}

	if resp.StatusCode != http.StatusCreated {
		t.Errorf("Incidents.Create returned status %d, want %d", resp.StatusCode, http.StatusCreated)
	}

	if incident.ID != "01FDAG4SAP5TYPT98WGR2N7W91" {
		t.Errorf("Incidents.Create returned ID %s, want %s", incident.ID, "01FDAG4SAP5TYPT98WGR2N7W91")
	}

	if incident.Name != "New Database Issue" {
		t.Errorf("Incidents.Create returned Name %s, want %s", incident.Name, "New Database Issue")
	}
}

func TestIncidentsService_Update(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	incidentID := "01FDAG4SAP5TYPT98WGR2N7W91"
	status := "resolved"
	summary := "Issue resolved"

	input := &UpdateIncidentOptions{
		Status:  &status,
		Summary: &summary,
	}

	mux.HandleFunc(fmt.Sprintf("/v2/incidents/%s", incidentID), func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "PUT")
		testHeader(t, r, "Content-Type", "application/json")

		var received UpdateIncidentOptions
		err := json.NewDecoder(r.Body).Decode(&received)
		if err != nil {
			t.Errorf("error decoding request body: %v", err)
		}

		if !reflect.DeepEqual(received, *input) {
			t.Errorf("Request body = %+v, want %+v", received, *input)
		}

		response := `{
			"incident": {
				"id": "01FDAG4SAP5TYPT98WGR2N7W91",
				"name": "Database Connection Issues",
				"summary": "Issue resolved",
				"type": "incident",
				"status": "resolved",
				"created_at": "2021-08-17T13:28:57.801578Z",
				"updated_at": "2021-08-17T14:28:57.801578Z",
				"closed_at": "2021-08-17T14:28:57.801578Z",
				"mode": "real",
				"visibility": "public"
			}
		}`

		_, _ = fmt.Fprint(w, response)
	})

	ctx := context.Background()
	incident, _, err := client.Incidents.Update(ctx, incidentID, input)
	if err != nil {
		t.Errorf("Incidents.Update returned error: %v", err)
	}

	if incident.Status != "resolved" {
		t.Errorf("Incidents.Update returned Status %s, want %s", incident.Status, "resolved")
	}

	if incident.Summary != "Issue resolved" {
		t.Errorf("Incidents.Update returned Summary %s, want %s", incident.Summary, "Issue resolved")
	}

	if incident.ClosedAt == nil {
		t.Error("Incidents.Update returned nil ClosedAt, want timestamp")
	}
}

func TestIncidentsService_Delete(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	incidentID := "01FDAG4SAP5TYPT98WGR2N7W91"

	mux.HandleFunc(fmt.Sprintf("/v2/incidents/%s", incidentID), func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "DELETE")
		w.WriteHeader(http.StatusNoContent)
	})

	ctx := context.Background()
	resp, err := client.Incidents.Delete(ctx, incidentID)
	if err != nil {
		t.Errorf("Incidents.Delete returned error: %v", err)
	}

	if resp.StatusCode != http.StatusNoContent {
		t.Errorf("Incidents.Delete returned status %d, want %d", resp.StatusCode, http.StatusNoContent)
	}
}

func TestIncidentsService_ErrorHandling(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/v2/incidents/not-found", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		response := `{
			"type": "validation_error",
			"status": 404,
			"detail": "Incident not found",
			"errors": [
				{
					"code": "not_found",
					"detail": "No incident with ID 'not-found'",
					"source": {
						"pointer": "/id"
					}
				}
			]
		}`
		_, _ = fmt.Fprint(w, response)
	})

	ctx := context.Background()
	_, _, err := client.Incidents.Get(ctx, "not-found")

	if err == nil {
		t.Error("Expected error, got nil")
	}

	errResp := &ErrorResponse{}
	ok := errors.As(err, &errResp)
	if !ok {
		t.Errorf("Error type = %T, want *ErrorResponse", err)
	}

	if errResp.Status != 404 {
		t.Errorf("Error status = %d, want %d", errResp.Status, 404)
	}

	if errResp.Detail != "Incident not found" {
		t.Errorf("Error detail = %s, want %s", errResp.Detail, "Incident not found")
	}

	if len(errResp.Errors) != 1 {
		t.Errorf("Error count = %d, want %d", len(errResp.Errors), 1)
	}
}

func TestTimestamp_MarshalJSON(t *testing.T) {
	tests := []struct {
		name string
		time time.Time
		want string
	}{
		{
			name: "regular timestamp",
			time: parseTime("2021-08-17T13:28:57.801578Z"),
			want: `"2021-08-17T13:28:57Z"`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ts := Timestamp{tt.time}
			got, err := json.Marshal(ts)
			if err != nil {
				t.Errorf("Timestamp.MarshalJSON() error = %v", err)
				return
			}
			if string(got) != tt.want {
				t.Errorf("Timestamp.MarshalJSON() = %s, want %s", string(got), tt.want)
			}
		})
	}
}

func TestTimestamp_UnmarshalJSON(t *testing.T) {
	tests := []struct {
		name    string
		json    string
		want    time.Time
		wantErr bool
	}{
		{
			name: "valid timestamp",
			json: `"2021-08-17T13:28:57.801578Z"`,
			want: parseTime("2021-08-17T13:28:57.801578Z"),
		},
		{
			name:    "invalid timestamp",
			json:    `"not-a-timestamp"`,
			wantErr: true,
		},
		{
			name:    "invalid json",
			json:    `not-json`,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var ts Timestamp
			err := json.Unmarshal([]byte(tt.json), &ts)
			if (err != nil) != tt.wantErr {
				t.Errorf("Timestamp.UnmarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && !ts.Equal(tt.want) {
				t.Errorf("Timestamp.UnmarshalJSON() = %v, want %v", ts.Time, tt.want)
			}
		})
	}
}

func TestErrorResponse_Error(t *testing.T) {
	req, _ := http.NewRequest(http.MethodGet, "https://api.incident.io/v2/incidents/123", nil)
	resp := &http.Response{
		Request:    req,
		StatusCode: http.StatusNotFound,
	}

	err := &ErrorResponse{
		Response: resp,
		Detail:   "Incident not found",
	}

	want := "GET https://api.incident.io/v2/incidents/123: 404 Incident not found"
	if got := err.Error(); got != want {
		t.Errorf("ErrorResponse.Error() = %q, want %q", got, want)
	}
}

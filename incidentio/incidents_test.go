package incidentio

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
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

func testQuery(t *testing.T, r *http.Request, want url.Values) {
	t.Helper()
	if got := r.URL.Query(); !reflect.DeepEqual(got, want) {
		t.Errorf("Query = %v, want %v", got, want)
	}
}

func testBody(t *testing.T, r *http.Request, want map[string]any) {
	t.Helper()
	raw, err := io.ReadAll(r.Body)
	if err != nil {
		t.Fatalf("reading body: %v", err)
	}
	var got map[string]any
	if err := json.Unmarshal(raw, &got); err != nil {
		t.Fatalf("decoding body %q: %v", raw, err)
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Body = %v, want %v", got, want)
	}
}

func parseTime(s string) time.Time {
	t, _ := time.Parse(time.RFC3339, s)
	return t
}

const incidentJSON = `{
	"id": "01FDAG4SAP5TYPT98WGR2N7W91",
	"reference": "INC-123",
	"name": "Database Connection Issues",
	"summary": "Users experiencing connection timeouts",
	"permalink": "https://app.incident.io/incidents/123",
	"mode": "standard",
	"visibility": "public",
	"incident_status": {
		"id": "status-1",
		"name": "Investigating",
		"category": "active",
		"rank": 2,
		"created_at": "2021-08-17T13:28:57.801578Z",
		"updated_at": "2021-08-17T13:28:57.801578Z"
	},
	"severity": {"id": "sev-1", "name": "Critical", "rank": 1},
	"incident_role_assignments": [
		{"role": {"id": "role-1", "name": "Lead"}, "assignee": {"id": "user-1", "name": "Jane"}}
	],
	"custom_field_entries": [
		{
			"custom_field": {"id": "cf-1", "name": "Team", "field_type": "single_select"},
			"values": [{"value_option": {"id": "opt-1", "custom_field_id": "cf-1", "value": "Platform", "sort_key": 1}}]
		}
	],
	"team_ids": ["team-1"],
	"creator": {"user": {"id": "user-1", "name": "Jane"}},
	"created_at": "2021-08-17T13:28:57.801578Z",
	"updated_at": "2021-08-17T14:28:57.801578Z",
	"last_activity_at": "2021-08-17T14:28:57.801578Z"
}`

func checkIncident(t *testing.T, got *Incident) {
	t.Helper()
	if got == nil {
		t.Fatal("incident is nil")
	}
	if got.ID != "01FDAG4SAP5TYPT98WGR2N7W91" || got.Reference != "INC-123" {
		t.Errorf("unexpected id/reference: %q %q", got.ID, got.Reference)
	}
	if got.IncidentStatus == nil || got.IncidentStatus.Category != "active" || got.IncidentStatus.Rank != 2 {
		t.Errorf("unexpected status: %+v", got.IncidentStatus)
	}
	if got.Severity == nil || got.Severity.ID != "sev-1" {
		t.Errorf("unexpected severity: %+v", got.Severity)
	}
	if len(got.CustomFieldEntries) != 1 ||
		got.CustomFieldEntries[0].CustomField.ID != "cf-1" ||
		got.CustomFieldEntries[0].Values[0].ValueOption.Value != "Platform" {
		t.Errorf("unexpected custom field entries: %+v", got.CustomFieldEntries)
	}
	if !reflect.DeepEqual(got.TeamIDs, []string{"team-1"}) {
		t.Errorf("TeamIDs = %v", got.TeamIDs)
	}
	if got.Creator == nil || got.Creator.User == nil || got.Creator.User.ID != "user-1" {
		t.Errorf("unexpected creator: %+v", got.Creator)
	}
	if !got.UpdatedAt.Equal(parseTime("2021-08-17T14:28:57.801578Z")) {
		t.Errorf("UpdatedAt = %v", got.UpdatedAt)
	}
}

func TestIncidentsService_List(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/v2/incidents", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testHeader(t, r, "Authorization", "Bearer test-key")
		testQuery(t, r, url.Values{})
		_, _ = fmt.Fprintf(w, `{"incidents": [%s]}`, incidentJSON)
	})

	got, _, err := client.Incidents.List(context.Background(), nil)
	if err != nil {
		t.Fatalf("Incidents.List returned error: %v", err)
	}
	if len(got) != 1 {
		t.Fatalf("got %d incidents, want 1", len(got))
	}
	checkIncident(t, got[0])
}

func TestIncidentsService_List_Options(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/v2/incidents", func(w http.ResponseWriter, r *http.Request) {
		testQuery(t, r, url.Values{
			"page_size":                     {"25"},
			"after":                         {"01ABC"},
			"sort_by":                       {"created_at_oldest_first"},
			"filter_mode":                   {"any"},
			"status[one_of]":                {"s1", "s2"},
			"status_category[one_of]":       {"active"},
			"created_at[gte]":               {"2024-01-01"},
			"severity[not_in]":              {"sev-9"},
			"mode[one_of]":                  {"standard"},
			"incident_role[role-1][one_of]": {"user-1"},
			"custom_field[cf-1][one_of]":    {"opt-1"},
		})
		_, _ = fmt.Fprint(w, `{"incidents": []}`)
	})

	opts := &IncidentListOptions{
		ListOptions:    ListOptions{PageSize: 25, After: "01ABC"},
		SortBy:         "created_at_oldest_first",
		FilterMode:     "any",
		Status:         Filter{"one_of": {"s1", "s2"}},
		StatusCategory: Filter{"one_of": {"active"}},
		CreatedAt:      Filter{"gte": {"2024-01-01"}},
		Severity:       Filter{"not_in": {"sev-9"}},
		Mode:           Filter{"one_of": {"standard"}},
		IncidentRole:   map[string]Filter{"role-1": {"one_of": {"user-1"}}},
		CustomField:    map[string]Filter{"cf-1": {"one_of": {"opt-1"}}},
	}
	got, _, err := client.Incidents.List(context.Background(), opts)
	if err != nil {
		t.Fatalf("Incidents.List returned error: %v", err)
	}
	if len(got) != 0 {
		t.Errorf("got %d incidents, want 0", len(got))
	}
}

func TestIncidentsService_Get(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/v2/incidents/01FDAG4SAP5TYPT98WGR2N7W91", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		_, _ = fmt.Fprintf(w, `{"incident": %s}`, incidentJSON)
	})

	got, _, err := client.Incidents.Get(context.Background(), "01FDAG4SAP5TYPT98WGR2N7W91")
	if err != nil {
		t.Fatalf("Incidents.Get returned error: %v", err)
	}
	checkIncident(t, got)
}

func TestIncidentsService_Create(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/v2/incidents", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
		testBody(t, r, map[string]any{
			"idempotency_key": "my-key",
			"visibility":      "public",
			"name":            "Database Connection Issues",
			"severity_id":     "sev-1",
			"incident_role_assignments": []any{
				map[string]any{
					"incident_role_id": "role-1",
					"assignee":         map[string]any{"id": "user-1"},
				},
			},
			"custom_field_entries": []any{
				map[string]any{
					"custom_field_id": "cf-1",
					"values":          []any{map[string]any{"value_option_id": "opt-1"}},
				},
			},
		})
		w.WriteHeader(http.StatusCreated)
		_, _ = fmt.Fprintf(w, `{"incident": %s}`, incidentJSON)
	})

	got, _, err := client.Incidents.Create(context.Background(), &CreateIncidentOptions{
		IdempotencyKey: "my-key",
		Visibility:     "public",
		Name:           "Database Connection Issues",
		SeverityID:     "sev-1",
		IncidentRoleAssignments: []CreateRoleAssignment{
			{IncidentRoleID: "role-1", Assignee: UserReference{ID: "user-1"}},
		},
		CustomFieldEntries: []CustomFieldEntryPayload{
			{CustomFieldID: "cf-1", Values: []CustomFieldValuePayload{{ValueOptionID: "opt-1"}}},
		},
	})
	if err != nil {
		t.Fatalf("Incidents.Create returned error: %v", err)
	}
	checkIncident(t, got)
}

func TestIncidentsService_Create_GeneratesIdempotencyKey(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	var keys []string
	mux.HandleFunc("/v2/incidents", func(w http.ResponseWriter, r *http.Request) {
		var body map[string]any
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Errorf("decoding body: %v", err)
		}
		key, _ := body["idempotency_key"].(string)
		keys = append(keys, key)
		_, _ = fmt.Fprintf(w, `{"incident": %s}`, incidentJSON)
	})

	opts := &CreateIncidentOptions{Visibility: "private", Name: "x"}
	for range 2 {
		if _, _, err := client.Incidents.Create(context.Background(), opts); err != nil {
			t.Fatalf("Incidents.Create returned error: %v", err)
		}
	}

	if opts.IdempotencyKey != "" {
		t.Errorf("caller options were mutated: %q", opts.IdempotencyKey)
	}
	if len(keys) != 2 || keys[0] == "" || keys[1] == "" {
		t.Fatalf("keys = %v, want two non-empty keys", keys)
	}
	if keys[0] == keys[1] {
		t.Errorf("expected distinct keys per call, got %q twice", keys[0])
	}
}

func TestIncidentsService_Update(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/v2/incidents/01FDAG4SAP5TYPT98WGR2N7W91/actions/edit", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
		testBody(t, r, map[string]any{
			"incident": map[string]any{
				"name":               "Renamed",
				"incident_status_id": "status-2",
			},
			"notify_incident_channel": true,
		})
		_, _ = fmt.Fprintf(w, `{"incident": %s}`, incidentJSON)
	})

	name, status := "Renamed", "status-2"
	got, _, err := client.Incidents.Update(context.Background(), "01FDAG4SAP5TYPT98WGR2N7W91", &UpdateIncidentOptions{
		Name:                  &name,
		IncidentStatusID:      &status,
		NotifyIncidentChannel: true,
	})
	if err != nil {
		t.Fatalf("Incidents.Update returned error: %v", err)
	}
	checkIncident(t, got)
}

func TestIncidentsService_ImportPostmortemDocument(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/v2/incidents/inc-1/actions/import_postmortem_document", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
		testBody(t, r, map[string]any{"title": "Postmortem", "content": "# Hello"})
		_, _ = fmt.Fprint(w, `{"postmortem_document": {"id": "pm-1", "incident_id": "inc-1", "title": "Postmortem", "document_url": "https://example.com/pm"}}`)
	})

	got, _, err := client.Incidents.ImportPostmortemDocument(context.Background(), "inc-1", "Postmortem", "# Hello")
	if err != nil {
		t.Fatalf("ImportPostmortemDocument returned error: %v", err)
	}
	if got.ID != "pm-1" || got.DocumentURL != "https://example.com/pm" {
		t.Errorf("unexpected document: %+v", got)
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
		t.Fatal("Expected error, got nil")
	}

	errResp := &ErrorResponse{}
	ok := errors.As(err, &errResp)
	if !ok {
		t.Fatalf("Error type = %T, want *ErrorResponse", err)
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

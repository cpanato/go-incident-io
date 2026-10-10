package incidentio

import (
	"context"
	"encoding/json"
	"net/http"
	"net/url"
	"testing"
)

func TestIncidentUpdatesService_List(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/v2/incident_updates", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testQuery(t, r, url.Values{"incident_id": {"inc1"}, "page_size": {"5"}, "after": {"u0"}})
		_, _ = w.Write([]byte(`{"incident_updates":[{"id":"u1","incident_id":"inc1","message":"hi",
			"new_severity":{"id":"s1","name":"Major"},"updater":{"user":{"id":"usr","name":"A"}},
			"created_at":"2024-01-01T00:00:00Z"}]}`))
	})

	got, _, err := client.IncidentUpdates.List(context.Background(), &IncidentUpdateListOptions{
		ListOptions: ListOptions{PageSize: 5, After: "u0"},
		IncidentID:  "inc1",
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != 1 || got[0].ID != "u1" || got[0].Message != "hi" || got[0].NewSeverity == nil || got[0].NewSeverity.ID != "s1" {
		t.Errorf("unexpected result: %+v", got)
	}
	if got[0].Updater == nil || got[0].Updater.User == nil || got[0].Updater.User.ID != "usr" {
		t.Errorf("updater not decoded: %+v", got[0].Updater)
	}
}

func TestIncidentUpdatesService_Create(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/v2/incident_updates", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
		testBody(t, r, map[string]any{
			"incident_id": "inc1", "idempotency_key": "k", "message": "m", "to_severity_id": "s2",
		})
		_, _ = w.Write([]byte(`{"incident_update":{"id":"u2","incident_id":"inc1"}}`))
	})

	got, _, err := client.IncidentUpdates.Create(context.Background(), &CreateIncidentUpdateOptions{
		IncidentID: "inc1", IdempotencyKey: "k", Message: "m", ToSeverityID: "s2",
	})
	if err != nil || got.ID != "u2" {
		t.Fatalf("got %+v, err %v", got, err)
	}
}

func TestIncidentUpdatesService_CreateGeneratesKey(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/v2/incident_updates", func(w http.ResponseWriter, r *http.Request) {
		var body map[string]any
		_ = json.NewDecoder(r.Body).Decode(&body)
		if k, _ := body["idempotency_key"].(string); k == "" {
			t.Error("expected generated idempotency key")
		}
		_, _ = w.Write([]byte(`{"incident_update":{"id":"u3"}}`))
	})

	opts := &CreateIncidentUpdateOptions{IncidentID: "inc1"}
	if _, _, err := client.IncidentUpdates.Create(context.Background(), opts); err != nil {
		t.Fatal(err)
	}
	if opts.IdempotencyKey != "" {
		t.Error("caller options must not be mutated")
	}
}

func TestIncidentTimelineItemsService(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/v2/incident_timeline_items", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			testQuery(t, r, url.Values{"incident_id": {"inc1"}})
			_, _ = w.Write([]byte(`{"incident_timeline_items":[{"id":"t1","incident_id":"inc1","title":"T",
				"timestamp":"2024-01-01T00:00:00Z","created_at":"2024-01-01T00:00:00Z","updated_at":"2024-01-01T00:00:00Z"}]}`))
		case http.MethodPost:
			testBody(t, r, map[string]any{
				"incident_id": "inc1", "title": "T", "timestamp": "2024-01-01T00:00:00Z",
				"idempotency_key": "k", "description": "d",
			})
			_, _ = w.Write([]byte(`{"incident_timeline_item":{"id":"t2","title":"T"}}`))
		}
	})

	items, _, err := client.IncidentTimelineItems.List(context.Background(), &IncidentTimelineItemListOptions{IncidentID: "inc1"})
	if err != nil || len(items) != 1 || items[0].ID != "t1" {
		t.Fatalf("got %+v, err %v", items, err)
	}

	item, _, err := client.IncidentTimelineItems.Create(context.Background(), &CreateIncidentTimelineItemOptions{
		IncidentID: "inc1", Title: "T", Timestamp: Timestamp{parseTime("2024-01-01T00:00:00Z")},
		IdempotencyKey: "k", Description: "d",
	})
	if err != nil || item.ID != "t2" {
		t.Fatalf("got %+v, err %v", item, err)
	}
}

func TestIncidentParticipantsService_List(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/v2/incident_participants", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testQuery(t, r, url.Values{"incident_id": {"inc1"}})
		_, _ = w.Write([]byte(`{"incident_participants":{
			"active":[{"participant_type":"responder","user":{"id":"u1","name":"A"}}],
			"passive":[{"participant_type":"observer","user":{"id":"u2","name":"B"}}]}}`))
	})

	got, _, err := client.IncidentParticipants.List(context.Background(), "inc1")
	if err != nil {
		t.Fatal(err)
	}
	if len(got.Active) != 1 || got.Active[0].ParticipantType != "responder" || got.Passive[0].User.ID != "u2" {
		t.Errorf("unexpected result: %+v", got)
	}
}

func TestUsersService_Get(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/v2/users/u1", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		_, _ = w.Write([]byte(`{"user":{"id":"u1","name":"A","email":"a@x.io","role":"owner","is_active":true,"slack_user_id":"U1"}}`))
	})

	got, _, err := client.Users.Get(context.Background(), "u1")
	if err != nil {
		t.Fatal(err)
	}
	if got.ID != "u1" || !got.IsActive || got.SlackUserID != "U1" {
		t.Errorf("unexpected user: %+v", got)
	}
}

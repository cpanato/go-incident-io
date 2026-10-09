package incidentio

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"testing"
	"time"
)

const overrideJSON = `{
	"id": "ovr-1",
	"schedule_id": "sch-1",
	"rotation_id": "rot-1",
	"layer_id": "lay-1",
	"user": {"id": "user-1", "name": "John", "email": "john@example.com"},
	"start_at": "2021-08-20T09:00:00Z",
	"end_at": "2021-08-20T17:00:00Z",
	"created_at": "2021-08-17T13:28:57.801578Z",
	"updated_at": "2021-08-17T13:28:57.801578Z"
}`

func checkOverride(t *testing.T, got *Override) {
	t.Helper()
	if got == nil {
		t.Fatal("override is nil")
	}
	if got.ID != "ovr-1" || got.ScheduleID != "sch-1" || got.RotationID != "rot-1" || got.LayerID != "lay-1" {
		t.Errorf("unexpected override: %+v", got)
	}
	if got.User == nil || got.User.Email != "john@example.com" {
		t.Errorf("unexpected user: %+v", got.User)
	}
	if !got.StartAt.Equal(parseTime("2021-08-20T09:00:00Z")) || !got.EndAt.Equal(parseTime("2021-08-20T17:00:00Z")) {
		t.Errorf("unexpected times: %v - %v", got.StartAt, got.EndAt)
	}
}

func TestSchedulesService_ListEntries(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	start := time.Date(2021, 8, 17, 0, 0, 0, 0, time.UTC)
	end := time.Date(2021, 8, 24, 0, 0, 0, 0, time.UTC)

	mux.HandleFunc("/v2/schedule_entries", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testQuery(t, r, url.Values{
			"schedule_id":        {"sch-1"},
			"entry_window_start": {"2021-08-17T00:00:00Z"},
			"entry_window_end":   {"2021-08-24T00:00:00Z"},
		})
		_, _ = fmt.Fprint(w, `{"schedule_entries": {
			"scheduled": [{"entry_id": "e1", "fingerprint": "f1", "rotation_id": "rot-1",
				"start_at": "2021-08-17T13:00:00Z", "end_at": "2021-08-24T13:00:00Z",
				"user": {"id": "user-1", "name": "John"}}],
			"overrides": [],
			"final": [{"entry_id": "e2", "rotation_id": "rot-1",
				"start_at": "2021-08-17T13:00:00Z", "end_at": "2021-08-24T13:00:00Z"}]
		}}`)
	})

	got, _, err := client.Schedules.ListEntries(context.Background(), "sch-1", &ScheduleEntriesOptions{
		EntryWindow: &TimeWindow{StartAt: start, EndAt: end},
	})
	if err != nil {
		t.Fatalf("Schedules.ListEntries returned error: %v", err)
	}
	if len(got.Scheduled) != 1 || got.Scheduled[0].EntryID != "e1" || got.Scheduled[0].User.ID != "user-1" {
		t.Errorf("unexpected scheduled entries: %+v", got.Scheduled)
	}
	if len(got.Overrides) != 0 {
		t.Errorf("Overrides = %d, want 0", len(got.Overrides))
	}
	if len(got.Final) != 1 || got.Final[0].EntryID != "e2" {
		t.Errorf("unexpected final entries: %+v", got.Final)
	}
}

func TestSchedulesService_ListEntries_NoWindow(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/v2/schedule_entries", func(w http.ResponseWriter, r *http.Request) {
		testQuery(t, r, url.Values{"schedule_id": {"sch-1"}})
		_, _ = fmt.Fprint(w, `{"schedule_entries": {"scheduled": [], "overrides": [], "final": []}}`)
	})

	if _, _, err := client.Schedules.ListEntries(context.Background(), "sch-1", nil); err != nil {
		t.Fatalf("Schedules.ListEntries returned error: %v", err)
	}
}

func TestSchedulesService_ListOverrides(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/v2/schedule_overrides", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testQuery(t, r, url.Values{"schedule_id": {"sch-1"}})
		_, _ = fmt.Fprintf(w, `{"overrides": [%s]}`, overrideJSON)
	})

	got, _, err := client.Schedules.ListOverrides(context.Background(), "sch-1", nil)
	if err != nil {
		t.Fatalf("Schedules.ListOverrides returned error: %v", err)
	}
	if len(got) != 1 {
		t.Fatalf("got %d overrides, want 1", len(got))
	}
	checkOverride(t, got[0])
}

func TestSchedulesService_ListOverrides_Options(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/v2/schedule_overrides", func(w http.ResponseWriter, r *http.Request) {
		testQuery(t, r, url.Values{
			"schedule_id": {"sch-1"},
			"rotation_id": {"rot-1"},
			"layer_id":    {"lay-1"},
			"page_size":   {"10"},
			"after":       {"ovr-0"},
		})
		_, _ = fmt.Fprint(w, `{"overrides": []}`)
	})

	_, _, err := client.Schedules.ListOverrides(context.Background(), "sch-1", &ListOverridesOptions{
		ListOptions: ListOptions{PageSize: 10, After: "ovr-0"},
		RotationID:  "rot-1",
		LayerID:     "lay-1",
	})
	if err != nil {
		t.Fatalf("Schedules.ListOverrides returned error: %v", err)
	}
}

func TestSchedulesService_GetOverride(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/v2/schedule_overrides/ovr-1", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		_, _ = fmt.Fprintf(w, `{"override": %s}`, overrideJSON)
	})

	got, _, err := client.Schedules.GetOverride(context.Background(), "ovr-1")
	if err != nil {
		t.Fatalf("Schedules.GetOverride returned error: %v", err)
	}
	checkOverride(t, got)
}

func TestSchedulesService_CreateOverride(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/v2/schedule_overrides", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
		testBody(t, r, map[string]any{
			"schedule_id": "sch-1",
			"rotation_id": "rot-1",
			"layer_id":    "lay-1",
			"user":        map[string]any{"email": "john@example.com"},
			"start_at":    "2021-08-20T09:00:00Z",
			"end_at":      "2021-08-20T17:00:00Z",
		})
		w.WriteHeader(http.StatusCreated)
		_, _ = fmt.Fprintf(w, `{"override": %s}`, overrideJSON)
	})

	got, resp, err := client.Schedules.CreateOverride(context.Background(), &CreateOverrideOptions{
		ScheduleID: "sch-1",
		RotationID: "rot-1",
		LayerID:    "lay-1",
		User:       UserReference{Email: "john@example.com"},
		StartAt:    Timestamp{parseTime("2021-08-20T09:00:00Z")},
		EndAt:      Timestamp{parseTime("2021-08-20T17:00:00Z")},
	})
	if err != nil {
		t.Fatalf("Schedules.CreateOverride returned error: %v", err)
	}
	if resp.StatusCode != http.StatusCreated {
		t.Errorf("status = %d, want %d", resp.StatusCode, http.StatusCreated)
	}
	checkOverride(t, got)
}

func TestSchedulesService_UpdateOverride(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/v2/schedule_overrides/ovr-1", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "PUT")
		testBody(t, r, map[string]any{
			"rotation_id": "rot-1",
			"layer_id":    "lay-1",
			"user":        map[string]any{"id": "user-1"},
			"start_at":    "2021-08-20T09:00:00Z",
			"end_at":      "2021-08-20T17:00:00Z",
		})
		_, _ = fmt.Fprintf(w, `{"override": %s}`, overrideJSON)
	})

	got, _, err := client.Schedules.UpdateOverride(context.Background(), "ovr-1", &UpdateOverrideOptions{
		RotationID: "rot-1",
		LayerID:    "lay-1",
		User:       UserReference{ID: "user-1"},
		StartAt:    Timestamp{parseTime("2021-08-20T09:00:00Z")},
		EndAt:      Timestamp{parseTime("2021-08-20T17:00:00Z")},
	})
	if err != nil {
		t.Fatalf("Schedules.UpdateOverride returned error: %v", err)
	}
	checkOverride(t, got)
}

func TestSchedulesService_DeleteOverride(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/v2/schedule_overrides/ovr-1", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "DELETE")
		w.WriteHeader(http.StatusNoContent)
	})

	resp, err := client.Schedules.DeleteOverride(context.Background(), "ovr-1")
	if err != nil {
		t.Fatalf("Schedules.DeleteOverride returned error: %v", err)
	}
	if resp.StatusCode != http.StatusNoContent {
		t.Errorf("status = %d, want %d", resp.StatusCode, http.StatusNoContent)
	}
}

func TestSchedulesService_List_Pagination(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/v2/schedules", func(w http.ResponseWriter, r *http.Request) {
		testQuery(t, r, url.Values{"page_size": {"5"}, "after": {"sch-0"}})
		_, _ = fmt.Fprint(w, `{"schedules": []}`)
	})

	_, _, err := client.Schedules.List(context.Background(), &ScheduleListOptions{
		ListOptions: ListOptions{PageSize: 5, After: "sch-0"},
	})
	if err != nil {
		t.Fatalf("Schedules.List returned error: %v", err)
	}
}

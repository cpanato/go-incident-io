package incidentio

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"
)

const alertSourceJSON = `{
	"id": "src1",
	"name": "Datadog",
	"source_type": "datadog",
	"secret_token": "tok",
	"alert_events_url": "https://api.incident.io/v2/alert_events/datadog/src1",
	"auto_resolve_incident_alerts": true,
	"auto_resolve_timeout_minutes": 30,
	"owning_team_ids": ["team1"],
	"template": {"title": {"literal": "x"}}
}`

func TestAlertSourcesService_ListGet(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/v2/alert_sources", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		_, _ = w.Write([]byte(`{"alert_sources":[` + alertSourceJSON + `]}`))
	})
	mux.HandleFunc("/v2/alert_sources/src1", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		_, _ = w.Write([]byte(`{"alert_source":` + alertSourceJSON + `}`))
	})

	list, _, err := client.AlertSources.List(context.Background())
	if err != nil || len(list) != 1 {
		t.Fatalf("got %+v, err %v", list, err)
	}
	s := list[0]
	if s.SourceType != "datadog" || !s.AutoResolveIncidentAlerts || s.AutoResolveTimeoutMinutes != 30 || s.OwningTeamIDs[0] != "team1" {
		t.Errorf("unexpected source: %+v", s)
	}
	var tpl map[string]any
	if err := json.Unmarshal(s.Template, &tpl); err != nil || tpl["title"] == nil {
		t.Errorf("template not preserved: %s (%v)", s.Template, err)
	}

	got, _, err := client.AlertSources.Get(context.Background(), "src1")
	if err != nil || got.ID != "src1" {
		t.Fatalf("got %+v, err %v", got, err)
	}
}

func TestAlertSourcesService_Create(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/v2/alert_sources", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
		testBody(t, r, map[string]any{
			"name":                    "Datadog",
			"source_type":             "datadog",
			"template":                map[string]any{"title": map[string]any{"literal": "x"}},
			"filter_condition_groups": []any{},
		})
		_, _ = w.Write([]byte(`{"alert_source":` + alertSourceJSON + `}`))
	})

	got, _, err := client.AlertSources.Create(context.Background(), &CreateAlertSourceOptions{
		Name:                  "Datadog",
		SourceType:            "datadog",
		Template:              json.RawMessage(`{"title":{"literal":"x"}}`),
		FilterConditionGroups: json.RawMessage(`[]`),
	})
	if err != nil || got.ID != "src1" {
		t.Fatalf("got %+v, err %v", got, err)
	}
}

func TestAlertSourcesService_UpdateDelete(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/v2/alert_sources/src1", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPut:
			testBody(t, r, map[string]any{
				"name":     "Renamed",
				"template": map[string]any{},
				"disabled": true,
			})
			_, _ = w.Write([]byte(`{"alert_source":` + alertSourceJSON + `}`))
		case http.MethodDelete:
			w.WriteHeader(http.StatusNoContent)
		default:
			t.Errorf("unexpected method %s", r.Method)
		}
	})

	disabled := true
	got, _, err := client.AlertSources.Update(context.Background(), "src1", &UpdateAlertSourceOptions{
		Name: "Renamed", Template: json.RawMessage(`{}`), Disabled: &disabled,
	})
	if err != nil || got.ID != "src1" {
		t.Fatalf("got %+v, err %v", got, err)
	}

	resp, err := client.AlertSources.Delete(context.Background(), "src1")
	if err != nil || resp.StatusCode != http.StatusNoContent {
		t.Fatalf("resp %v, err %v", resp, err)
	}
}

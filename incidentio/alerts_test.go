package incidentio

import (
	"context"
	"net/http"
	"net/url"
	"testing"
)

const alertJSON = `{
	"id": "al1",
	"alert_source_id": "src1",
	"deduplication_key": "dk",
	"status": "firing",
	"title": "CPU high",
	"source_url": "https://example.com",
	"alert_group_ids": ["g1"],
	"attributes": [
		{"attribute": {"id": "a1", "name": "Team", "type": "String", "array": false, "required": true},
		 "value": {"literal": "infra", "label": "Infra"}},
		{"attribute": {"id": "a2", "name": "Hosts", "type": "String", "array": true, "required": false},
		 "array_value": [{"literal": "h1"}, {"literal": "h2"}]}
	],
	"tags": [{"id": "t1", "name": "prod"}],
	"created_at": "2024-01-01T00:00:00Z",
	"updated_at": "2024-01-02T00:00:00Z",
	"resolved_at": "2024-01-03T00:00:00Z"
}`

func TestAlertsService_List(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/v2/alerts", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testQuery(t, r, url.Values{
			"page_size":            {"10"},
			"status[one_of]":       {"firing"},
			"alert_source[not_in]": {"s1", "s2"},
			"created_at[gte]":      {"2024-01-01"},
		})
		_, _ = w.Write([]byte(`{"alerts":[` + alertJSON + `]}`))
	})

	got, _, err := client.Alerts.List(context.Background(), &AlertListOptions{
		ListOptions: ListOptions{PageSize: 10},
		Status:      Filter{"one_of": {"firing"}},
		AlertSource: Filter{"not_in": {"s1", "s2"}},
		CreatedAt:   Filter{"gte": {"2024-01-01"}},
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != 1 {
		t.Fatalf("got %d alerts", len(got))
	}
	a := got[0]
	if a.ID != "al1" || a.Status != "firing" || a.Title != "CPU high" || a.AlertGroupIDs[0] != "g1" {
		t.Errorf("unexpected alert: %+v", a)
	}
	if a.Attributes[0].Value.Literal != "infra" || len(a.Attributes[1].ArrayValue) != 2 {
		t.Errorf("attributes not decoded: %+v", a.Attributes)
	}
	if a.Tags[0].Name != "prod" || a.ResolvedAt == nil {
		t.Errorf("tags/resolved_at not decoded: %+v", a)
	}
}

func TestAlertsService_Get(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/v2/alerts/al1", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		_, _ = w.Write([]byte(`{"alert":` + alertJSON + `}`))
	})

	got, _, err := client.Alerts.Get(context.Background(), "al1")
	if err != nil || got.ID != "al1" {
		t.Fatalf("got %+v, err %v", got, err)
	}
}

func TestAlertsService_Actions(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/v2/alerts/al1/actions/resolve", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
		testBody(t, r, map[string]any{})
		_, _ = w.Write([]byte(`{"alert":` + alertJSON + `}`))
	})
	for _, action := range []string{"add_tags", "remove_tags", "set_tags"} {
		mux.HandleFunc("/v2/alerts/al1/actions/"+action, func(w http.ResponseWriter, r *http.Request) {
			testMethod(t, r, "POST")
			testBody(t, r, map[string]any{"tags": []any{"a", "b"}})
			_, _ = w.Write([]byte(`{"alert":` + alertJSON + `}`))
		})
	}

	ctx := context.Background()
	calls := map[string]func() (*Alert, error){
		"resolve": func() (*Alert, error) { a, _, err := client.Alerts.Resolve(ctx, "al1"); return a, err },
		"add": func() (*Alert, error) {
			a, _, err := client.Alerts.AddTags(ctx, "al1", []string{"a", "b"})
			return a, err
		},
		"remove": func() (*Alert, error) {
			a, _, err := client.Alerts.RemoveTags(ctx, "al1", []string{"a", "b"})
			return a, err
		},
		"set": func() (*Alert, error) {
			a, _, err := client.Alerts.SetTags(ctx, "al1", []string{"a", "b"})
			return a, err
		},
	}
	for name, call := range calls {
		a, err := call()
		if err != nil || a.ID != "al1" {
			t.Errorf("%s: got %+v, err %v", name, a, err)
		}
	}
}

func TestAlertsService_SetTagsNilSendsEmptyArray(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/v2/alerts/al1/actions/set_tags", func(w http.ResponseWriter, r *http.Request) {
		testBody(t, r, map[string]any{"tags": []any{}})
		_, _ = w.Write([]byte(`{"alert":` + alertJSON + `}`))
	})

	if _, _, err := client.Alerts.SetTags(context.Background(), "al1", nil); err != nil {
		t.Fatal(err)
	}
}

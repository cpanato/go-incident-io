package incidentio

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"testing"
)

const actionJSON = `{
	"id": "act-1",
	"incident_id": "inc-1",
	"description": "Fix it",
	"status": "outstanding",
	"assignee": {"id": "user-1", "name": "Jane"},
	"creator": {"user": {"id": "user-2", "name": "Bob"}},
	"created_at": "2021-08-17T13:28:57.801578Z",
	"updated_at": "2021-08-17T13:28:57.801578Z"
}`

func checkAction(t *testing.T, got *Action) {
	t.Helper()
	if got == nil {
		t.Fatal("action is nil")
	}
	if got.ID != "act-1" || got.IncidentID != "inc-1" || got.Status != "outstanding" {
		t.Errorf("unexpected action: %+v", got)
	}
	if got.Assignee == nil || got.Assignee.ID != "user-1" {
		t.Errorf("unexpected assignee: %+v", got.Assignee)
	}
	if got.Creator == nil || got.Creator.User == nil || got.Creator.User.ID != "user-2" {
		t.Errorf("unexpected creator: %+v", got.Creator)
	}
	if got.CompletedAt != nil {
		t.Errorf("CompletedAt = %v, want nil", got.CompletedAt)
	}
}

func TestActionsService_List(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/v3/actions", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testQuery(t, r, url.Values{
			"incident_id":     {"inc-1"},
			"incident_mode":   {"standard"},
			"page_size":       {"20"},
			"after":           {"act-0"},
			"created_at[gte]": {"2024-01-01"},
			"updated_at[lte]": {"2024-02-01"},
		})
		_, _ = fmt.Fprintf(w, `{"actions": [%s]}`, actionJSON)
	})

	got, _, err := client.Actions.List(context.Background(), &ActionListOptions{
		ListOptions:  ListOptions{PageSize: 20, After: "act-0"},
		IncidentID:   "inc-1",
		IncidentMode: "standard",
		CreatedAt:    Filter{"gte": {"2024-01-01"}},
		UpdatedAt:    Filter{"lte": {"2024-02-01"}},
	})
	if err != nil {
		t.Fatalf("Actions.List returned error: %v", err)
	}
	if len(got) != 1 {
		t.Fatalf("got %d actions, want 1", len(got))
	}
	checkAction(t, got[0])
}

func TestActionsService_Get(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/v3/actions/act-1", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		_, _ = fmt.Fprintf(w, `{"action": %s}`, actionJSON)
	})

	got, _, err := client.Actions.Get(context.Background(), "act-1")
	if err != nil {
		t.Fatalf("Actions.Get returned error: %v", err)
	}
	checkAction(t, got)
}

func TestActionsService_Create(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/v3/actions", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
		testBody(t, r, map[string]any{
			"incident_id": "inc-1",
			"description": "Fix it",
			"assignee_id": "user-1",
		})
		w.WriteHeader(http.StatusCreated)
		_, _ = fmt.Fprintf(w, `{"action": %s}`, actionJSON)
	})

	got, _, err := client.Actions.Create(context.Background(), &CreateActionOptions{
		IncidentID:  "inc-1",
		Description: "Fix it",
		AssigneeID:  "user-1",
	})
	if err != nil {
		t.Fatalf("Actions.Create returned error: %v", err)
	}
	checkAction(t, got)
}

func TestActionsService_Update(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/v3/actions/act-1", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "PUT")
		testBody(t, r, map[string]any{"description": "Fix it", "status": "completed"})
		_, _ = fmt.Fprintf(w, `{"action": %s}`, actionJSON)
	})

	got, _, err := client.Actions.Update(context.Background(), "act-1", &UpdateActionOptions{
		Description: "Fix it",
		Status:      "completed",
	})
	if err != nil {
		t.Fatalf("Actions.Update returned error: %v", err)
	}
	checkAction(t, got)
}

func TestActionsService_Delete(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/v3/actions/act-1", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "DELETE")
		w.WriteHeader(http.StatusNoContent)
	})

	resp, err := client.Actions.Delete(context.Background(), "act-1")
	if err != nil {
		t.Fatalf("Actions.Delete returned error: %v", err)
	}
	if resp.StatusCode != http.StatusNoContent {
		t.Errorf("status = %d, want %d", resp.StatusCode, http.StatusNoContent)
	}
}

const customFieldJSON = `{
	"id": "cf-1",
	"name": "Team",
	"description": "Owning team",
	"field_type": "single_select",
	"catalog_type_id": "cat-1",
	"created_at": "2021-08-17T13:28:57.801578Z",
	"updated_at": "2021-08-17T13:28:57.801578Z"
}`

func checkCustomField(t *testing.T, got *CustomField) {
	t.Helper()
	if got == nil {
		t.Fatal("custom field is nil")
	}
	if got.ID != "cf-1" || got.Name != "Team" || got.FieldType != "single_select" || got.CatalogTypeID != "cat-1" {
		t.Errorf("unexpected custom field: %+v", got)
	}
}

func TestCustomFieldsService_List(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/v2/custom_fields", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		_, _ = fmt.Fprintf(w, `{"custom_fields": [%s]}`, customFieldJSON)
	})

	got, _, err := client.CustomFields.List(context.Background())
	if err != nil {
		t.Fatalf("CustomFields.List returned error: %v", err)
	}
	if len(got) != 1 {
		t.Fatalf("got %d custom fields, want 1", len(got))
	}
	checkCustomField(t, got[0])
}

func TestCustomFieldsService_Get(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/v2/custom_fields/cf-1", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		_, _ = fmt.Fprintf(w, `{"custom_field": %s}`, customFieldJSON)
	})

	got, _, err := client.CustomFields.Get(context.Background(), "cf-1")
	if err != nil {
		t.Fatalf("CustomFields.Get returned error: %v", err)
	}
	checkCustomField(t, got)
}

func TestCustomFieldsService_Create(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/v2/custom_fields", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
		testBody(t, r, map[string]any{
			"name":            "Team",
			"description":     "Owning team",
			"field_type":      "single_select",
			"catalog_type_id": "cat-1",
		})
		w.WriteHeader(http.StatusCreated)
		_, _ = fmt.Fprintf(w, `{"custom_field": %s}`, customFieldJSON)
	})

	got, _, err := client.CustomFields.Create(context.Background(), &CreateCustomFieldOptions{
		Name:          "Team",
		Description:   "Owning team",
		FieldType:     "single_select",
		CatalogTypeID: "cat-1",
	})
	if err != nil {
		t.Fatalf("CustomFields.Create returned error: %v", err)
	}
	checkCustomField(t, got)
}

func TestCustomFieldsService_Update(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/v2/custom_fields/cf-1", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "PUT")
		testBody(t, r, map[string]any{"name": "Team", "description": "Owning team"})
		_, _ = fmt.Fprintf(w, `{"custom_field": %s}`, customFieldJSON)
	})

	got, _, err := client.CustomFields.Update(context.Background(), "cf-1", &UpdateCustomFieldOptions{
		Name:        "Team",
		Description: "Owning team",
	})
	if err != nil {
		t.Fatalf("CustomFields.Update returned error: %v", err)
	}
	checkCustomField(t, got)
}

func TestCustomFieldsService_Delete(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/v2/custom_fields/cf-1", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "DELETE")
		w.WriteHeader(http.StatusNoContent)
	})

	resp, err := client.CustomFields.Delete(context.Background(), "cf-1")
	if err != nil {
		t.Fatalf("CustomFields.Delete returned error: %v", err)
	}
	if resp.StatusCode != http.StatusNoContent {
		t.Errorf("status = %d, want %d", resp.StatusCode, http.StatusNoContent)
	}
}

const workflowJSON = `{
	"id": "wf-1",
	"name": "Notify on create",
	"version": 3,
	"state": "active",
	"trigger": {"name": "incident.updated", "label": "Incident updated"},
	"once_for": [{"key": "incident.id"}],
	"runs_on_incidents": "newly_created",
	"runs_on_incident_modes": ["standard"],
	"continue_on_step_error": true,
	"condition_groups": [],
	"expressions": [],
	"steps": [{"name": "slack.post"}]
}`

func checkWorkflow(t *testing.T, got *Workflow) {
	t.Helper()
	if got == nil {
		t.Fatal("workflow is nil")
	}
	if got.ID != "wf-1" || got.Version != 3 || got.State != "active" || !got.ContinueOnStepError {
		t.Errorf("unexpected workflow: %+v", got)
	}
	if got.Trigger == nil || got.Trigger.Name != "incident.updated" {
		t.Errorf("unexpected trigger: %+v", got.Trigger)
	}
	if len(got.OnceFor) != 1 {
		t.Errorf("OnceFor = %v", got.OnceFor)
	}
	if string(got.Steps) != `[{"name": "slack.post"}]` {
		t.Errorf("Steps = %s", got.Steps)
	}
}

func TestWorkflowsService_List(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/v2/workflows", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		_, _ = fmt.Fprint(w, `{"workflows": [{"id": "wf-1", "name": "Notify on create", "version": 3, "state": "active"}]}`)
	})

	got, _, err := client.Workflows.List(context.Background())
	if err != nil {
		t.Fatalf("Workflows.List returned error: %v", err)
	}
	if len(got) != 1 || got[0].ID != "wf-1" {
		t.Errorf("unexpected workflows: %+v", got)
	}
}

func TestWorkflowsService_Get(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/v2/workflows/wf-1", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testQuery(t, r, url.Values{})
		_, _ = fmt.Fprintf(w, `{"workflow": %s}`, workflowJSON)
	})

	got, _, err := client.Workflows.Get(context.Background(), "wf-1", false)
	if err != nil {
		t.Fatalf("Workflows.Get returned error: %v", err)
	}
	checkWorkflow(t, got)
}

func TestWorkflowsService_Get_SkipStepUpgrades(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/v2/workflows/wf-1", func(w http.ResponseWriter, r *http.Request) {
		testQuery(t, r, url.Values{"skip_step_upgrades": {"true"}})
		_, _ = fmt.Fprintf(w, `{"workflow": %s}`, workflowJSON)
	})

	if _, _, err := client.Workflows.Get(context.Background(), "wf-1", true); err != nil {
		t.Fatalf("Workflows.Get returned error: %v", err)
	}
}

func TestWorkflowsService_Create(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/v2/workflows", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
		testBody(t, r, map[string]any{
			"name":                   "Notify on create",
			"trigger":                "incident.updated",
			"once_for":               []any{"incident.id"},
			"condition_groups":       []any{},
			"steps":                  []any{map[string]any{"name": "slack.post"}},
			"expressions":            []any{},
			"runs_on_incident_modes": []any{"standard"},
			"runs_on_incidents":      "newly_created",
			"continue_on_step_error": true,
		})
		w.WriteHeader(http.StatusCreated)
		_, _ = fmt.Fprintf(w, `{"workflow": %s}`, workflowJSON)
	})

	got, _, err := client.Workflows.Create(context.Background(), &CreateWorkflowOptions{
		Name:                "Notify on create",
		Trigger:             "incident.updated",
		OnceFor:             []string{"incident.id"},
		ConditionGroups:     json.RawMessage(`[]`),
		Steps:               json.RawMessage(`[{"name":"slack.post"}]`),
		Expressions:         json.RawMessage(`[]`),
		RunsOnIncidentModes: []string{"standard"},
		RunsOnIncidents:     "newly_created",
		ContinueOnStepError: true,
	})
	if err != nil {
		t.Fatalf("Workflows.Create returned error: %v", err)
	}
	checkWorkflow(t, got)
}

func TestWorkflowsService_Update(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/v2/workflows/wf-1", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "PUT")
		testBody(t, r, map[string]any{
			"name":                   "Renamed",
			"once_for":               []any{},
			"condition_groups":       []any{},
			"steps":                  []any{},
			"expressions":            []any{},
			"runs_on_incident_modes": []any{"standard"},
			"runs_on_incidents":      "newly_created",
			"continue_on_step_error": false,
			"state":                  "disabled",
		})
		_, _ = fmt.Fprintf(w, `{"workflow": %s}`, workflowJSON)
	})

	got, _, err := client.Workflows.Update(context.Background(), "wf-1", &UpdateWorkflowOptions{
		Name:                "Renamed",
		OnceFor:             []string{},
		ConditionGroups:     json.RawMessage(`[]`),
		Steps:               json.RawMessage(`[]`),
		Expressions:         json.RawMessage(`[]`),
		RunsOnIncidentModes: []string{"standard"},
		RunsOnIncidents:     "newly_created",
		State:               "disabled",
	})
	if err != nil {
		t.Fatalf("Workflows.Update returned error: %v", err)
	}
	checkWorkflow(t, got)
}

func TestWorkflowsService_Delete(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/v2/workflows/wf-1", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "DELETE")
		w.WriteHeader(http.StatusNoContent)
	})

	resp, err := client.Workflows.Delete(context.Background(), "wf-1")
	if err != nil {
		t.Fatalf("Workflows.Delete returned error: %v", err)
	}
	if resp.StatusCode != http.StatusNoContent {
		t.Errorf("status = %d, want %d", resp.StatusCode, http.StatusNoContent)
	}
}

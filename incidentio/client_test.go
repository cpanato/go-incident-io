package incidentio

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestNewClient_Defaults(t *testing.T) {
	c := NewClient("key")
	if got, want := c.BaseURL.String(), defaultBaseURL; got != want {
		t.Errorf("BaseURL = %q, want %q", got, want)
	}
	if c.UserAgent != userAgent {
		t.Errorf("UserAgent = %q, want %q", c.UserAgent, userAgent)
	}
	for name, svc := range map[string]any{
		"Incidents": c.Incidents, "Severities": c.Severities, "IncidentTypes": c.IncidentTypes,
		"IncidentRoles": c.IncidentRoles, "CustomFields": c.CustomFields, "Actions": c.Actions,
		"Workflows": c.Workflows, "Schedules": c.Schedules, "Users": c.Users, "Webhooks": c.Webhooks,
	} {
		if svc == nil {
			t.Errorf("service %s is nil", name)
		}
	}
}

func TestNewClient_Options(t *testing.T) {
	hc := &http.Client{Timeout: time.Second}
	c := NewClient("key", WithHTTPClient(hc), WithBaseURL("https://example.com/api/"))

	if c.client != hc {
		t.Error("WithHTTPClient did not set the HTTP client")
	}
	if got := c.BaseURL.String(); got != "https://example.com/api/" {
		t.Errorf("BaseURL = %q", got)
	}
}

func TestClient_NewRequest(t *testing.T) {
	c := NewClient("secret")

	t.Run("without body", func(t *testing.T) {
		req, err := c.NewRequest("GET", "v2/incidents", nil)
		if err != nil {
			t.Fatal(err)
		}
		if got := req.URL.String(); got != "https://api.incident.io/v2/incidents" {
			t.Errorf("URL = %q", got)
		}
		if got := req.Header.Get("Content-Type"); got != "" {
			t.Errorf("Content-Type = %q, want empty", got)
		}
		if got := req.Header.Get("Authorization"); got != "Bearer secret" {
			t.Errorf("Authorization = %q", got)
		}
		if got := req.Header.Get("Accept"); got != "application/json" {
			t.Errorf("Accept = %q", got)
		}
	})

	t.Run("with body", func(t *testing.T) {
		req, err := c.NewRequest("POST", "v2/incidents", map[string]string{"a": "<b>"})
		if err != nil {
			t.Fatal(err)
		}
		if got := req.Header.Get("Content-Type"); got != "application/json" {
			t.Errorf("Content-Type = %q", got)
		}
	})

	t.Run("unencodable body", func(t *testing.T) {
		if _, err := c.NewRequest("POST", "v2/incidents", make(chan int)); err == nil {
			t.Error("expected error")
		}
	})

	t.Run("bad url", func(t *testing.T) {
		if _, err := c.NewRequest("GET", "http://[::1", nil); err == nil {
			t.Error("expected error")
		}
	})
}

func TestClient_Do_ContextCanceled(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/v2/incidents", func(w http.ResponseWriter, _ *http.Request) {
		_, _ = fmt.Fprint(w, `{"incidents": []}`)
	})

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	_, _, err := client.Incidents.List(ctx, nil)
	if !errors.Is(err, context.Canceled) {
		t.Errorf("error = %v, want context.Canceled", err)
	}
}

func TestClient_Do_NoContentWithTarget(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/v2/custom_fields", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	})

	got, _, err := client.CustomFields.List(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got) != 0 {
		t.Errorf("got %d fields, want 0", len(got))
	}
}

func TestSeveritiesService_List(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/v1/severities", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		_, _ = fmt.Fprint(w, `{"severities": [{"id": "sev-1", "name": "Critical", "rank": 1}]}`)
	})

	got, _, err := client.Severities.List(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != 1 || got[0].ID != "sev-1" || got[0].Rank != 1 {
		t.Errorf("unexpected severities: %+v", got)
	}
}

func TestIncidentTypesService_List(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/v1/incident_types", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		_, _ = fmt.Fprint(w, `{"incident_types": [{"id": "type-1", "name": "Outage", "description": "d"}]}`)
	})

	got, _, err := client.IncidentTypes.List(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != 1 || got[0].ID != "type-1" || got[0].Name != "Outage" {
		t.Errorf("unexpected incident types: %+v", got)
	}
}

func TestIncidentRolesService_List(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/v2/incident_roles", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		_, _ = fmt.Fprint(w, `{"incident_roles": [{"id": "role-1", "name": "Lead", "required": true}]}`)
	})

	got, _, err := client.IncidentRoles.List(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != 1 || got[0].ID != "role-1" || !got[0].Required {
		t.Errorf("unexpected roles: %+v", got)
	}
}

// apiCalls lists every service method so error handling is verified uniformly.
func apiCalls(c *Client) map[string]func(ctx context.Context) error {
	return map[string]func(ctx context.Context) error{
		"Incidents.List": func(ctx context.Context) error { _, _, err := c.Incidents.List(ctx, nil); return err },
		"Incidents.Get":  func(ctx context.Context) error { _, _, err := c.Incidents.Get(ctx, "x"); return err },
		"Incidents.Create": func(ctx context.Context) error {
			_, _, err := c.Incidents.Create(ctx, &CreateIncidentOptions{})
			return err
		},
		"Incidents.Update": func(ctx context.Context) error { _, _, err := c.Incidents.Update(ctx, "x", nil); return err },
		"Incidents.ImportPostmortemDocument": func(ctx context.Context) error {
			_, _, err := c.Incidents.ImportPostmortemDocument(ctx, "x", "t", "c")
			return err
		},
		"Severities.List":    func(ctx context.Context) error { _, _, err := c.Severities.List(ctx); return err },
		"IncidentTypes.List": func(ctx context.Context) error { _, _, err := c.IncidentTypes.List(ctx); return err },
		"IncidentRoles.List": func(ctx context.Context) error { _, _, err := c.IncidentRoles.List(ctx); return err },
		"Users.List":         func(ctx context.Context) error { _, _, err := c.Users.List(ctx, nil); return err },
		"CustomFields.List":  func(ctx context.Context) error { _, _, err := c.CustomFields.List(ctx); return err },
		"CustomFields.Get":   func(ctx context.Context) error { _, _, err := c.CustomFields.Get(ctx, "x"); return err },
		"CustomFields.Create": func(ctx context.Context) error {
			_, _, err := c.CustomFields.Create(ctx, &CreateCustomFieldOptions{})
			return err
		},
		"CustomFields.Update": func(ctx context.Context) error {
			_, _, err := c.CustomFields.Update(ctx, "x", &UpdateCustomFieldOptions{})
			return err
		},
		"CustomFields.Delete": func(ctx context.Context) error { _, err := c.CustomFields.Delete(ctx, "x"); return err },
		"Actions.List":        func(ctx context.Context) error { _, _, err := c.Actions.List(ctx, nil); return err },
		"Actions.Get":         func(ctx context.Context) error { _, _, err := c.Actions.Get(ctx, "x"); return err },
		"Actions.Create": func(ctx context.Context) error {
			_, _, err := c.Actions.Create(ctx, &CreateActionOptions{})
			return err
		},
		"Actions.Update": func(ctx context.Context) error {
			_, _, err := c.Actions.Update(ctx, "x", &UpdateActionOptions{})
			return err
		},
		"Actions.Delete": func(ctx context.Context) error { _, err := c.Actions.Delete(ctx, "x"); return err },
		"Workflows.List": func(ctx context.Context) error { _, _, err := c.Workflows.List(ctx); return err },
		"Workflows.Get":  func(ctx context.Context) error { _, _, err := c.Workflows.Get(ctx, "x", false); return err },
		"Workflows.Create": func(ctx context.Context) error {
			_, _, err := c.Workflows.Create(ctx, &CreateWorkflowOptions{})
			return err
		},
		"Workflows.Update": func(ctx context.Context) error {
			_, _, err := c.Workflows.Update(ctx, "x", &UpdateWorkflowOptions{})
			return err
		},
		"Workflows.Delete": func(ctx context.Context) error { _, err := c.Workflows.Delete(ctx, "x"); return err },
		"Schedules.List":   func(ctx context.Context) error { _, _, err := c.Schedules.List(ctx, nil); return err },
		"Schedules.Get":    func(ctx context.Context) error { _, _, err := c.Schedules.Get(ctx, "x"); return err },
		"Schedules.Create": func(ctx context.Context) error {
			_, _, err := c.Schedules.Create(ctx, &CreateScheduleOptions{})
			return err
		},
		"Schedules.Update": func(ctx context.Context) error {
			_, _, err := c.Schedules.Update(ctx, "x", &UpdateScheduleOptions{})
			return err
		},
		"Schedules.Delete": func(ctx context.Context) error { _, err := c.Schedules.Delete(ctx, "x"); return err },
		"Schedules.ListEntries": func(ctx context.Context) error {
			_, _, err := c.Schedules.ListEntries(ctx, "x", nil)
			return err
		},
		"Schedules.ListOverrides": func(ctx context.Context) error {
			_, _, err := c.Schedules.ListOverrides(ctx, "x", nil)
			return err
		},
		"Schedules.GetOverride": func(ctx context.Context) error { _, _, err := c.Schedules.GetOverride(ctx, "x"); return err },
		"Schedules.CreateOverride": func(ctx context.Context) error {
			_, _, err := c.Schedules.CreateOverride(ctx, &CreateOverrideOptions{})
			return err
		},
		"Schedules.UpdateOverride": func(ctx context.Context) error {
			_, _, err := c.Schedules.UpdateOverride(ctx, "x", &UpdateOverrideOptions{})
			return err
		},
		"Schedules.DeleteOverride": func(ctx context.Context) error { _, err := c.Schedules.DeleteOverride(ctx, "x"); return err },
	}
}

func newTestClient(handler http.HandlerFunc) (*Client, func()) {
	server := httptest.NewServer(handler)
	c := NewClient("test-key", WithBaseURL(server.URL+"/v2"))
	return c, server.Close
}

func TestServices_APIErrors(t *testing.T) {
	c, closeServer := newTestClient(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = fmt.Fprint(w, `{"type": "internal_error", "status": 500, "detail": "boom"}`)
	})
	defer closeServer()

	for name, call := range apiCalls(c) {
		t.Run(name, func(t *testing.T) {
			err := call(context.Background())
			var errResp *ErrorResponse
			if !errors.As(err, &errResp) {
				t.Fatalf("error = %v (%T), want *ErrorResponse", err, err)
			}
			if errResp.Detail != "boom" || errResp.Response.StatusCode != http.StatusInternalServerError {
				t.Errorf("unexpected error response: %+v", errResp)
			}
		})
	}
}

func TestServices_InvalidJSON(t *testing.T) {
	c, closeServer := newTestClient(func(w http.ResponseWriter, _ *http.Request) {
		_, _ = fmt.Fprint(w, `{not json`)
	})
	defer closeServer()

	for name, call := range apiCalls(c) {
		// Delete methods do not decode a body.
		if strings.Contains(name, "Delete") {
			continue
		}
		t.Run(name, func(t *testing.T) {
			if err := call(context.Background()); err == nil {
				t.Error("expected a decode error")
			}
		})
	}
}

func TestServices_TransportErrors(t *testing.T) {
	server := httptest.NewServer(http.NotFoundHandler())
	c := NewClient("test-key", WithBaseURL(server.URL+"/v2"))
	server.Close() // nothing listens any more

	for name, call := range apiCalls(c) {
		t.Run(name, func(t *testing.T) {
			if err := call(context.Background()); err == nil {
				t.Error("expected a transport error")
			}
		})
	}
}

func TestListOptions_apply_Nil(t *testing.T) {
	var o *ListOptions
	q := map[string][]string{}
	o.apply(q)
	if len(q) != 0 {
		t.Errorf("query = %v, want empty", q)
	}
}

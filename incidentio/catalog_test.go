package incidentio

import (
	"context"
	"net/http"
	"net/url"
	"testing"
)

const catalogTypeJSON = `{
	"id": "ct1",
	"name": "Service",
	"description": "Services",
	"type_name": "Custom[\"Service\"]",
	"ranked": true,
	"is_editable": true,
	"use_name_as_identifier": false,
	"icon": "server",
	"color": "blue",
	"categories": ["service"],
	"annotations": {"k": "v"},
	"schema": {"version": 3, "attributes": [
		{"id": "a1", "name": "Tier", "type": "String", "array": false, "mode": "api"},
		{"id": "a2", "name": "Owner", "type": "User", "array": false, "mode": "path",
		 "path": [{"attribute_id": "p1", "attribute_name": "Team"}]}
	]},
	"created_at": "2024-01-01T00:00:00Z",
	"updated_at": "2024-01-02T00:00:00Z"
}`

func TestCatalogTypesService_ListGet(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/v3/catalog_types", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		_, _ = w.Write([]byte(`{"catalog_types":[` + catalogTypeJSON + `]}`))
	})
	mux.HandleFunc("/v3/catalog_types/ct1", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		_, _ = w.Write([]byte(`{"catalog_type":` + catalogTypeJSON + `}`))
	})

	list, _, err := client.CatalogTypes.List(context.Background())
	if err != nil || len(list) != 1 {
		t.Fatalf("got %+v, err %v", list, err)
	}
	ct := list[0]
	if ct.Name != "Service" || !ct.Ranked || ct.Annotations["k"] != "v" || ct.Schema.Version != 3 {
		t.Errorf("unexpected type: %+v", ct)
	}
	if got := ct.Schema.Attributes[1].Path[0].AttributeName; got != "Team" {
		t.Errorf("path not decoded: %q", got)
	}

	got, _, err := client.CatalogTypes.Get(context.Background(), "ct1")
	if err != nil || got.ID != "ct1" {
		t.Fatalf("got %+v, err %v", got, err)
	}
}

func TestCatalogTypesService_Mutations(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/v3/catalog_types", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
		testBody(t, r, map[string]any{"name": "Service", "description": "Services", "ranked": true, "icon": "server"})
		_, _ = w.Write([]byte(`{"catalog_type":` + catalogTypeJSON + `}`))
	})
	mux.HandleFunc("/v3/catalog_types/ct1", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPut:
			testBody(t, r, map[string]any{"name": "S2", "description": "d", "categories": []any{"service"}})
			_, _ = w.Write([]byte(`{"catalog_type":` + catalogTypeJSON + `}`))
		case http.MethodDelete:
			w.WriteHeader(http.StatusNoContent)
		default:
			t.Errorf("unexpected method %s", r.Method)
		}
	})
	mux.HandleFunc("/v3/catalog_types/ct1/actions/update_schema", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
		testBody(t, r, map[string]any{
			"version":    float64(3),
			"attributes": []any{map[string]any{"name": "Tier", "type": "String", "array": false}},
		})
		_, _ = w.Write([]byte(`{"catalog_type":` + catalogTypeJSON + `}`))
	})

	ctx := context.Background()
	ranked := true
	if ct, _, err := client.CatalogTypes.Create(ctx, &CreateCatalogTypeOptions{
		Name: "Service", Description: "Services", Ranked: &ranked, Icon: "server",
	}); err != nil || ct.ID != "ct1" {
		t.Fatalf("create: %+v, %v", ct, err)
	}
	if ct, _, err := client.CatalogTypes.Update(ctx, "ct1", &UpdateCatalogTypeOptions{
		Name: "S2", Description: "d", Categories: []string{"service"},
	}); err != nil || ct.ID != "ct1" {
		t.Fatalf("update: %+v, %v", ct, err)
	}
	if ct, _, err := client.CatalogTypes.UpdateSchema(ctx, "ct1", &UpdateCatalogTypeSchemaOptions{
		Version:    3,
		Attributes: []*CatalogTypeAttribute{{Name: "Tier", Type: "String"}},
	}); err != nil || ct.ID != "ct1" {
		t.Fatalf("update schema: %+v, %v", ct, err)
	}
	if resp, err := client.CatalogTypes.Delete(ctx, "ct1"); err != nil || resp.StatusCode != http.StatusNoContent {
		t.Fatalf("delete: %v, %v", resp, err)
	}
}

const catalogEntryJSON = `{
	"id": "ce1",
	"catalog_type_id": "ct1",
	"name": "api",
	"external_id": "ext",
	"aliases": ["a"],
	"rank": 2,
	"attribute_values": {
		"a1": {"value": {"literal": "gold", "label": "gold"}},
		"a2": {"array_value": [{"literal": "x", "label": "X"}, {"literal": "y", "label": "Y"}]}
	},
	"created_at": "2024-01-01T00:00:00Z",
	"updated_at": "2024-01-02T00:00:00Z"
}`

func TestCatalogEntriesService_ListGet(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/v3/catalog_entries", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testQuery(t, r, url.Values{"catalog_type_id": {"ct1"}, "identifier": {"api"}, "page_size": {"25"}, "after": {"c0"}})
		_, _ = w.Write([]byte(`{"catalog_entries":[` + catalogEntryJSON + `]}`))
	})
	mux.HandleFunc("/v3/catalog_entries/ce1", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		_, _ = w.Write([]byte(`{"catalog_entry":` + catalogEntryJSON + `}`))
	})

	list, _, err := client.CatalogEntries.List(context.Background(), &CatalogEntryListOptions{
		ListOptions:   ListOptions{PageSize: 25, After: "c0"},
		CatalogTypeID: "ct1",
		Identifier:    "api",
	})
	if err != nil || len(list) != 1 {
		t.Fatalf("got %+v, err %v", list, err)
	}
	e := list[0]
	if e.Rank != 2 || e.ExternalID != "ext" || e.AttributeValues["a1"].Value.Literal != "gold" || len(e.AttributeValues["a2"].ArrayValue) != 2 {
		t.Errorf("unexpected entry: %+v", e)
	}

	got, _, err := client.CatalogEntries.Get(context.Background(), "ce1")
	if err != nil || got.ID != "ce1" {
		t.Fatalf("got %+v, err %v", got, err)
	}
}

func TestCatalogEntriesService_Mutations(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/v3/catalog_entries", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
		testBody(t, r, map[string]any{
			"catalog_type_id": "ct1",
			"name":            "api",
			"rank":            float64(0),
			"attribute_values": map[string]any{
				"a1": map[string]any{"value": map[string]any{"literal": "gold"}},
				"a2": map[string]any{"array_value": []any{map[string]any{"literal": "x"}}},
			},
		})
		_, _ = w.Write([]byte(`{"catalog_entry":` + catalogEntryJSON + `}`))
	})
	mux.HandleFunc("/v3/catalog_entries/ce1", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPut:
			testBody(t, r, map[string]any{
				"name":              "api2",
				"attribute_values":  map[string]any{},
				"update_attributes": []any{"a1"},
			})
			_, _ = w.Write([]byte(`{"catalog_entry":` + catalogEntryJSON + `}`))
		case http.MethodDelete:
			w.WriteHeader(http.StatusNoContent)
		default:
			t.Errorf("unexpected method %s", r.Method)
		}
	})
	mux.HandleFunc("/v3/catalog_entries/actions/bulk_update", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
		testBody(t, r, map[string]any{
			"catalog_type_id": "ct1",
			"entries": []any{map[string]any{
				"entry_id":         "ce1",
				"attribute_values": map[string]any{"a1": map[string]any{"value": map[string]any{"literal": "silver"}}},
			}},
			"update_attributes": []any{"a1"},
		})
		w.WriteHeader(http.StatusNoContent)
	})

	ctx := context.Background()
	zero := 0
	if e, _, err := client.CatalogEntries.Create(ctx, &CreateCatalogEntryOptions{
		CatalogTypeID: "ct1",
		Name:          "api",
		Rank:          &zero,
		AttributeValues: map[string]*CatalogAttributeBinding{
			"a1": {Value: &CatalogBindingValue{Literal: "gold"}},
			"a2": {ArrayValue: []*CatalogBindingValue{{Literal: "x"}}},
		},
	}); err != nil || e.ID != "ce1" {
		t.Fatalf("create: %+v, %v", e, err)
	}
	if e, _, err := client.CatalogEntries.Update(ctx, "ce1", &UpdateCatalogEntryOptions{
		Name:             "api2",
		AttributeValues:  map[string]*CatalogAttributeBinding{},
		UpdateAttributes: []string{"a1"},
	}); err != nil || e.ID != "ce1" {
		t.Fatalf("update: %+v, %v", e, err)
	}
	if resp, err := client.CatalogEntries.Delete(ctx, "ce1"); err != nil || resp.StatusCode != http.StatusNoContent {
		t.Fatalf("delete: %v, %v", resp, err)
	}
	resp, err := client.CatalogEntries.BulkUpdate(ctx, &BulkUpdateCatalogEntriesOptions{
		CatalogTypeID: "ct1",
		Entries: []*PartialCatalogEntry{{
			EntryID:         "ce1",
			AttributeValues: map[string]*CatalogAttributeBinding{"a1": {Value: &CatalogBindingValue{Literal: "silver"}}},
		}},
		UpdateAttributes: []string{"a1"},
	})
	if err != nil || resp.StatusCode != http.StatusNoContent {
		t.Fatalf("bulk update: %v, %v", resp, err)
	}
}

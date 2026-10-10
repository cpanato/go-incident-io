package incidentio

import (
	"context"
	"fmt"
	"net/http"
	"testing"
)

func TestGetKey_MissingKeyAndNull(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/v2/incidents/missing", func(w http.ResponseWriter, _ *http.Request) {
		_, _ = fmt.Fprint(w, `{"other": {}}`)
	})
	mux.HandleFunc("/v2/incidents/null", func(w http.ResponseWriter, _ *http.Request) {
		_, _ = fmt.Fprint(w, `{"incident": null}`)
	})

	for _, id := range []string{"missing", "null"} {
		got, _, err := client.Incidents.Get(context.Background(), id)
		if err != nil {
			t.Errorf("%s: unexpected error: %v", id, err)
		}
		if got != nil {
			t.Errorf("%s: got %+v, want nil", id, got)
		}
	}
}

func TestGetKey_WrongShape(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/v2/incidents", func(w http.ResponseWriter, _ *http.Request) {
		_, _ = fmt.Fprint(w, `{"incidents": "not-a-list"}`)
	})

	if _, _, err := client.Incidents.List(context.Background(), nil); err == nil {
		t.Error("expected an error decoding a mismatched value")
	}
}

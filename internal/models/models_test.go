package models

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestParsePrice(t *testing.T) {
	cases := map[string]float64{
		"0":         0,
		"0.000002":  0.000002,
		"-1":        0, // routers use -1
		"":          0,
		"not-a-num": 0,
	}
	for in, want := range cases {
		if got := parsePrice(in); got != want {
			t.Errorf("parsePrice(%q) = %g, want %g", in, got, want)
		}
	}
}

func TestFetch(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if got := r.Header.Get("Authorization"); got != "Bearer test-key" {
			t.Errorf("auth header = %q", got)
		}
		fmt.Fprint(w, `{"data":[
			{"id":"z/model","name":"Z Model","context_length":1000,"pricing":{"prompt":"0.000001","completion":"0.000002"}},
			{"id":"a/free","name":"A Free","context_length":200000,"pricing":{"prompt":"0","completion":"0"}}
		]}`)
	}))
	defer srv.Close()

	// Point the package at the test server.
	old := catalogURL
	defer func() { catalogURL = old }()
	catalogURL = srv.URL

	catalog, err := Fetch(context.Background(), "test-key")
	if err != nil {
		t.Fatalf("Fetch: %v", err)
	}
	if len(catalog.Models) != 2 {
		t.Fatalf("models = %d, want 2", len(catalog.Models))
	}
	// Results are sorted by ID.
	if catalog.Models[0].ID != "a/free" {
		t.Errorf("first model = %q, want a/free", catalog.Models[0].ID)
	}
	if !catalog.Models[0].IsFree() {
		t.Error("a/free should be free")
	}
	if catalog.Models[1].IsFree() {
		t.Error("z/model should not be free")
	}
	if catalog.Models[1].ContextLength != 1000 {
		t.Errorf("context length = %d, want 1000", catalog.Models[1].ContextLength)
	}
}

func TestFetchBadStatus(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer srv.Close()

	old := catalogURL
	defer func() { catalogURL = old }()
	catalogURL = srv.URL

	if _, err := Fetch(context.Background(), "k"); err == nil {
		t.Fatal("expected error for non-200 status")
	}
}

func TestPreferFree(t *testing.T) {
	catalog := &Catalog{Models: []Model{
		{ID: "a/paid", PromptPrice: 0.000001},
		{ID: "b/free", PromptPrice: 0, CompletionPrice: 0},
		{ID: "c/free", PromptPrice: 0, CompletionPrice: 0},
	}}
	if got := catalog.PreferFree(); got != "b/free" {
		t.Errorf("PreferFree = %q, want b/free", got)
	}

	paidOnly := &Catalog{Models: []Model{{ID: "a/paid", PromptPrice: 0.000001}}}
	if got := paidOnly.PreferFree(); got != "" {
		t.Errorf("PreferFree = %q, want empty", got)
	}
}

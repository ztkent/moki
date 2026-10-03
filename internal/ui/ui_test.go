package ui

import (
	"strings"
	"testing"

	"github.com/ztkent/moki/internal/models"
)

func TestHumanTokens(t *testing.T) {
	cases := map[int]string{
		500:     "500",
		2000:    "2K",
		200000:  "200K",
		1000000: "1M",
		1048576: "1M",
	}
	for in, want := range cases {
		if got := humanTokens(in); got != want {
			t.Errorf("humanTokens(%d) = %q, want %q", in, got, want)
		}
	}
}

func TestModelItemDescription(t *testing.T) {
	free := modelItem{model: models.Model{ID: "a/free", Name: "A Free", ContextLength: 200000}}
	desc := free.Description()
	for _, want := range []string{"a/free", "free", "200K ctx"} {
		if !strings.Contains(desc, want) {
			t.Errorf("description %q missing %q", desc, want)
		}
	}

	paid := modelItem{model: models.Model{ID: "z/paid", Name: "Z Paid", ContextLength: 1000, PromptPrice: 0.000001}}
	if strings.Contains(paid.Description(), "free") {
		t.Errorf("paid description should not say free: %q", paid.Description())
	}
}

func TestNewPicker(t *testing.T) {
	catalog := &models.Catalog{Models: []models.Model{
		{ID: "a/one", Name: "One"},
		{ID: "b/two", Name: "Two"},
	}}
	p := NewPicker(catalog)
	if p.Chosen() != nil {
		t.Error("a fresh picker should have no selection")
	}
	if got := p.list.Items(); len(got) != 2 {
		t.Errorf("picker items = %d, want 2", len(got))
	}
}

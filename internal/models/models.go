// Package models fetches and caches the OpenRouter model catalog.
package models

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"time"
)

const cacheTTL = 24 * time.Hour

// catalogURL is a variable so tests can point it at a local server.
var catalogURL = "https://openrouter.ai/api/v1/models"

// Model is a single entry from the OpenRouter catalog.
type Model struct {
	ID              string
	Name            string
	Description     string
	ContextLength   int
	PromptPrice     float64 // USD per input token
	CompletionPrice float64 // USD per output token
}

// IsFree reports whether the model costs nothing to use.
func (m Model) IsFree() bool { return m.PromptPrice == 0 && m.CompletionPrice == 0 }

// PreferFree returns the first free model in the catalog, or "" when none
// exist. The catalog is sorted by ID, so the choice is deterministic.
func (c *Catalog) PreferFree() string {
	for _, m := range c.Models {
		if m.IsFree() {
			return m.ID
		}
	}
	return ""
}

// Catalog is a snapshot of the available models.
type Catalog struct {
	FetchedAt time.Time `json:"fetched_at"`
	Models    []Model   `json:"models"`
}

// Load returns the model catalog, preferring a fresh on-disk cache. When
// refresh is true, or the cache is stale, it fetches from OpenRouter and falls
// back to a stale cache if the network is unavailable.
func Load(ctx context.Context, apiKey string, refresh bool) (*Catalog, error) {
	path := cachePath()

	if !refresh && path != "" {
		if cached, err := readCache(path); err == nil && time.Since(cached.FetchedAt) < cacheTTL {
			return cached, nil
		}
	}

	catalog, err := Fetch(ctx, apiKey)
	if err != nil {
		if path != "" {
			if cached, cerr := readCache(path); cerr == nil {
				return cached, nil
			}
		}
		return nil, err
	}
	if path != "" {
		_ = writeCache(path, catalog)
	}
	return catalog, nil
}

// Fetch retrieves the catalog directly from OpenRouter.
func Fetch(ctx context.Context, apiKey string) (*Catalog, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, catalogURL, nil)
	if err != nil {
		return nil, fmt.Errorf("build models request: %w", err)
	}
	if apiKey != "" {
		req.Header.Set("Authorization", "Bearer "+apiKey)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("fetch models: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("fetch models: unexpected status %s", resp.Status)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read models response: %w", err)
	}

	var wire struct {
		Data []struct {
			ID            string `json:"id"`
			Name          string `json:"name"`
			Description   string `json:"description"`
			ContextLength int    `json:"context_length"`
			Pricing       struct {
				Prompt     string `json:"prompt"`
				Completion string `json:"completion"`
			} `json:"pricing"`
		} `json:"data"`
	}
	if err := json.Unmarshal(body, &wire); err != nil {
		return nil, fmt.Errorf("decode models response: %w", err)
	}

	list := make([]Model, 0, len(wire.Data))
	for _, w := range wire.Data {
		list = append(list, Model{
			ID:              w.ID,
			Name:            w.Name,
			Description:     w.Description,
			ContextLength:   w.ContextLength,
			PromptPrice:     parsePrice(w.Pricing.Prompt),
			CompletionPrice: parsePrice(w.Pricing.Completion),
		})
	}
	sort.Slice(list, func(i, j int) bool { return list[i].ID < list[j].ID })

	return &Catalog{FetchedAt: time.Now(), Models: list}, nil
}

// parsePrice converts OpenRouter's per-token price string to a float. Negative
// or unparseable values (used for routers) become zero.
func parsePrice(s string) float64 {
	v, err := strconv.ParseFloat(s, 64)
	if err != nil || v < 0 {
		return 0
	}
	return v
}

func cachePath() string {
	dir, err := os.UserCacheDir()
	if err != nil {
		return ""
	}
	return filepath.Join(dir, "moki", "models.json")
}

func readCache(path string) (*Catalog, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var catalog Catalog
	if err := json.Unmarshal(data, &catalog); err != nil {
		return nil, err
	}
	return &catalog, nil
}

func writeCache(path string, catalog *Catalog) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	data, err := json.MarshalIndent(catalog, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0o644)
}

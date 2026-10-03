// Package prefs persists user preferences between runs.
package prefs

import (
	"encoding/json"
	"os"
	"path/filepath"
)

// Prefs is the persisted user configuration.
type Prefs struct {
	Model string `json:"model"`
}

func path() string {
	dir, err := os.UserConfigDir()
	if err != nil {
		return ""
	}
	return filepath.Join(dir, "moki", "prefs.json")
}

// Load reads the saved preferences, returning empty defaults when none exist.
func Load() *Prefs {
	p := &Prefs{}
	f := path()
	if f == "" {
		return p
	}
	data, err := os.ReadFile(f)
	if err != nil {
		return p
	}
	_ = json.Unmarshal(data, p)
	return p
}

// Save writes the preferences to disk.
func (p *Prefs) Save() error {
	f := path()
	if f == "" {
		return nil
	}
	if err := os.MkdirAll(filepath.Dir(f), 0o755); err != nil {
		return err
	}
	data, err := json.MarshalIndent(p, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(f, data, 0o644)
}

// Model returns the saved model, or "" when none is set.
func Model() string { return Load().Model }

// SetModel saves the model preference.
func SetModel(id string) error {
	p := Load()
	p.Model = id
	return p.Save()
}

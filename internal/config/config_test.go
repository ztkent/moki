package config

import "testing"

func TestParseDefaults(t *testing.T) {
	cfg, err := Parse([]string{"hello", "world"})
	if err != nil {
		t.Fatalf("Parse: %v", err)
	}
	if cfg.Question != "hello world" {
		t.Errorf("question = %q, want %q", cfg.Question, "hello world")
	}
	if cfg.Temperature != DefaultTemperature {
		t.Errorf("temperature = %g, want %g", cfg.Temperature, DefaultTemperature)
	}
	if cfg.MaxTokens != DefaultMaxTokens {
		t.Errorf("max-tokens = %d, want %d", cfg.MaxTokens, DefaultMaxTokens)
	}
	if cfg.Conversation {
		t.Error("conversation should default to false")
	}
}

func TestParseFlags(t *testing.T) {
	cfg, err := Parse([]string{"-c", "-m", "openai/gpt-6-sol", "-t", "0.2", "--max-tokens", "512", "explain", "this"})
	if err != nil {
		t.Fatalf("Parse: %v", err)
	}
	if !cfg.Conversation {
		t.Error("expected conversation mode")
	}
	if cfg.Model != "openai/gpt-6-sol" {
		t.Errorf("model = %q", cfg.Model)
	}
	if cfg.Temperature != 0.2 {
		t.Errorf("temperature = %g, want 0.2", cfg.Temperature)
	}
	if cfg.MaxTokens != 512 {
		t.Errorf("max-tokens = %d, want 512", cfg.MaxTokens)
	}
	if cfg.Question != "explain this" {
		t.Errorf("question = %q", cfg.Question)
	}
}

func TestParseUnknownFlag(t *testing.T) {
	if _, err := Parse([]string{"--nope"}); err == nil {
		t.Fatal("expected error for unknown flag")
	}
}

func TestValidate(t *testing.T) {
	cases := []struct {
		name    string
		cfg     Config
		wantErr bool
	}{
		{"ok", Config{Temperature: 0.7, MaxTokens: 100}, false},
		{"temp too high", Config{Temperature: 3, MaxTokens: 100}, true},
		{"temp negative", Config{Temperature: -1, MaxTokens: 100}, true},
		{"zero tokens", Config{Temperature: 0.7, MaxTokens: 0}, true},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.cfg.Validate()
			if (err != nil) != tc.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tc.wantErr)
			}
		})
	}
}

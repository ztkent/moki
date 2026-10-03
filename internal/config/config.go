// Package config parses and validates moki's command-line options.
package config

import (
	"fmt"
	"io"
	"strings"

	"github.com/spf13/pflag"
)

const (
	// DefaultModel is used when no model is chosen via flag, picker, or env.
	// OpenRouter's auto router picks a capable model per request.
	DefaultModel = "openrouter/auto"

	// DefaultTemperature balances focus and creativity for developer answers.
	DefaultTemperature = 0.7

	// DefaultMaxTokens caps a single response. Reasoning models spend tokens on
	// hidden thinking, so this is generous by default.
	DefaultMaxTokens = 16384
)

// Config holds the resolved options for a single moki invocation.
type Config struct {
	Conversation  bool
	Model         string
	Temperature   float64
	MaxTokens     int
	ListModels    bool
	RefreshModels bool
	NoPicker      bool
	SetModel      bool
	Help          bool
	Version       bool

	// Question is the positional prompt, joined with spaces.
	Question string
}

// Parse reads args (without the program name) into a Config.
func Parse(args []string) (*Config, error) {
	cfg := &Config{}
	fs := pflag.NewFlagSet("moki", pflag.ContinueOnError)
	fs.SetOutput(io.Discard)
	fs.SortFlags = false

	fs.BoolVarP(&cfg.Help, "help", "h", false, "Show help")
	fs.BoolVarP(&cfg.Conversation, "conversation", "c", false, "Start an interactive conversation")
	fs.StringVarP(&cfg.Model, "model", "m", "", "Model to use (opens the picker when empty)")
	fs.Float64VarP(&cfg.Temperature, "temperature", "t", DefaultTemperature, "Sampling temperature (0.0-2.0)")
	fs.IntVar(&cfg.MaxTokens, "max-tokens", DefaultMaxTokens, "Maximum tokens to generate per response")
	fs.BoolVar(&cfg.ListModels, "list-models", false, "List available models and exit")
	fs.BoolVar(&cfg.RefreshModels, "refresh-models", false, "Refresh the cached model catalog")
	fs.BoolVar(&cfg.NoPicker, "no-picker", false, "Skip the interactive model picker")
	fs.BoolVar(&cfg.SetModel, "set-model", false, "Choose a model with the picker, save it, and exit")
	fs.BoolVarP(&cfg.Version, "version", "v", false, "Print the version and exit")

	if err := fs.Parse(args); err != nil {
		return nil, err
	}
	cfg.Question = strings.Join(fs.Args(), " ")
	return cfg, nil
}

// Validate checks that the option values are usable.
func (c *Config) Validate() error {
	if c.Temperature < 0 || c.Temperature > 2 {
		return fmt.Errorf("temperature must be between 0 and 2, got %g", c.Temperature)
	}
	if c.MaxTokens < 1 {
		return fmt.Errorf("max-tokens must be positive, got %d", c.MaxTokens)
	}
	return nil
}

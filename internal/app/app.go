// Package app wires Moki's configuration, model catalog, and UI together.
package app

import (
	"context"
	"fmt"
	"os"
	"strings"

	aiutil "github.com/ztkent/ai-util"
	"github.com/ztkent/moki/internal/chat"
	"github.com/ztkent/moki/internal/config"
	"github.com/ztkent/moki/internal/models"
	"github.com/ztkent/moki/internal/prefs"
	"github.com/ztkent/moki/internal/prompts"
	"github.com/ztkent/moki/internal/tools"
	"github.com/ztkent/moki/internal/ui"
)

// Version is the released version of Moki. It is overridden at build time via
// -ldflags "-X github.com/ztkent/moki/internal/app.Version=<version>".
var Version = "1.8.3"

// Run executes a single Moki invocation.
func Run(ctx context.Context, cfg *config.Config) error {
	if cfg.Help {
		fmt.Println(tools.HelpMessage)
		return nil
	}
	if cfg.Version {
		fmt.Println("moki " + Version)
		return nil
	}

	apiKey := os.Getenv("OPENROUTER_API_KEY")
	if apiKey == "" {
		return fmt.Errorf("OPENROUTER_API_KEY is not set; create one at https://openrouter.ai/keys")
	}
	client := aiutil.New(apiKey)

	if cfg.ListModels {
		return listModels(ctx, apiKey, cfg.RefreshModels)
	}
	if cfg.SetModel {
		return setModel(ctx, apiKey, cfg.RefreshModels)
	}

	model, err := resolveModel(ctx, cfg, apiKey)
	if err != nil {
		return err
	}

	if cfg.Conversation {
		return runConversation(ctx, client, cfg, model)
	}
	return runOneShot(ctx, client, cfg, model)
}

// resolveModel picks a model from the flag, environment, saved preference,
// picker, or default, in that order.
func resolveModel(ctx context.Context, cfg *config.Config, apiKey string) (string, error) {
	if cfg.Model != "" {
		return cfg.Model, nil
	}
	if env := os.Getenv("MOKI_MODEL"); env != "" {
		return env, nil
	}
	if saved := prefs.Model(); saved != "" {
		return saved, nil
	}
	if cfg.NoPicker {
		return defaultModel(ctx, apiKey, cfg.RefreshModels), nil
	}

	catalog, err := models.Load(ctx, apiKey, cfg.RefreshModels)
	if err != nil {
		// A missing catalog shouldn't block a request; fall back to the default.
		return config.DefaultModel, nil
	}
	chosen, err := ui.RunPicker(catalog, "", os.Stderr)
	if err != nil {
		return "", err
	}
	if chosen != nil {
		_ = prefs.SetModel(chosen.ID)
		return chosen.ID, nil
	}
	return defaultModel(ctx, apiKey, cfg.RefreshModels), nil
}

// setModel opens the picker, saves the choice, and exits without a request.
func setModel(ctx context.Context, apiKey string, refresh bool) error {
	catalog, err := models.Load(ctx, apiKey, refresh)
	if err != nil {
		return err
	}
	chosen, err := ui.RunPicker(catalog, prefs.Model(), os.Stderr)
	if err != nil {
		return err
	}
	if chosen == nil {
		fmt.Println("No model selected.")
		return nil
	}
	if err := prefs.SetModel(chosen.ID); err != nil {
		return err
	}
	fmt.Println("Saved model:", chosen.ID)
	return nil
}

// defaultModel prefers a free model from the catalog, falling back to the
// OpenRouter auto router when none is available.
func defaultModel(ctx context.Context, apiKey string, refresh bool) string {
	catalog, err := models.Load(ctx, apiKey, refresh)
	if err != nil {
		return config.DefaultModel
	}
	if free := catalog.PreferFree(); free != "" {
		return free
	}
	return config.DefaultModel
}

// runOneShot answers a single question, streaming to stdout.
func runOneShot(ctx context.Context, client *aiutil.Client, cfg *config.Config, model string) error {
	question := buildQuestion(cfg.Question)
	if question == "" {
		return fmt.Errorf("no question provided; run `moki --help` for usage")
	}

	session := chat.New(client, model, prompts.RequestPrompt, cfg.Temperature, cfg.MaxTokens)
	reply, err := session.SendStream(ctx, question, func(text string) {
		fmt.Print(text)
	})
	if err != nil {
		return err
	}
	fmt.Println()
	fmt.Println(ui.Footer(reply.Model, reply.Usage.TotalTokens))
	return nil
}

// buildQuestion combines a piped stdin payload with the positional question.
func buildQuestion(question string) string {
	piped := strings.TrimSpace(tools.ReadFromStdinPipe())
	switch {
	case piped == "":
		return strings.TrimSpace(question)
	case question == "":
		return piped
	default:
		return question + "\n\n" + piped
	}
}

// runConversation starts the interactive chat TUI.
func runConversation(ctx context.Context, client *aiutil.Client, cfg *config.Config, model string) error {
	catalog, _ := models.Load(ctx, os.Getenv("OPENROUTER_API_KEY"), cfg.RefreshModels)

	fmt.Println(ui.Header())
	fmt.Println(ui.Tagline(model))
	fmt.Println()

	session := chat.New(client, model, prompts.ConversationPrompt, cfg.Temperature, cfg.MaxTokens)
	return ui.RunChat(ctx, session, catalog, prompts.IntroPrompt, model, func(id string) {
		_ = prefs.SetModel(id)
	})
}

// listModels prints the catalog to stdout.
func listModels(ctx context.Context, apiKey string, refresh bool) error {
	catalog, err := models.Load(ctx, apiKey, refresh)
	if err != nil {
		return err
	}
	for _, m := range catalog.Models {
		price := "paid"
		if m.IsFree() {
			price = "free"
		}
		fmt.Printf("%-55s %-6s %s\n", m.ID, price, m.Name)
	}
	return nil
}

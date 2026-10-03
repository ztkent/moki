package ui

import "fmt"

// banner is Moki's ASCII wordmark.
const banner = `	      _    _
  /\/\   ___ | | _(_)
 /    \ / _ \| |/ / |
/ /\/\ \ (_) |   <| |  AI Assistant for the Command Line
\/    \/\___/|_|\_\_|  [https://github.com/ztkent/moki]`

// Header returns the styled startup banner.
func Header() string {
	return headerStyle.Render(banner)
}

// Tagline returns a short line describing the active model.
func Tagline(model string) string {
	return taglineStyle.Render(fmt.Sprintf("model: %s", model))
}

// Footer renders the model and token usage shown after a one-shot answer.
func Footer(model string, tokens int) string {
	if model == "" {
		model = "unknown"
	}
	text := model
	if tokens > 0 {
		text = fmt.Sprintf("%s · %d tokens", model, tokens)
	}
	return taglineStyle.Render("— " + text)
}

package app

import "testing"

func TestBuildQuestion(t *testing.T) {
	// With no piped stdin, the positional question is returned trimmed.
	if got := buildQuestion("  hello  "); got != "hello" {
		t.Errorf("buildQuestion = %q, want %q", got, "hello")
	}
	if got := buildQuestion(""); got != "" {
		t.Errorf("buildQuestion = %q, want empty", got)
	}
}

// Package tools provides small helpers shared by the CLI.
package tools

import (
	"bufio"
	"os"
	"strings"
)

// ReadFromStdinPipe returns piped stdin, or "" when stdin is a terminal.
func ReadFromStdinPipe() string {
	info, err := os.Stdin.Stat()
	if err != nil || info.Mode()&os.ModeNamedPipe == 0 {
		return ""
	}

	var b strings.Builder
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Buffer(make([]byte, 0, 64*1024), 1024*1024)
	for scanner.Scan() {
		b.WriteString(scanner.Text())
		b.WriteByte('\n')
	}
	return b.String()
}

// HelpMessage is printed for -h/--help.
const HelpMessage = `Moki — an AI assistant for the command line.

USAGE
    moki [flags] [question]

MODES
    (default)             Answer a single question and exit
    -c, --conversation    Start an interactive conversation

FLAGS
    -m, --model <id>      Model to use (opens the picker when omitted)
    -t, --temperature     Sampling temperature, 0.0-2.0 (default 0.7)
        --max-tokens      Maximum tokens per response (default 4096)
        --list-models     List available models and exit
        --refresh-models  Refresh the cached model catalog
        --no-picker       Skip the interactive model picker
    -h, --help            Show this help
        --version         Print the version

EXAMPLES
    moki "how do I undo the last git commit?"
    moki -c
    moki -m anthropic/claude-sonnet-5.5 "review this function"
    cat main.go | moki "explain this code"

ENVIRONMENT
    OPENROUTER_API_KEY    Required. Create one at https://openrouter.ai/keys
    MOKI_MODEL            Default model when none is chosen

Models are served through OpenRouter (https://openrouter.ai).`

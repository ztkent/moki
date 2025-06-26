package tools

import (
	"bufio"
	"os"
	"strings"
)

func ReadFromStdinPipe() string {
	info, _ := os.Stdin.Stat()
	if (info.Mode() & os.ModeNamedPipe) != 0 {
		scanner := bufio.NewScanner(os.Stdin)
		var input strings.Builder
		for scanner.Scan() {
			input.WriteString(scanner.Text())
			input.WriteRune('\n')
		}
		return input.String()
	}
	return ""
}

const HelpMessage = `Moki - AI Assistant for the Command Line

USAGE:
    moki [FLAGS] [QUESTION]

FLAGS:
    -h, --help                     Show this help message
    -c, --conversation            Start an interactive conversation with Moki
    -llm <provider>              Select the LLM provider: openai, replicate, or google (default: openai)
    -m <model>                   Set the model to use (uses provider default if empty)
    -t <temperature>             Set the temperature for the LLM response (default: 0.7)
    --max-tokens <number>        Set the maximum number of tokens to generate (default: 4096)
    --flags                      Log the flags used for this request

EXAMPLES:
    moki "How do I list files in Linux?"
    moki -c                                    # Start interactive conversation
    moki -llm google "Explain quantum computing"
    moki -m gpt-4 "Write a Python function"

ENVIRONMENT VARIABLES:
    OPENAI_API_KEY              OpenAI API key (required for OpenAI provider)
    REPLICATE_API_TOKEN         Replicate API token (required for Replicate provider)
    GOOGLE_API_KEY              Google AI API key (required for Google provider)
    GOOGLE_PROJECT_ID           Google Cloud Project ID (optional for Gemini API)
    LOG_LEVEL                   Set log level: debug, info, error (default: info)

For more information, visit: https://github.com/ztkent/moki`

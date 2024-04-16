package tools

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	aiutil "github.com/Ztkent/ai-util/pkg/aiutil"
)

func ConnectAIClient(aiFlag *string, modelFlag *string, temperatureFlag *float64) (*aiutil.Client, error) {
	var client *aiutil.Client
	if *aiFlag == "openai" {
		err := aiutil.MustLoadAPIKey(true, false)
		if err != nil {
			return nil, fmt.Errorf("failed to load OpenAI API key: %s", err)
		}
		if model, ok := aiutil.IsOpenAIModel(*modelFlag); ok {
			client = aiutil.MustConnectOpenAI(model, float32(*temperatureFlag))
		} else {
			client = aiutil.MustConnectOpenAI(aiutil.GPT35Turbo, float32(*temperatureFlag))
		}
	} else if *aiFlag == "anyscale" {
		err := aiutil.MustLoadAPIKey(false, true)
		if err != nil {
			return nil, fmt.Errorf("failed to load Anyscale API key: %s", err)
		}
		if model, ok := aiutil.IsAnyscaleModel(*modelFlag); ok {
			client = aiutil.MustConnectAnyscale(model, float32(*temperatureFlag))
		} else {
			client = aiutil.MustConnectAnyscale(aiutil.Mixtral8x7BInstruct, float32(*temperatureFlag))
		}
	} else {
		return nil, fmt.Errorf("invalid AI provider: %s provided, select either anyscale or openai", *aiFlag)
	}
	return client, nil
}

func ReadFromStdin() string {
	info, _ := os.Stdin.Stat()
	if (info.Mode() & os.ModeCharDevice) == 0 {
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

var HelpMessage = `Usage:
	# Ask the assistant a question
	moki [your message]
	cat moki.go | moki [tell me about this code]
  
	# Start a conversation with the assistant
	moki -c
	moki -m=turbo -c -max-tokens=100000 -t=0.5

Flags:
	-h:                        Show this message
	-c:                        Start a conversation with Moki
	-llm [openai, anyscale]:   Set the LLM Provider
	-m [string]:               Set the model to use for the LLM response
	-max-tokens [int]:         Set the maximum number of tokens to generate per response
	-t [0.0-1.0]:              Set the temperature for the LLM response
	-d:                        Show debug logging

API Keys:
	Set your API keys as environment variables:
		- export OPENAI_API_KEY=<your key>
		- export ANYSCALE_API_KEY=<your key>

Model Options:
	- OpenAI:
		- gpt-3.5-turbo, aka: turbo35
		- gpt-4-turbo-preview, aka: turbopreview
		- gpt-4-turbo, aka: turbo
	- Anyscale:
		- mistralai/Mistral-7B-Instruct-v0.1, aka: m7b
		- mistralai/Mixtral-8x7B-Instruct-v0.1, aka: m8x7b
		- codellama/CodeLlama-70b-Instruct-hf, aka: cl70b
`

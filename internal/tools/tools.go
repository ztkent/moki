package tools

import (
	"fmt"

	aiclient "github.com/Ztkent/go-openai-extended"
)

func ConnectAIClient(aiFlag *string, modelFlag *string, temperatureFlag *float64) (*aiclient.Client, error) {
	var client *aiclient.Client
	if *aiFlag == "openai" {
		err := aiclient.MustLoadAPIKey(true, false)
		if err != nil {
			return nil, fmt.Errorf("failed to load OpenAI API key: %s", err)
		}
		if model, ok := aiclient.IsOpenAIModel(*modelFlag); ok {
			client = aiclient.MustConnectOpenAI(model, float32(*temperatureFlag))
		} else {
			client = aiclient.MustConnectOpenAI(aiclient.GPT35Turbo, float32(*temperatureFlag))
		}
	} else if *aiFlag == "anyscale" {
		err := aiclient.MustLoadAPIKey(false, true)
		if err != nil {
			return nil, fmt.Errorf("failed to load Anyscale API key: %s", err)
		}
		if model, ok := aiclient.IsAnyscaleModel(*modelFlag); ok {
			client = aiclient.MustConnectAnyscale(model, float32(*temperatureFlag))
		} else {
			client = aiclient.MustConnectAnyscale(aiclient.CodeLlama34b, float32(*temperatureFlag))
		}
	} else {
		return nil, fmt.Errorf("invalid AI provider: %s provided, select either anyscale or openai", *aiFlag)
	}
	return client, nil
}

var HelpMessage = `Usage:
	moki [your question]

Flags:
	-h:                        Show this message
	-c:                        Start a conversation with Moki
	-llm [openai, anyscale]:   Set the LLM Provider
	-m [string]:               Set the model to use for the LLM response
	-max-messages [int]:       Set the maximum conversation context length
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
		- gpt-4-turbo-preview, aka: turbo
	- Anyscale:
		- mistralai/Mistral-7B-Instruct-v0.1, aka: m7b
		- mistralai/Mixtral-8x7B-Instruct-v0.1, aka: m8x7b
		- meta-llama/Llama-2-7b-chat-hf, aka: l7b
		- meta-llama/Llama-2-13b-chat-hf, aka: l13b
		- meta-llama/Llama-2-70b-chat-hf, aka: l70b
		- codellama/CodeLlama-34b-Instruct-hf, aka: cl34b
		- codellama/CodeLlama-70b-Instruct-hf, aka: cl70b
`

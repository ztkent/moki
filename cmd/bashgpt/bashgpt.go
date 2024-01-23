package main

import (
	"flag"
	"fmt"

	"github.com/Ztkent/bash-gpt/cmd/bashgpt/prompt"
	"github.com/Ztkent/bash-gpt/internal/tools"
	"github.com/Ztkent/bash-gpt/pkg/aiclient"
)

/*
Command-line interface for a BashGPT conversation.

Usage:
  bashgpt [your question]

Flags:
  -h:                        Show this message
  -c:                        Start a conversation with BashGPT
  -ai [openai, anyscale]:    Set the AI Provider
  -m [string]:               Set the model to use for the AI response
  -max-messages [int]:       Set the maximum conversation context length
  -max-tokens [int]:         Set the maximum number of tokens to generate
  -t [0.0-1.0]:              Set the temperature for the AI response

  Model Options:
    -openai:
	  - gpt-3.5-turbo, aka: turbo
	-anyscale:
	  - mistralai/Mistral-7B-Instruct-v0.1, aka: m7b
	  - mistralai/Mixtral-8x7B-Instruct-v0.1, aka: m8x7b
	  - meta-llama/Llama-2-7b-chat-hf, aka: l7b
	  - meta-llama/Llama-2-13b-chat-hf, aka: l13b
.*/

func main() {
	// Define the flags
	helpFlag := flag.Bool("-h", false, "Show this message")
	convFlag := flag.Bool("-c", false, "Start a conversation with BashGPT")
	aiFlag := flag.String("-ai", "anyscale", "Selct the AI provider, either OpenAI or Anyscale")
	modelFlag := flag.String("-m", "m7b", "Set the model to use for the AI response")
	temperatureFlag := flag.Float64("-t", 0.2, "Set the temperature for the AI response")
	maxMessagesFlag := flag.Int("-max-messages", 0, "Set the maximum conversation context length")
	maxTokensFlag := flag.Int("-max-tokens", 0, "Set the maximum number of tokens to generate")

	// Parse the flags
	flag.Parse()

	// Show the help message
	if *helpFlag {
		fmt.Println(
			`
			Usage:
				bashgpt [your question]

			Flags:
				-h:                        Show this message
				-c:                        Start a conversation with BashGPT
				-ai [openai, anyscale]:    Set the AI Provider
				-m [string]:               Set the model to use for the AI response
				-max-messages [int]:       Set the maximum conversation context length
				-max-tokens [int]:         Set the maximum number of tokens to generate
				-t [0.0-1.0]:              Set the temperature for the AI response

			Model Options:
				- OpenAI:
					- gpt-3.5-turbo, aka: turbo
				- Anyscale:
					- mistralai/Mistral-7B-Instruct-v0.1, aka: m7b
					- mistralai/Mixtral-8x7B-Instruct-v0.1, aka: m8x7b
					- meta-llama/Llama-2-7b-chat-hf, aka: l7b
					- meta-llama/Llama-2-13b-chat-hf, aka: l13b
			`)
		return
	}

	var client *aiclient.Client
	if *aiFlag == "openai" {
		if model, ok := aiclient.IsOpenAIModel(*modelFlag); ok {
			fmt.Printf("Starting conversation with OpenAI-%s\n", model)
			err := aiclient.MustLoadAPIKey(true, false)
			if err != nil {
				fmt.Printf("Failed to load OpenAI API key: %s\n", err)
				return
			}
			client = aiclient.MustConnectOpenAI(model, float32(*temperatureFlag))
		} else {
			fmt.Println(fmt.Sprintf("Invalid OpenAI model: %s provided, please provide a valid model", *modelFlag))
			return
		}
	} else if *aiFlag == "anyscale" {
		if model, ok := aiclient.IsAnyscaleModel(*modelFlag); ok {
			fmt.Printf("Starting conversation with Anyscale-%s\n", model)
			err := aiclient.MustLoadAPIKey(false, true)
			if err != nil {
				fmt.Printf("Failed to load Anyscale API key: %s\n", err)
				return
			}
			client = aiclient.MustConnectAnyscale(model, float32(*temperatureFlag))
		} else {
			fmt.Println(fmt.Sprintf("Invalid Anyscale model: %s provided, please provide a valid model", *modelFlag))
			return
		}
	} else {
		fmt.Println(fmt.Sprintf("Invalid AI provider: %s provided, select either anyscale or openai", *aiFlag))
		return
	}

	if *convFlag {
		// Start a conversation with the BashGPT prompt
		conv := aiclient.NewConversation(prompt.BashGPTPrompt, *maxMessagesFlag, *maxTokensFlag)
		err := tools.StartConversationCLI(client, conv)
		if err != nil {
			fmt.Println(err)
		}
		return
	}

	// Send a single request to the AI, return it to the user.
	// Format the response in the style of a Bash CLI response.
}

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
  -set [openai, anyscale]:   Set the AI Provider
  -model [string]:           Set the model to use for the AI response
  -max-messages [int]:       Set the maximum conversation context length
  -max-tokens [int]:         Set the maximum number of tokens to generate
  -temperature [0.0-1.0]:    Set the temperature for the AI response

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
	// Check which AI provider to use
	openAIFlag := flag.Bool("openai", false, "Use OpenAI")
	anyscaleFlag := flag.Bool("anyscale", false, "Use Anyscale")
	flag.Parse()

	var client *aiclient.Client
	if *openAIFlag && *anyscaleFlag {
		fmt.Println("Only provide one flag, please provide either -openai or -anyscale")
		return
	} else if *openAIFlag {
		fmt.Println("Starting conversation with OpenAI-GPT35Turbo")
		err := aiclient.MustLoadAPIKey(true, false)
		if err != nil {
			fmt.Println(err)
			return
		}
		client = aiclient.MustConnectOpenAI(aiclient.GPT35Turbo, 0.1)
	} else if *anyscaleFlag {
		fmt.Println("Starting conversation with Anyscale-Mistral7BInstruct")
		err := aiclient.MustLoadAPIKey(false, true)
		if err != nil {
			fmt.Println(err)
			return
		}
		client = aiclient.MustConnectAnyscale(aiclient.Mistral7BInstruct, 0.1)
	} else {
		fmt.Println("No flag provided, please provide either -openai or -anyscale")
		return
	}

	// Start a conversation with the BashGPT prompt
	err := tools.StartConversationCLI(client, prompt.BashGPTPrompt)
	if err != nil {
		fmt.Println(err)
	}
}

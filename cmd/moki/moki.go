package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	aiclient "github.com/Ztkent/go-openai-extended"
	"github.com/Ztkent/moki/internal/prompts"
	"github.com/Ztkent/moki/internal/tools"
	"github.com/sirupsen/logrus"
)

/*
Moki - An AI assistant for the command line.
*/

var logger = logrus.New()

func init() {
	// Setup the logger, so it can be parsed by datadog
	logger.Formatter = &logrus.JSONFormatter{}
	logger.SetLevel(logrus.InfoLevel)
	logger.SetOutput(os.Stdout)
}

func main() {
	// Define the flags
	helpFlag := flag.Bool("h", false, "Show this message")
	convFlag := flag.Bool("c", false, "Start a conversation with Moki")
	aiFlag := flag.String("llm", "openai", "Selct the LLM provider, either OpenAI or Anyscale")
	modelFlag := flag.String("m", "turbo35", "Set the model to use for the LLM response")
	temperatureFlag := flag.Float64("t", 0.2, "Set the temperature for the LLM response")
	maxMessagesFlag := flag.Int("max-messages", 0, "Set the maximum conversation context length")
	maxTokensFlag := flag.Int("max-tokens", 1000, "Set the maximum number of tokens to generate per response")
	ragFlag := flag.Bool("r", true, "Enable RAG functionality")

	// Parse the flags
	flag.Parse()

	// Show the help message
	if *helpFlag {
		fmt.Println(HelpMessage)
		return
	}

	//  Connect to AI Client
	client, err := ConnectAIClient(aiFlag, modelFlag, temperatureFlag)
	if err != nil {
		logger.WithFields(logrus.Fields{
			"error": err,
		}).Error("Failed to connect to the AI client")
		return
	}

	// Seed the conversation with some initial context to improve the AI responses
	conv := aiclient.NewConversation(prompts.BashGPTPrompt, *maxMessagesFlag, *maxTokensFlag, *ragFlag)
	conv.SeedConversation(map[string]string{
		"install Python 3.9 on Ubuntu":                         "sudo apt update && sudo apt install python3.9",
		"python regex to match a URL?":                         "^https?://[^/\\s]+/\\S+$",
		"list all files in a directory":                        "ls -la",
		"ammend specific old commit with commit sha":           "git rebase -i <commit-sha>",
		"run a specific command on a specific day of the week": "echo \"0 0 * * <day-of-week> <command>\" | sudo tee -a /etc/crontab",
	})

	if *convFlag {
		// Start a conversation with the Moki assistant
		err := tools.StartConversationCLI(client, conv)
		if err != nil {
			logger.WithFields(logrus.Fields{
				"error": err,
			}).Error("Conversation Failed")
		}
		return
	}

	// Require an input
	if len(flag.Args()) == 0 {
		fmt.Println("Please provide a question to ask Moki")
		return
	}

	// Response with a single request to Moki
	err = tools.LogChatStream(client, conv, strings.Join(flag.Args(), " "))
	if err != nil {
		logger.WithFields(logrus.Fields{
			"error": err,
		}).Error("Failed to log new chat stream")
	}
}

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
	logger.WithFields(logrus.Fields{
		"Model":    *modelFlag,
		"Provider": *aiFlag,
	}).Info("Starting AI Client")
	return client, nil
}

var HelpMessage = `
Usage:
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

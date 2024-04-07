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
	logger.SetOutput(os.Stdout)
	// Set the log level
	logLevel := strings.ToLower(os.Getenv("LOG_LEVEL"))
	switch logLevel {
	case "debug":
		logger.SetLevel(logrus.DebugLevel)
	case "info":
		logger.SetLevel(logrus.InfoLevel)
	case "error":
		logger.SetLevel(logrus.ErrorLevel)
	default:
		logger.SetLevel(logrus.InfoLevel)
	}
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
		fmt.Println(tools.HelpMessage)
		return
	}

	//  Connect to AI Client
	client, err := tools.ConnectAIClient(aiFlag, modelFlag, temperatureFlag)
	if err != nil {
		logger.WithFields(logrus.Fields{
			"error": err,
		}).Error("Failed to connect to the AI client")
		return
	}
	logger.WithFields(logrus.Fields{
		"Model":    *modelFlag,
		"Provider": *aiFlag,
	}).Info("Starting AI Client")

	// Seed the conversation with some initial context to improve the AI responses
	conv := aiclient.NewConversation(prompts.MokiPrompt, *maxMessagesFlag, *maxTokensFlag, *ragFlag)
	conv.SeedConversation(map[string]string{
		"install Python 3.9 on Ubuntu":                         "sudo apt update && sudo apt install python3.9",
		"python regex to match a URL?":                         "^https?://[^/\\s]+/\\S+$",
		"list all files in a directory":                        "ls -la",
		"ammend specific old commit with commit sha":           "git rebase -i <commit-sha>",
		"run a specific command on a specific day of the week": "echo \"0 0 * * <day-of-week> <command>\" | sudo tee -a /etc/crontab",
	})

	if *convFlag {
		// Start a conversation with Moki
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

	// Respond with a single request to Moki
	err = tools.LogChatStream(client, conv, strings.Join(flag.Args(), " "))
	if err != nil {
		logger.WithFields(logrus.Fields{
			"error": err,
		}).Error("Failed to log new chat stream")
	}
}

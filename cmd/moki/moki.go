package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	aiutil "github.com/Ztkent/ai-util/pkg/aiutil"
	"github.com/Ztkent/moki/internal/conversation"
	"github.com/Ztkent/moki/internal/prompts"
	"github.com/Ztkent/moki/internal/request"
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

const (
	DefaultModel     = "turbo35"
	DefaultProvider  = "openai"
	DefaultTemp      = 0.2
	DefaultMaxTokens = 100000
)

func main() {
	// Define the flags
	helpFlag := flag.Bool("h", false, "Show this message")
	convFlag := flag.Bool("c", false, "Start a conversation with Moki")
	aiFlag := flag.String("llm", DefaultProvider, "Selct the LLM provider, either OpenAI or Anyscale")
	modelFlag := flag.String("m", DefaultModel, "Set the model to use for the LLM response")
	temperatureFlag := flag.Float64("t", DefaultTemp, "Set the temperature for the LLM response")
	maxTokensFlag := flag.Int("max-tokens", DefaultMaxTokens, "Set the maximum number of tokens to generate per response")
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
		}).Errorln("Failed to connect to the AI client")
		return
	}
	logger.WithFields(logrus.Fields{
		"Model":    *modelFlag,
		"Provider": *aiFlag,
	}).Debugln("Starting AI Client")

	if *convFlag {
		// Create a new conversation with Moki
		conv := aiutil.NewConversation(prompts.ConversationPrompt, *maxTokensFlag, *ragFlag)
		// Check if there is any input from stdin
		stdinInput := tools.ReadFromStdin()
		if stdinInput != "" {
			conv.AddReference("User Input", stdinInput)
			logger.WithFields(logrus.Fields{
				"Reference": stdinInput,
			}).Debugln("Added new reference from stdin")
		}

		// Start the conversation
		err := conversation.StartConversationCLI(client, conv)
		if err != nil {
			logger.WithFields(logrus.Fields{
				"error": err,
			}).Errorln("Conversation Failed")
		}
		return
	} else {
		// Create a new conversation with Moki
		conv := aiutil.NewConversation(prompts.RequestPrompt, *maxTokensFlag, *ragFlag)
		// Seed the conversation with some initial context to improve the AI responses
		conv.SeedConversation(map[string]string{
			"install Python 3.9 on Ubuntu":                         "sudo apt update && sudo apt install python3.9",
			"python regex to match a URL?":                         "^https?://[^/\\s]+/\\S+$",
			"list all files in a directory":                        "ls -la",
			"ammend specific old commit with commit sha":           "git rebase -i <commit-sha>",
			"run a specific command on a specific day of the week": "echo \"0 0 * * <day-of-week> <command>\" | sudo tee -a /etc/crontab",
		})

		// Require an input
		if len(flag.Args()) == 0 {
			fmt.Println("Please provide a question to ask Moki")
			return
		}

		// Check if there is any input from stdin
		stdinInput := tools.ReadFromStdin()
		if stdinInput != "" {
			conv.AddReference("User Input", stdinInput)
			logger.WithFields(logrus.Fields{
				"Reference": stdinInput,
			}).Debugln("Added new reference from stdin")
		}

		// Respond with a single request to Moki
		err = request.LogChatStream(client, conv, strings.Join(flag.Args(), " "))
		if err != nil {
			logger.WithFields(logrus.Fields{
				"error": err,
			}).Errorln("Failed to log new chat stream")
		}
	}
}

package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
	aiutil "github.com/ztkent/ai-util"
	"github.com/ztkent/ai-util/types"
	"github.com/ztkent/moki/internal/conversation"
	"github.com/ztkent/moki/internal/prompts"
	"github.com/ztkent/moki/internal/tools"
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
	aiFlag := flag.String("llm", "google", "Select the LLM provider: openai, replicate, or google")
	modelFlag := flag.String("m", "", "Set the model to use for the LLM response (uses provider default if empty)")
	temperatureFlag := flag.Float64("t", 0.7, "Set the temperature for the LLM response")
	maxTokensFlag := flag.Int("max-tokens", 4096, "Set the maximum number of tokens to generate per response")
	flagFlag := flag.Bool("flags", false, "Log the flags used for this request")

	// Parse the flags
	flag.Parse()

	// Log the flags for this request
	if *flagFlag {
		logger.WithFields(logrus.Fields{
			"helpFlag":        *helpFlag,
			"convFlag":        *convFlag,
			"aiFlag":          *aiFlag,
			"modelFlag":       *modelFlag,
			"temperatureFlag": *temperatureFlag,
			"maxTokensFlag":   *maxTokensFlag,
		}).Infoln("Flags")
	}

	// Show the help message
	if *helpFlag {
		fmt.Println(tools.HelpMessage)
		return
	}

	// Create AI client using the new builder pattern
	client, err := createAIClient(*aiFlag, *modelFlag, *temperatureFlag, *maxTokensFlag)
	if err != nil {
		logger.WithFields(logrus.Fields{
			"error": err,
		}).Errorln("Failed to connect to the AI client")
		return
	}
	defer client.Close()

	logger.WithFields(logrus.Fields{
		"provider": *aiFlag,
		"model":    *modelFlag,
	}).Debugln("Started AI Client")

	if *convFlag {
		// Create a new conversation with Moki
		conv := client.NewConversation(&aiutil.ConversationConfig{
			SystemPrompt: prompts.ConversationPrompt,
			MaxTokens:    *maxTokensFlag,
			AutoTruncate: true,
		})
		err := conversation.StartConversationCLI(client, conv)
		if err != nil {
			logger.WithFields(logrus.Fields{
				"error": err,
			}).Errorln("Conversation Failed")
		}
		return
	}

	// Send a request to Moki
	conv := client.NewConversation(&aiutil.ConversationConfig{
		SystemPrompt: prompts.RequestPrompt,
		MaxTokens:    *maxTokensFlag,
		AutoTruncate: true,
	})

	// Seed the conversation with some initial context to improve the AI responses
	seedMessages := map[string]string{
		"install Python 3.9 on Ubuntu":                         "sudo apt update && sudo apt install python3.9",
		"python regex to match a URL?":                         "^https?://[^/\\s]+/\\S+$",
		"list all files in a directory":                        "ls -la",
		"ammend specific old commit with commit sha":           "git rebase -i <commit-sha>",
		"run a specific command on a specific day of the week": "echo \"0 0 * * <day-of-week> <command>\" | sudo tee -a /etc/crontab",
	}

	for question, answer := range seedMessages {
		conv.AddUserMessage(question)
		conv.AddAssistantMessage(answer)
	}

	// Require an input
	if len(flag.Args()) == 0 {
		fmt.Println("Please provide a question to ask Moki")
		return
	}

	// Respond with a single request to Moki
	err = LogChatStream(client, conv, strings.Join(flag.Args(), " "))
	if err != nil {
		logger.WithFields(logrus.Fields{
			"error": err,
		}).Errorln("Failed to log new chat stream")
	}
}

// createAIClient creates an AI client using the new builder pattern
func createAIClient(provider, model string, temperature float64, maxTokens int) (*aiutil.Client, error) {
	builder := aiutil.NewAIClient().
		WithDefaultProvider(provider).
		WithDefaultTemperature(temperature).
		WithDefaultMaxTokens(maxTokens)

	// Set model if provided
	if model != "" {
		builder = builder.WithDefaultModel(model)
	}

	// Configure providers based on the selected provider
	switch provider {
	case "openai":
		apiKey := os.Getenv("OPENAI_API_KEY")
		if apiKey == "" {
			return nil, fmt.Errorf("OPENAI_API_KEY environment variable is required")
		}
		builder = builder.WithOpenAI(apiKey)
		if model == "" {
			builder = builder.WithDefaultModel("gpt-3.5-turbo")
		}

	case "replicate":
		apiKey := os.Getenv("REPLICATE_API_TOKEN")
		if apiKey == "" {
			return nil, fmt.Errorf("REPLICATE_API_TOKEN environment variable is required")
		}
		builder = builder.WithReplicate(apiKey)
		if model == "" {
			builder = builder.WithDefaultModel("meta/meta-llama-3-8b-instruct")
		}

	case "google":
		apiKey := os.Getenv("GOOGLE_API_KEY")
		if apiKey == "" {
			return nil, fmt.Errorf("GOOGLE_API_KEY environment variable is required")
		}
		projectID := os.Getenv("GOOGLE_PROJECT_ID") // Optional for Gemini API
		builder = builder.WithGoogle(apiKey, projectID)
		if model == "" {
			builder = builder.WithDefaultModel("gemini-2.0-flash")
		}

	default:
		return nil, fmt.Errorf("unsupported provider: %s (supported: openai, replicate, google)", provider)
	}

	return builder.Build()
}

func LogChatStream(client *aiutil.Client, conv *aiutil.Conversation, userInput string) error {
	oneMin, cancel := context.WithTimeout(context.Background(), time.Second*60)
	defer cancel()

	// Use the new streaming API with correct types
	err := conv.SendStream(oneMin, userInput, "", func(ctx context.Context, response *types.StreamResponse) error {
		if response.Delta != nil && response.Delta.TextData != "" {
			fmt.Print(response.Delta.TextData)
		}
		return nil
	})

	if err != nil {
		return err
	}

	fmt.Println() // Add newline after streaming is complete
	return nil
}

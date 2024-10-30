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
	aiFlag := flag.String("llm", string(aiutil.OpenAI), "Selct the LLM provider, either OpenAI or Replicate")
	modelFlag := flag.String("m", "", "Set the model to use for the LLM response")
	temperatureFlag := flag.Float64("t", aiutil.DefaultTemp, "Set the temperature for the LLM response")
	maxTokensFlag := flag.Int("max-tokens", aiutil.DefaultMaxTokens, "Set the maximum number of tokens to generate per response")
	resourcesFlag := flag.Bool("r", true, "Enable resources functionality")
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
			"resourcesFlag":   *resourcesFlag,
		}).Infoln("Flags")
	}

	// Show the help message
	if *helpFlag {
		fmt.Println(tools.HelpMessage)
		return
	}

	//  Connect to AI Client
	client, err := aiutil.NewAIClient(*aiFlag, *modelFlag, *temperatureFlag)
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
		conv := aiutil.NewConversation(prompts.ConversationPrompt, *maxTokensFlag, *resourcesFlag)
		// Start the conversation
		err := conversation.StartConversationCLI(client, conv)
		if err != nil {
			logger.WithFields(logrus.Fields{
				"error": err,
			}).Errorln("Conversation Failed")
		}
		return
	}

	// Send a request to Moki
	conv := aiutil.NewConversation(prompts.RequestPrompt, *maxTokensFlag, *resourcesFlag)
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

	// Respond with a single request to Moki
	err = LogChatStream(client, conv, strings.Join(flag.Args(), " "))
	if err != nil {
		logger.WithFields(logrus.Fields{
			"error": err,
		}).Errorln("Failed to log new chat stream")
	}
}

func LogChatStream(client aiutil.Client, conv *aiutil.Conversation, userInput string) error {
	oneMin, cancel := context.WithTimeout(context.Background(), time.Second*60)
	defer cancel()

	// Start the chat with a fresh conversation, and the users prompt
	responseChan, errChan := make(chan string), make(chan error)

	// Check if the user's input contains a resource command
	// If so, manage the resource and add the result to the conversation
	modifiedInput, resourcesAdded, err := aiutil.ManageResources(conv, userInput)
	if err != nil {
		return err
	}
	if len(modifiedInput) == 0 {
		fmt.Println("Please provide a message to continue the conversation.")
		return nil
	} else if len(resourcesAdded) > 0 {
		fmt.Println("Resources added to conversation: ", strings.Join(resourcesAdded, ","))
	}

	go client.SendStreamRequest(oneMin, conv, modifiedInput, responseChan, errChan)
	// Read the response from the channel as it is streamed
	for {
		select {
		case response, ok := <-responseChan:
			if !ok {
				// Request channel closed
				fmt.Println()
				return nil
			}
			fmt.Print(response)
		case err := <-errChan:
			fmt.Println()
			return err
		}
	}
}

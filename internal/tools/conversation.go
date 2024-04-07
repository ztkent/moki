package tools

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	aiclient "github.com/Ztkent/go-openai-extended"
	"github.com/Ztkent/moki/internal/prompts"
)

var exitCommands = []string{"exit", "quit", "bye", ":q", "end", "q"}
var helpCommands = []string{"help", "?"}

// Start a conversation with Moki via the CLI
func StartConversationCLI(client *aiclient.Client, conv *aiclient.Conversation) error {
	// This is the maximum conversation time
	thirtyMin, cancel0 := context.WithTimeout(context.Background(), time.Minute*30)
	defer cancel0()

	oneMin, cancel := context.WithTimeout(thirtyMin, time.Minute*1)
	defer cancel()

	// Start the chat with a fresh conversation, and get the system greeting
	introChat, err := client.SendCompletionRequest(oneMin, aiclient.NewConversation(prompts.MokiPrompt, 0, 0, false), "We're starting a conversation. Introduce yourself.")
	if err != nil {
		return err
	}
	fmt.Println("Moki: " + introChat)

	// Lets start a conversation with the user via CLI
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("Request: ")
		// Ask for the user's input
		userInput, _ := reader.ReadString('\n')
		userInput = strings.TrimSpace(userInput)

		// Check if the user wants to exit
		if strings.Contains(strings.Join(exitCommands, "|"), strings.ToLower(userInput)) {
			break
		} else if strings.Contains(strings.Join(helpCommands, "|"), strings.ToLower(userInput)) {
			fmt.Println("--------------------------------------------------")
			fmt.Println("moki: ")
			fmt.Println("    Type 'exit', 'quit', or 'bye' to end the conversation.")
			fmt.Println("    Type your message to continue the conversation.")
			continue
		}

		// Check if the user's input contains a resource command
		// If so, manage the resource and add the result to the conversation
		userInput = conv.ManageRAG(userInput)
		// Check if the user provided a message
		if len(userInput) == 0 {
			fmt.Println("Please provide a message to continue the conversation.")
			continue
		}

		// Send the user's input to the LLM 🤖, wait at most 1 minute.
		oneMin, cancel = context.WithTimeout(thirtyMin, time.Minute*1)
		defer cancel()
		responseChan, errChan := make(chan string), make(chan error)
		go client.SendStreamRequest(oneMin, conv, userInput, responseChan, errChan)
		fmt.Print("Moki: ")

		// Read the response from the channel as it is streamed
		done := false
		for !done {
			select {
			case response, ok := <-responseChan:
				if !ok {
					// Request channel closed
					done = true
					break
				}
				fmt.Print(response)
			case err := <-errChan:
				if err != nil {
					return err
				}
			}
		}
		fmt.Println()
	}
	return nil
}

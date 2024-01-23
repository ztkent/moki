package tools

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/Ztkent/bash-gpt/pkg/aiclient"
)

func StartConversationCLI(client *aiclient.Client, systemPrompt string) error {
	var exitCommands = []string{"exit", "quit", "bye"}
	var helpCommands = []string{"help", "?"}

	// This is the maximum conversation time
	thirtyMin, cancel0 := context.WithTimeout(context.Background(), time.Minute*30)
	defer cancel0()

	oneMin, cancel := context.WithTimeout(thirtyMin, time.Minute*1)
	defer cancel()

	// Start the chat with a fresh conversation, and get the system greeting
	conv := aiclient.NewConversation(systemPrompt)
	introChat, err := client.SendCompletionRequest(oneMin, conv, "")
	if err != nil {
		return err
	}
	fmt.Print(introChat)

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
			fmt.Println("bashgpt: ")
			fmt.Println("    Type 'exit', 'quit', or 'bye' to end the conversation.")
			fmt.Println("    Type your message to continue the conversation.")
			continue
		}

		// Send the user's input to the AI 🤖, wait at most 1 minute.
		oneMin, cancel = context.WithTimeout(thirtyMin, time.Minute*1)
		defer cancel()
		responseChan, errChan := make(chan string), make(chan error)
		go client.SendStreamRequest(oneMin, conv, userInput, responseChan, errChan)
		fmt.Print("Response: ")

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
	}
	return nil
}

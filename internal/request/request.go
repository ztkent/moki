package request

import (
	"context"
	"fmt"
	"time"

	aiclient "github.com/Ztkent/go-openai-extended"
)

func LogChatStream(client *aiclient.Client, conv *aiclient.Conversation, chatPrompt string) error {
	oneMin, cancel := context.WithTimeout(context.Background(), time.Second*60)
	defer cancel()

	// Start the chat with a fresh conversation, and the users prompt
	responseChan, errChan := make(chan string), make(chan error)

	// Check if the user's input contains a resource command
	// If so, manage the resource and add the result to the conversation
	userInput := conv.ManageRAG(chatPrompt)
	// Check if the user provided a message
	if len(userInput) == 0 {
		return fmt.Errorf("Please provide a message to continue the conversation.")
	}

	go client.SendStreamRequest(oneMin, conv, chatPrompt, responseChan, errChan)
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

package request

import (
	"context"
	"fmt"
	"strings"
	"time"

	aiutil "github.com/Ztkent/ai-util"
)

func LogChatStream(client aiutil.Client, conv *aiutil.Conversation, userInput string) error {
	oneMin, cancel := context.WithTimeout(context.Background(), time.Second*60)
	defer cancel()

	// Start the chat with a fresh conversation, and the users prompt
	responseChan, errChan := make(chan string), make(chan error)

	// Check if the user's input contains a resource command
	// If so, manage the resource and add the result to the conversation
	modifiedInput, resourcesAdded, err := aiutil.ManageRAG(conv, userInput)
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

package conversation

import (
	"context"
	"fmt"
	"time"

	aiutil "github.com/Ztkent/ai-util/pkg/aiutil"
	"github.com/Ztkent/moki/internal/prompts"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

const (
	MaxConversationTime = time.Minute * 30
	SingleRequestTime   = time.Minute * 1
)

// Define exit commands as a slice
var exitCommands = []string{"exit", "quit", "bye", ":q", "end", "q"}
var helpCommands = []string{"help", "?"}

// StartConversationCLI starts a conversation with Moki via the CLI
func StartConversationCLI(client *aiutil.Client, conv *aiutil.Conversation) error {
	// Set the maximum conversation time
	ctx, cancel := context.WithTimeout(context.Background(), MaxConversationTime)
	defer cancel()

	// Start the chat with a fresh conversation, and get the system greeting
	introChat, err := GetIntroduction(client, ctx)
	if err != nil {
		return err
	}
	fmt.Println("Moki: " + introChat)

	// Start a conversation with the user via CLI
	return StartChat(ctx, client, conv)
}

// StartChat handles the conversation with the user
func StartChat(ctx context.Context, client *aiutil.Client, conv *aiutil.Conversation) error {
	for {
		done, err := func() (bool, error) {
			textInput := textinput.New()
			textInput.Prompt = "You: "
			m := MokiModel{Model: textInput, quit: false}
			m.Model.Focus()
			p := tea.NewProgram(m)
			defer p.RestoreTerminal()
			defer p.Quit()
			if resModel, err := p.Run(); err != nil {
				return true, err
			} else if resModel == nil {
				return true, fmt.Errorf("failed to continue the conversation.")
			} else {
				m = resModel.(MokiModel)
				if m.quit {
					fmt.Println("Goodbye!")
					return true, nil
				}
				fmt.Println("You: " + m.Value())
			}
			// Handle user's message
			err := HandleUserMessage(client, conv, ctx, m.Value())
			if err != nil {
				return false, err
			}
			return false, nil
		}()
		if err != nil {
			fmt.Println("Request Failed: ", err)
		}
		if done {
			break
		}
	}
	return nil
}

// GetIntroduction sends the initial message to start the conversation
func GetIntroduction(client *aiutil.Client, ctx context.Context) (string, error) {
	ctxWithTimeout, cancel := context.WithTimeout(ctx, SingleRequestTime)
	defer cancel()

	introChat, err := client.SendCompletionRequest(ctxWithTimeout, aiutil.NewConversation(prompts.MokiPrompt, 0, 0, false), "We're starting a conversation. Introduce yourself.")
	if err != nil {
		return "", err
	}

	return introChat, nil
}

// handleUserMessage handles the user's message
func HandleUserMessage(client *aiutil.Client, conv *aiutil.Conversation, ctx context.Context, userInput string) error {
	// Check if the user's input contains a resource command
	// If so, manage the resource and add the result to the conversation
	userInput = conv.ManageRAG(userInput)

	// Check if the user provided a message
	if len(userInput) == 0 {
		fmt.Println("Please provide a message to continue the conversation.")
		return nil
	}

	// Send the user's input to the LLM 🤖, wait at most 1 minute.
	ctxWithTimeout, cancel := context.WithTimeout(ctx, SingleRequestTime)
	defer cancel()

	responseChan, errChan := make(chan string), make(chan error)
	go client.SendStreamRequest(ctxWithTimeout, conv, userInput, responseChan, errChan)

	fmt.Print("Moki: ")
	defer fmt.Println()
	// Read the response from the channel as it is streamed
	for {
		select {
		case response, ok := <-responseChan:
			if !ok {
				// Request channel closed
				return nil
			}
			fmt.Print(response)
		case err := <-errChan:
			if err != nil {
				return err
			}
		}
	}
}

// Check if the user's input is a help command
func isHelpCommand(input string) bool {
	for _, command := range helpCommands {
		if input == command {
			return true
		}
	}
	return false
}

// Check if the user's input is an exit command
func isExitCommand(input string) bool {
	for _, command := range exitCommands {
		if input == command {
			return true
		}
	}
	return false
}

// printHelpMessage prints the help message
func printHelpMessage() {
	fmt.Println("--------------------------------------------------")
	fmt.Println("moki: ")
	fmt.Println("    Type 'exit', 'quit', or 'bye' to end the conversation.")
	fmt.Println("    Type your message to continue the conversation.")
}

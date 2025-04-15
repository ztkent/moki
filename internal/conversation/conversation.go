package conversation

import (
	"context"
	"fmt"
	"slices"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	aiutil "github.com/ztkent/ai-util"
	"github.com/ztkent/moki/internal/prompts"
	"github.com/ztkent/moki/internal/tools"
)

const (
	MaxConversationTime = time.Minute * 30
	SingleRequestTime   = time.Minute * 1
	MokiHeader          = `	      _    _
  /\/\   ___ | | _(_)
 /    \ / _ \| |/ / |
/ /\/\ \ (_) |   <| |  AI Assistant for the Command Line
\/    \/\___/|_|\_\_|  [https://github.com/ztkent/moki]`
)

var exitCommands = []string{"exit", "quit", ":q!"}

// StartConversationCLI starts a conversation with Moki via the CLI
func StartConversationCLI(client aiutil.Client, conv *aiutil.Conversation) error {
	ctx, cancel := context.WithTimeout(context.Background(), MaxConversationTime)
	defer cancel()

	fmt.Print(MokiHeader + "\n\n")
	introChat, err := GetIntroduction(client, ctx)
	if err != nil {
		return err
	}
	fmt.Println("Moki: " + introChat)

	return StartChat(ctx, client, conv)
}

// StartChat starts a chat session with Moki
// It handles user input and manages the conversation flow.
func StartChat(ctx context.Context, client aiutil.Client, conv *aiutil.Conversation) error {
	for {
		done, err := func() (bool, error) {
			textInput := textinput.New()
			textInput.Prompt = "You: "
			m := MokiModel{Model: textInput, quit: false}
			m.Model.Focus()
			p := tea.NewProgram(m)
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
			shouldExit, err := HandleUserMessage(client, conv, ctx, m.Value())
			if shouldExit {
				return true, nil
			}
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

// GetIntroduction sends an introduction request to Moki and returns the response.
func GetIntroduction(client aiutil.Client, ctx context.Context) (string, error) {
	ctxWithTimeout, cancel := context.WithTimeout(ctx, SingleRequestTime)
	defer cancel()

	introChat, err := client.SendCompletionRequest(ctxWithTimeout, aiutil.NewConversation(prompts.ConversationPrompt, 0, false), "We're starting a conversation. Introduce yourself. Your name is Moki. Only refer to yourself as Moki.")
	if err != nil {
		return introChat, err
	}
	return introChat, err
}

// HandleUserMessage handles the user's message and returns true if the user wants to exit.
func HandleUserMessage(client aiutil.Client, conv *aiutil.Conversation, ctx context.Context, userInput string) (bool, error) {
	modifiedInput, resourcesAdded, err := tools.ManageResources(conv, userInput)
	if err != nil {
		return false, err
	}

	if slices.Contains(exitCommands, strings.ToLower(strings.TrimSpace(modifiedInput))) {
		fmt.Println("Goodbye!")
		return true, nil
	}

	if len(modifiedInput) == 0 && len(resourcesAdded) == 0 {
		fmt.Println("Please provide a message or command to continue the conversation.")
		return false, nil
	} else if len(resourcesAdded) > 0 {
		fmt.Println("Resources added to conversation: ", strings.Join(resourcesAdded, ","))
	}

	ctxWithTimeout, cancel := context.WithTimeout(ctx, SingleRequestTime)
	defer cancel()

	responseChan, errChan := make(chan string), make(chan error)
	go client.SendStreamRequest(ctxWithTimeout, conv, modifiedInput, responseChan, errChan)

	firstResponse := true
	for {
		select {
		case response, ok := <-responseChan:
			if !ok {
				return false, nil
			}
			if firstResponse {
				fmt.Print("Moki: ")
				defer fmt.Println()
				firstResponse = false
			}
			fmt.Print(response)
		case err := <-errChan:
			if err != nil {
				return false, err
			}
		}
	}
}

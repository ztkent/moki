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
	"github.com/ztkent/ai-util/types"
	"github.com/ztkent/moki/internal/prompts"
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
func StartConversationCLI(client *aiutil.Client, conv *aiutil.Conversation) error {
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
func StartChat(ctx context.Context, client *aiutil.Client, conv *aiutil.Conversation) error {
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
func GetIntroduction(client *aiutil.Client, ctx context.Context) (string, error) {
	ctxWithTimeout, cancel := context.WithTimeout(ctx, SingleRequestTime)
	defer cancel()

	// Create a temporary conversation for the introduction
	conv := client.NewConversation(&aiutil.ConversationConfig{
		SystemPrompt: prompts.ConversationPrompt,
		MaxTokens:    1000,
		AutoTruncate: true,
	})

	resp, err := conv.Send(ctxWithTimeout, "We're starting a conversation. Introduce yourself. Your name is Moki. Only refer to yourself as Moki.", "")
	if err != nil {
		return "", err
	}

	return resp.Message.GetText(), nil
}

// HandleUserMessage handles the user's message and returns true if the user wants to exit.
func HandleUserMessage(client *aiutil.Client, conv *aiutil.Conversation, ctx context.Context, userInput string) (bool, error) {
	if slices.Contains(exitCommands, strings.ToLower(strings.TrimSpace(userInput))) {
		fmt.Println("Goodbye!")
		return true, nil
	}

	if len(userInput) == 0 {
		fmt.Println("Please provide a message or command to continue the conversation.")
		return false, nil
	}

	ctxWithTimeout, cancel := context.WithTimeout(ctx, SingleRequestTime)
	defer cancel()

	fmt.Print("Moki: ")
	defer fmt.Println()

	// Use the streaming API
	err := conv.SendStream(ctxWithTimeout, userInput, "", func(ctx context.Context, response *types.StreamResponse) error {
		if response.Delta != nil && response.Delta.TextData != "" {
			fmt.Print(response.Delta.TextData)
		}
		return nil
	})

	return false, err
}

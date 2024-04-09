package conversation

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	aiutil "github.com/Ztkent/ai-util/pkg/aiutil"
	"github.com/Ztkent/moki/internal/prompts"
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
	reader := bufio.NewReader(os.Stdin)

	for {
		userInput := ""
		fmt.Print("Request: ")
		// Read the user's input character by character
		for {
			char, _, err := reader.ReadRune()
			if err != nil {
				return err
			}

			// Break the loop if the user hits enter
			if char == '\n' || char == '\r' {
				break
			}

			// Trigger resource selecting when the user types '@'
			if char == '@' {
				// Manage resource selection w/ bubbletea
				m := resourceSelectionModel{resourceTypes: []string{"url", "file"}}
				p := tea.NewProgram(m)
				if m, err := p.Run(); err != nil {
					return err
				} else {
					if !m.(resourceSelectionModel).selected {
						continue
					}
				}

				resourceType := m.resourceTypes[m.cursor]

				// Prompt the user to enter the resource
				fmt.Print("Enter the " + resourceType + ": ")
				resource, _ := reader.ReadString('\n')
				resource = strings.TrimSpace(resource)

				// Add the resource to the user's input
				userInput = userInput + " -" + resourceType + ":" + resource
			} else {
				userInput += string(char)
			}

			// Handle exit and help commands
			if isExitCommand(strings.ToLower(userInput)) {
				return nil
			} else if isHelpCommand(strings.ToLower(userInput)) {
				printHelpMessage()
				break
			}
		}

		// Handle user's message
		err := handleUserMessage(client, conv, ctx, userInput)
		if err != nil {
			return err
		}
	}
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
func handleUserMessage(client *aiutil.Client, conv *aiutil.Conversation, ctx context.Context, userInput string) error {
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

type resourceSelectionModel struct {
	resourceTypes []string
	cursor        int
	selected      bool
}

func (m resourceSelectionModel) Init() tea.Cmd {
	return nil
}
func (m resourceSelectionModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch key := msg.String(); key {
		case "ctrl+c", "q", "\x03":
			return m, tea.Quit
		case "enter", "\r":
			m.selected = true
			return m, nil
		case "esc", "\x1b":
			return m, tea.Quit
		case "down", "\x1b[B":
			if m.cursor < len(m.resourceTypes)-1 {
				m.cursor++
			}
		case "up", "\x1b[A":
			if m.cursor > 0 {
				m.cursor--
			}
		}
	}
	return m, nil
}

func (m resourceSelectionModel) View() string {
	view := ""
	for i, resourceType := range m.resourceTypes {
		cursor := " "
		if m.cursor == i {
			cursor = ">"
		}
		view += fmt.Sprintf("%s %s\n", cursor, resourceType)
	}
	return view
}

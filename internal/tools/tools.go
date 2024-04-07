package tools

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strings"
	"time"

	aiclient "github.com/Ztkent/go-openai-extended"
	"github.com/Ztkent/moki/internal/prompts"
	"github.com/rs/zerolog/log"
	"github.com/sashabaranov/go-openai"
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
	introChat, err := client.SendCompletionRequest(oneMin, aiclient.NewConversation(prompts.BashGPTPrompt, 0, 0, false), "We're starting a conversation. Introduce yourself.")
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
		userInput = ManageRAG(userInput, client, conv)
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

func ManageRAG(userInput string, client *aiclient.Client, conv *aiclient.Conversation) string {
	var resourceCommands = []string{"url", "file"}
	resourcesFound := []string{}
	for _, cmd := range resourceCommands {
		re := regexp.MustCompile(fmt.Sprintf(`\-%s:(.*)`, cmd))
		matches := re.FindAllStringSubmatch(userInput, -1)
		for _, match := range matches {
			if len(match) > 1 {
				resource := strings.TrimSpace(match[1])
				resourcesFound = append(resourcesFound, cmd+":"+resource)
				GenerateResource(client, conv, resource)
				userInput = strings.Replace(userInput, "-"+cmd+":"+resource, "", -1)
			}
		}
	}

	if len(resourcesFound) > 0 {
		fmt.Println("Resources found: " + strings.Join(resourcesFound, ", "))
		fmt.Println("User input: " + userInput)
	}
	return userInput
}

// Log the results of a fresh chat stream
func LogNewChatStream(client *aiclient.Client, conv *aiclient.Conversation, chatPrompt string) error {
	oneMin, cancel := context.WithTimeout(context.Background(), time.Second*60)
	defer cancel()

	// Start the chat with a fresh conversation, and the users prompt
	responseChan, errChan := make(chan string), make(chan error)
	log.Debug().Msg(fmt.Sprintf("prompt: " + chatPrompt))

	// Check if the user's input contains a resource command
	// If so, manage the resource and add the result to the conversation
	userInput := ManageRAG(chatPrompt, client, conv)
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

// Fetch and add a resource to the conversation
// There is a LIMIT to the number of tokens.
// Cant just send the whole page
func GenerateResource(client *aiclient.Client, conv *aiclient.Conversation, path string) {
	// Determine if the path is a valid url
	url, err := url.Parse(path)
	if err == nil && url.Scheme != "" && url.Host != "" {
		// Fetch the URL
		fmt.Println("Downloading file from URL: " + path)
		resp, err := http.Get(path)
		if err != nil {
			fmt.Println("Error fetching URL: " + err.Error())
			return
		}
		defer resp.Body.Close()
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			fmt.Println("Error reading response body: " + err.Error())
			return
		}

		// handle the html content
		html := string(body)
		// Build the System Message
		messageParts := make([]openai.ChatMessagePart, 0)
		messageParts = append(messageParts, openai.ChatMessagePart{
			Type: openai.ChatMessagePartTypeText,
			Text: "URL: " + path,
		})
		messageParts = append(messageParts, openai.ChatMessagePart{
			Type: openai.ChatMessagePartTypeText,
			Text: "Status: " + resp.Status,
		})
		messageParts = append(messageParts, openai.ChatMessagePart{
			Type: openai.ChatMessagePartTypeText,
			Text: "Content: " + html,
		})
		// Set the resource
		conv.Append(openai.ChatCompletionMessage{
			Name:         url.String(),
			Role:         openai.ChatMessageRoleSystem,
			MultiContent: messageParts,
			// Content:      html,
		})
	} else if _, err := os.Stat(path); !os.IsNotExist(err) {
		fmt.Println("Uploading file from path: " + path)
		// Open the file
		file, err := os.Open(path)
		if err != nil {
			fmt.Println("Error opening file: " + err.Error())
			return
		}
		// read the contents of the file
		fileInfo, err := file.Stat()
		if err != nil {
			fmt.Println("Error reading file: " + err.Error())
			return
		}
		fileSize := fileInfo.Size()
		fileContents := make([]byte, fileSize)
		_, err = file.Read(fileContents)
		if err != nil {
			fmt.Println("Error reading file: " + err.Error())
			return
		}
		// Close the file
		err = file.Close()
		if err != nil {
			fmt.Println("Error closing file: " + err.Error())
			return
		}
		// add the content, and path of the url to a json object
		resJson, err := json.Marshal(map[string]interface{}{
			"path":     path,
			"contents": string(fileContents),
		})
		if err != nil {
			fmt.Println("Error reading file: " + err.Error())
			return
		}
		conv.Append(openai.ChatCompletionMessage{
			Name:    path,
			Role:    openai.ChatMessageRoleSystem,
			Content: string(resJson),
		})
	} else {
		fmt.Println("Invalid path: " + path)
	}
}

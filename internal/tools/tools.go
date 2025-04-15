package tools

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"

	aiutil "github.com/ztkent/ai-util"
)

func ReadFromStdinPipe() string {
	info, _ := os.Stdin.Stat()
	if (info.Mode() & os.ModeNamedPipe) != 0 {
		scanner := bufio.NewScanner(os.Stdin)
		var input strings.Builder
		for scanner.Scan() {
			input.WriteString(scanner.Text())
			input.WriteRune('\n')
		}
		return input.String()
	}
	return ""
}

// Determine if the user's input contains a resource command
// There is usually some limit to the number of tokens
func ManageResources(conv *aiutil.Conversation, userInput string) (string, []string, error) {
	resourcesFound := []string{}
	if conv == nil {
		return userInput, resourcesFound, fmt.Errorf("Failed to ManageResources: Conversation is nil")
	} else if len(userInput) == 0 {
		return userInput, resourcesFound, nil
	}

	// Check if there is any input from stdin
	stdinInput := ReadFromStdinPipe()
	if stdinInput != "" {
		resourcesFound = append(resourcesFound, "stdin:"+stdinInput)
		conv.AddReference("User Input", stdinInput)
	}

	// Only supporting URL and File resources for now
	var resourceCommands = []string{"url", "file"}
	for _, cmd := range resourceCommands {
		re := regexp.MustCompile(fmt.Sprintf(`\-%s:(.*)`, cmd))
		matches := re.FindAllStringSubmatch(strings.ToLower(userInput), -1)
		for _, match := range matches {
			if len(match) > 1 {
				resource := strings.TrimSpace(match[1])
				resourcesFound = append(resourcesFound, cmd+":"+resource)
				err := aiutil.AddResource(conv, resource, cmd)
				if err != nil {
					return userInput, resourcesFound, err
				}
				userInput = strings.Replace(userInput, "-"+cmd+":"+resource, "", -1)
			}
		}
	}
	return userInput, resourcesFound, nil
}

var HelpMessage = `Usage:
	# Ask the assistant a question
	moki [your message]

	# Provide additional context
	cat moki.go | moki [tell me about this code]
	moki [tell me about this code]    -file:moki.go
	moki [tell me about this project] -url:https://github.com/ztkent/moki

	# Start a conversation with the assistant
	moki -c
	moki -c -m=turbo -max-tokens=100000 -t=0.5

Flags:
	-h:                        Show this message
	-c:                        Start a conversation with Moki
	-llm:                      Set the LLM Provider
	-m:                        Set the model to use for the LLM response
	-max-tokens: 	           Set the maximum number of tokens to generate per response
	-t:                        Set the temperature for the LLM response
	-d:                        Show debug logging

API Keys:
	- export OPENAI_API_KEY=<your key>
	- export REPLICATE_API_TOKEN=<your key>

Model Options:
	- OpenAI:
		- gpt-3.5-turbo, aka: turbo35
		- gpt-4-turbo, aka: turbo
		- gpt-4o
		- gpt-4o-mini
		- o1-preview
		- o1-mini
		- [Default] gpt-4.1
	- Replicate:
		- [Default] meta-llama-3-8b, aka: l3-8b
		- meta-llama-3-8b-instruct, aka: l3-8b-instruct
		- meta-llama-3-70b, aka: l3-70b
		- meta-llama-3-70b-instruct, aka: l3-70b-instruct
`

package tools

import (
	"bufio"
	"os"
	"strings"
)

func ReadFromStdin() string {
	info, _ := os.Stdin.Stat()
	if (info.Mode() & os.ModeCharDevice) == 0 {
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

var HelpMessage = `Usage:
	# Ask the assistant a question
	moki [your message]

	# Provide additional context
	cat moki.go | moki [tell me about this code]
	moki [tell me about this code]    -file:moki.go
	moki [tell me about this project] -url:https://github.com/Ztkent/moki

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
	- export ANYSCALE_API_KEY=<your key>
	- export REPLICATE_API_TOKEN=<your key>

Model Options:
	- OpenAI:
		- [Default] gpt-3.5-turbo, aka: turbo35
		- gpt-4-turbo, aka: turbo
	- Replicate:
		- [Default] meta-llama-3-8b, aka: l3-8b (default)
		- meta-llama-3-8b-instruct, aka: l3-8b-instruct
		- meta-llama-3-70b, aka: l3-70b
		- meta-llama-3-70b-instruct, aka: l3-70b-instruct
	- Anyscale:
		- [Default] mistralai/Mixtral-8x7B-Instruct-v0.1, aka: m8x7b (default)
		- mistralai/Mistral-7B-Instruct-v0.1, aka: m7b
		- codellama/CodeLlama-70b-Instruct-hf, aka: cl70b
`

# bashgpt
A GPT assistant that helps with command-line tasks like finding files, installing packages, and git.

## Usage

- Install bashgpt:  
  ```bash
  go install github.com/Ztkent/bashgpt/cmd/bashgpt
  ```
  
- Set your API key as an environment variable:
  ```bash
  export OPENAI_API_KEY=<your key>
  export ANYSCALE_API_KEY=<your key>
  ```

- Run the assistant:
  ```bash
  # Ask the assistant a question
  bashgpt [your message]

  # Start a conversation with the assistant
  bashgpt -c
  bashgpt -llm=openai -c -max-messages=250 -max-tokens=1000 -t=0.5
  ```

### Options
- There are two options for the API provider:  
  - OpenAI (https://platform.openai.com/docs/overview)  
  - Anyscale (https://www.anyscale.com/endpoints)  
```
Flags:
  -h:                        Show this message
  -c:                        Start a conversation with BashGPT
  -llm [openai, anyscale]:   Set the LLM Provider
  -m [string]:               Set the model to use for the LLM response
  -max-messages [int]:       Set the maximum conversation context length
  -max-tokens [int]:         Set the maximum number of tokens to generate per response
  -t [0.0-1.0]:              Set the temperature for the LLM response
  -d:                        Show debug logging

Model Options:
  - OpenAI:
    - gpt-3.5-turbo, aka: turbo
  - Anyscale:
    - mistralai/Mistral-7B-Instruct-v0.1, aka: m7b
    - mistralai/Mixtral-8x7B-Instruct-v0.1, aka: m8x7b
    - meta-llama/Llama-2-7b-chat-hf, aka: l7b
    - meta-llama/Llama-2-13b-chat-hf, aka: l13b
    - meta-llama/Llama-2-70b-chat-hf, aka: l70b
    - codellama/CodeLlama-34b-Instruct-hf, aka: cl34b
```

#### API Provider
By default the assistant will use the Anyscale API. To use OpenAI, run the assistant with a flag. 

```bash
bashgpt -llm=openai
bashgpt -llm=anyscale 
```

#### Model
Depending on the LLM Provider selected, different models are available.  
By default the Anyscale API uses `CodeLlama-34b`, and OpenAI uses `gpt-3.5-turbo`.
```bash
bashgpt -m=m8x7b
```

#### Conversation
The assistant can be used in conversation mode.  
This allows the assistant to remember previous messages and use them to generate more in-depth responses.
```bash
bashgpt -c
```

#### Conversation Context
Larger conversations require more tokens, by default the conversation context is limited to 100 messages.  
```bash
bashgpt -max-messages=250
```

#### Token Limit
Tokens cost money, by default the assistant will generate as many tokens as it needs for the converation.
```bash
bashgpt -max-tokens=10000
```

#### Temperature
The temperature of the LLM response is a measure of randomness. Adjust this value to taste.
Temperature is a float between 0 and 1. By default the temperature is 0.2
```bash
bashgpt -t=0.5
```
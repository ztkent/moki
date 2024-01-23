# bashgpt
An AI assitant that helps with command-line tasks like finding files, installing packages, and git.

### Usage

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
  bashgpt
  ```

### Options
- There are two options for the API provider:  
  - OpenAI (https://platform.openai.com/docs/overview)  
  - Anyscale (https://www.anyscale.com/endpoints)  
```
Flags:
  -h:                        Show this message
  -c:                        Start a conversation with BashGPT
  -set [openai, anyscale]:   Set the AI Provider
  -model [string]:           Set the model to use for the AI response
  -max-messages [int]:       Set the maximum conversation context length
  -max-tokens [int]:         Set the maximum number of tokens to generate
  -temperature [0.0-1.0]:    Set the temperature for the AI response

  Model Options:
    -openai:
	  - gpt-3.5-turbo, aka: turbo
	-anyscale:
	  - mistralai/Mistral-7B-Instruct-v0.1, aka: m7b
	  - mistralai/Mixtral-8x7B-Instruct-v0.1, aka: m8x7b
	  - meta-llama/Llama-2-7b-chat-hf, aka: l7b
	  - meta-llama/Llama-2-13b-chat-hf, aka: l13b
```

#### API Provider
By default the assistant will use the Anyscale API. To use OpenAI, run the assistant with a flag.  
Future requests to the assistant will use the selected provider. Conversation context will be maintained.

```bash
Flags
bashgpt -set openai
bashgpt -set anyscale 
```

#### Model
Depending on the AI Provider selected, different models are available.  
By default the Anyscale API uses the `Mistral-7B-Instruct-v0.1`, and OpenAI uses `gpt-3.5-turbo`.  
To adjust the model to `Mixtral-8x7B-Instruct-v0.1`, do this:
```bash
bashgpt -model m8x7b
```

#### Conversation Context
Larger conversations require more tokens, by default the conversation context is limited to 100 messages.  
To increase the limit, set this flag:
```bash
bashgpt -max-messages 250
```

#### Token Limit
Tokens cost money, by default the assistant will generate as many tokens as it needs for the converation.
To limit these tokens, set this flag:
```bash
bashgpt -max-tokens 10000
```

#### Temperature
The temperature of the AI response is a measure of randomness. Adjust this value to taste.
Temperature is a float between 0 and 1. By default the temperature is 0.2
```bash
bashgpt -temperature 0.5
```
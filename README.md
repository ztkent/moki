# BashGPT
A GPT assistant for the command line.  
Tuned to assist with developer tasks like finding files, installing packages, and git.


## Usage

- Install bashgpt:  
  ```bash
  go install github.com/Ztkent/bash-gpt/cmd/bashgpt@latest
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

## Examples
``` 
bashgpt install Python 3.9 on Ubuntu
- sudo apt update && sudo apt install python3.9

bashgpt given a text file, wrap each line in quotes. format it for display
- sed 's/.*/"&"/' file.txt

bashgpt given a text file, wrap each line in quotes. format it for display with python
- with open('file.txt', 'r') as f:
  lines = f.readlines()
  lines = ['"{}"'.format(line.strip()) for line in lines]
  print('\n'.join(lines))

bashgpt update git email and username
- git config --global user.email "youremail@example.com"
  git config --global user.name "Your Name"

bashgpt git re-edit a specific commit
- git commit --amend
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
    - gpt-3.5-turbo, aka: turbo3
    - gpt-4-turbo-preview, aka: turbo
  - Anyscale:
    - mistralai/Mistral-7B-Instruct-v0.1, aka: m7b
    - mistralai/Mixtral-8x7B-Instruct-v0.1, aka: m8x7b
    - meta-llama/Llama-2-7b-chat-hf, aka: l7b
    - meta-llama/Llama-2-13b-chat-hf, aka: l13b
    - meta-llama/Llama-2-70b-chat-hf, aka: l70b
    - codellama/CodeLlama-34b-Instruct-hf, aka: cl34b
    - codellama/CodeLlama-70b-Instruct-hf, aka: cl70b
```

#### API Provider
By default the assistant will use OpenAI. To use Anyscale, run the assistant with a flag. 

```bash
bashgpt -llm=openai
bashgpt -llm=anyscale 
```

#### Model
Depending on the LLM Provider selected, different models are available.  
By default the Anyscale API uses `Mistral-8x7b`, and OpenAI uses `gpt-4-turbo-preview`.
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
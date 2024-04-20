<a href="https://github.com/ztkent/moki/tags"><img src="https://img.shields.io/github/v/tag/ztkent/moki.svg" alt="Latest Release"></a>
<a href="https://github.com/ztkent/moki/actions"><img src="https://github.com/ztkent/moki/actions/workflows/build.yml/badge.svg?branch=main" alt="Build Status"></a>

# <img width="40" alt="logo_moki" src="https://github.com/Ztkent/moki/assets/7357311/f1dfb864-3c20-4384-898b-1acc4bb7c92f"> Moki

An AI assistant for the command line.  

Tuned to assist with developer tasks like finding files, installing packages, and git.   
Conversation mode can explain code snippets, generate unit tests, and scaffold new projects.

## Usage
- Install moki:  
  ```bash
  go install github.com/Ztkent/moki/cmd/moki@latest
  ```
  
- Set your API key as an environment variable:
  ```bash
  export OPENAI_API_KEY=<your key>
  export ANYSCALE_API_KEY=<your key>
  export REPLICATE_API_TOKEN=<your key>
  ```

- Run the assistant:
  ```bash
  # Ask the assistant a question
  moki [your message]

  # Provide additional context
  cat moki.go | moki [tell me about this code]
  moki [tell me about this code]    -file:moki.go
  moki [tell me about this project] -url:https://github.com/Ztkent/moki

  # Start a conversation with the assistant
  moki -c
  moki -c -m=turbo -max-tokens=100000 -t=0.5
  ```

## Example
https://github.com/Ztkent/moki/assets/7357311/2b839654-9d34-4307-a76c-598d9c09048e

## Configuration
- There are a few options for the API provider:  
  - OpenAI (https://platform.openai.com/docs/overview)  
  - Replicate (https://replicate.com/docs)
  - Anyscale (https://www.anyscale.com/endpoints)  
```
Flags:
  -c:                        Start a conversation with Moki
  -llm:                      Set the LLM Provider
  -m:                        Set the model to use for the LLM response
  -max-tokens:               Set the maximum number of tokens to generate
  -t:                        Set the temperature for the LLM response
  -d:                        Show debug logging

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
```

#### Conversation
The assistant can be used in conversation mode.  
This allows the assistant to generate more in-depth responses.
```bash
moki -c
```

#### API Provider
By default the assistant will use OpenAI. To use another, run the assistant with a flag. 
```bash
moki -llm=openai
moki -llm=anyscale 
moki -llm=replicate 
```

#### Model
Depending on the LLM Provider selected, different models are available.  
```bash
moki -m=turbo
moki -m=m8x7b
moki -m=l3-70b
```

#### Token Limit
Tokens cost money.   
By default the assistant will limit any conversation to 100k tokens.
```bash
moki -max-tokens=100000
```

#### Temperature
The temperature of an LLM response is a measure of randomness.   
The value float between 0 and 1. By default the temperature is 0.2
```bash
moki -t=0.5
```

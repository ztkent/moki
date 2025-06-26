<a href="https://github.com/ztkent/moki/tags"><img src="https://img.shields.io/github/v/tag/ztkent/moki.svg" alt="Latest Release"></a>
<a href="https://github.com/ztkent/moki/actions"><img src="https://github.com/ztkent/moki/actions/workflows/build.yml/badge.svg?branch=main" alt="Build Status"></a>
# <img width="40" alt="logo_moki" src="https://github.com/ztkent/moki/assets/7357311/f1dfb864-3c20-4384-898b-1acc4bb7c92f"> Moki

An AI Assistant for the command line. Tuned to assist with developer tasks.

## Installation

```bash
go install github.com/ztkent/moki/cmd/moki@latest
```

### Flags

| Flag | Description | Default |
|------|-------------|---------|
| `-h, --help` | Show help message | |
| `-c, --conversation` | Start interactive conversation | |
| `-llm <provider>` | Select LLM provider: openai, replicate, google | openai |
| `-m <model>` | Set specific model (uses provider default if empty) | |
| `-t <temperature>` | Set response temperature (0.0-1.0) | 0.7 |
| `--max-tokens <number>` | Set maximum tokens to generate | 4096 |
| `--flags` | Log the flags used for this request | |

## Configuration

### Environment Variables

Set up your API keys for the providers you want to use:

```bash
# OpenAI
export OPENAI_API_KEY="your-openai-api-key"

# Google (Gemini)
export GOOGLE_API_KEY="your-google-api-key"
export GOOGLE_PROJECT_ID="your-project-id"  # Optional

# Replicate
export REPLICATE_API_TOKEN="your-replicate-token"

# Logging (optional)
export LOG_LEVEL="info"  # debug, info, error
```

### Default Models

- **OpenAI**: `gpt-3.5-turbo`
- **Google**: `gemini-2.0-flash`
- **Replicate**: `meta/meta-llama-3-8b-instruct`

## Examples

```bash
# Get help with Git commands
moki "How do I undo the last commit in Git?"

# Python programming assistance
moki "Create a Python decorator for timing functions"

# System administration
moki "How to check disk usage on Ubuntu?"
```

### Interactive Mode

```bash
# Start a conversation
moki -c

# In conversation mode, you can:
# - Ask follow-up questions
# - Build on previous context
# - Get detailed explanations
```

### Provider-Specific Usage

```bash
# Use Google's Gemini for creative tasks
moki -llm google "Write a haiku about programming"

# Use Replicate for open-source models
moki -llm replicate "Explain machine learning basics"

# Use OpenAI with specific model
moki -llm openai -m gpt-4 "Review this code architecture"
```

### Advanced Configuration

```bash
# High creativity for creative writing
moki -t 0.9 "Write a creative story about a robot"

# Low temperature for factual queries
moki -t 0.1 "What is the capital of France?"

# Limit response length
moki --max-tokens 500 "Summarize quantum physics briefly"
```

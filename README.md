<a href="https://github.com/ztkent/moki/tags"><img src="https://img.shields.io/github/v/tag/ztkent/moki.svg" alt="Latest Release"></a>
<a href="https://github.com/ztkent/moki/actions"><img src="https://github.com/ztkent/moki/actions/workflows/build.yml/badge.svg?branch=main" alt="Build Status"></a>
# <img width="40" alt="logo_moki" src="https://github.com/ztkent/moki/assets/7357311/f1dfb864-3c20-4384-898b-1acc4bb7c92f"> Moki

An AI assistant for the command line.

Moki quickly answers a single question, or manages an interactive conversation.

## Installation

```bash
go install github.com/ztkent/moki/cmd/moki@latest
```

## Setup

Moki needs an OpenRouter API key:
```bash
export OPENROUTER_API_KEY="sk-or-..."
```

Optionally set a default model:

```bash
export MOKI_MODEL="anthropic/claude-sonnet-5.5"
```

## Usage

```bash
# Ask a question
moki "how do I undo the last git commit but keep the changes?"
moki "find files larger than 100MB in the current directory"
moki "what's the difference between git fetch and git pull?"

# Pipe in context
cat main.go | moki "explain this code"
git diff | moki "write a commit message for these changes"

# Interactive conversation
moki -c

# Choose a model
moki --set-model
moki -m openai/gpt-6-sol "review this function"
```

### Flags

| Flag | Description | Default |
|------|-------------|---------|
| `-h, --help` | Show help | |
| `-c, --conversation` | Start an interactive conversation | |
| `-m, --model <id>` | Model to use (opens the picker when omitted) | |
| `-t, --temperature <n>` | Sampling temperature, 0.0–2.0 | `0.7` |
| `--max-tokens <n>` | Maximum tokens per response | `16384` |
| `--set-model` | Choose a model with the picker, save it, and exit | |
| `--list-models` | List available models and exit | |
| `--refresh-models` | Refresh the cached model catalog | |
| `--no-picker` | Skip the interactive model picker | |
| `--version` | Print the version | |

### Conversation commands

Via `moki -c`:

| Command | Description |
|---------|-------------|
| `/model [id]` | Switch model (opens the picker when no id is given) |
| `/clear` | Clear the conversation history |
| `/help` | Show the command list |
| `/exit` | Quit (also `/quit`, `/q`, or `ctrl+c`) |

Changing the model with `/model` saves it as your default for next time.

## Environment variables

| Variable | Description |
|----------|-------------|
| `OPENROUTER_API_KEY` | **Required.** Your OpenRouter API key. |
| `MOKI_MODEL` | Default model when none is chosen. |

## Configuration files

| Path | Purpose |
|------|---------|
| `$XDG_CONFIG_HOME/moki/prefs.json` | Saved model preference. |
| `$XDG_CACHE_HOME/moki/models.json` | Cached model catalog (24h TTL). |

## Development

```bash
make build   # build ./moki
make test    # run tests with the race detector
make vet     # run go vet
make install # install to $GOPATH/bin
```

Moki is built on [`ai-util`](https://github.com/ztkent/ai-util)

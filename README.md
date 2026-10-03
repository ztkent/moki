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

### Ask a question

```bash
moki "how do I undo the last git commit but keep the changes?"
moki "find files larger than 100MB in the current directory"
moki "what's the difference between git fetch and git pull?"
```

Moki answers directly and prints the model and tokens used:

```
git reset --soft HEAD~1
— anthropic/claude-sonnet-5.5 · 96 tokens
```

### Pipe in context

Feed a file, a diff, or command output straight into the question:

```bash
cat main.go | moki "explain this code"
git diff | moki "write a commit message for these changes"
kubectl logs my-pod | moki "why is this crashing?"
```

### Start a conversation

```bash
moki -c
```

Follow-up questions keep their context. Inside a conversation:

```
/model          # switch models with the picker
/clear          # start fresh
/exit           # quit
```

### Choose a model

```bash
moki --set-model                    # pick once, saved for next time
moki -m openai/gpt-6-sol "..."      # use a model for one request
moki --list-models                  # see everything available
```

Your saved model is reused automatically, so you only pick once. The picker
marks it with a ★.

### Tune the response

```bash
moki -t 0.2 "what does this regex match?"     # more focused
moki -t 1.2 "brainstorm names for a CLI tool" # more creative
moki --max-tokens 500 "summarize this in a sentence"
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

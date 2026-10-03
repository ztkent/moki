<a href="https://github.com/ztkent/moki/tags"><img src="https://img.shields.io/github/v/tag/ztkent/moki.svg" alt="Latest Release"></a>
<a href="https://github.com/ztkent/moki/actions"><img src="https://github.com/ztkent/moki/actions/workflows/build.yml/badge.svg?branch=main" alt="Build Status"></a>
# <img width="40" alt="logo_moki" src="https://github.com/ztkent/moki/assets/7357311/f1dfb864-3c20-4384-898b-1acc4bb7c92f"> Moki

An AI assistant for the command line, tuned for developer tasks.

Moki answers a single question and exits, or drops into an interactive
conversation. Every model is served through [OpenRouter](https://openrouter.ai),
so one API key unlocks the whole catalog.

## Installation

```bash
go install github.com/ztkent/moki/cmd/moki@latest
```

## Setup

Moki needs an OpenRouter API key. Create one at
<https://openrouter.ai/keys> and export it:

```bash
export OPENROUTER_API_KEY="sk-or-..."
```

Optionally set a default model so the picker is skipped:

```bash
export MOKI_MODEL="anthropic/claude-sonnet-5.5"
```

## Usage

```bash
# One-shot: answer a question and exit
moki "how do I undo the last git commit?"

# Interactive conversation
moki -c

# Pipe context in
cat main.go | moki "explain this code"

# Choose a model explicitly
moki -m openai/gpt-6-sol "review this function"
```

On first run Moki fetches the OpenRouter model catalog and shows an
interactive picker. The catalog is cached for 24 hours.

### Flags

| Flag | Description | Default |
|------|-------------|---------|
| `-h, --help` | Show help | |
| `-c, --conversation` | Start an interactive conversation | |
| `-m, --model <id>` | Model to use (opens the picker when omitted) | |
| `-t, --temperature <n>` | Sampling temperature, 0.0–2.0 | `0.7` |
| `--max-tokens <n>` | Maximum tokens per response | `4096` |
| `--list-models` | List available models and exit | |
| `--refresh-models` | Refresh the cached model catalog | |
| `--no-picker` | Skip the interactive model picker | |
| `--version` | Print the version | |

### Conversation commands

Inside `moki -c`:

| Command | Description |
|---------|-------------|
| `/model [id]` | Switch model (opens the picker when no id is given) |
| `/clear` | Clear the conversation history |
| `/help` | Show the command list |
| `/exit` | Quit (also `/quit`, `/q`, or `ctrl+c`) |

## Environment variables

| Variable | Description |
|----------|-------------|
| `OPENROUTER_API_KEY` | **Required.** Your OpenRouter API key. |
| `MOKI_MODEL` | Default model when none is chosen. |

## Development

```bash
make build   # build ./moki
make test    # run tests with the race detector
make vet     # run go vet
make install # install to $GOPATH/bin
```

Moki is built on [`ai-util`](https://github.com/ztkent/ai-util), a small
OpenRouter client with streaming and tool-calling support.

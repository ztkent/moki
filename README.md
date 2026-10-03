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

# Pick a model once and save it as your default
moki --set-model
```

On first run Moki fetches the OpenRouter model catalog and shows an
interactive picker. The catalog is cached for 24 hours.

**Your model choice is remembered.** Whenever you pick a model — at startup,
with `--set-model`, or via `/model` in a conversation — Moki saves it and
reuses it next time, so you only choose once. The picker marks your saved
model with a ★ and pre-selects it.

After a one-shot answer, Moki prints the model that replied and the tokens
used:

```
— anthropic/claude-sonnet-5.5 · 412 tokens
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

Inside `moki -c`:

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

Moki is built on [`ai-util`](https://github.com/ztkent/ai-util), a small
OpenRouter client with streaming and tool-calling support.

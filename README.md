# nlm - NotebookLM CLI

An unofficial CLI client for [Google NotebookLM](https://notebooklm.google.com), written in Go.

Manage notebooks, add sources, chat with AI, and generate content — all from your terminal.

## Why?

The existing Python client ([notebooklm-py](https://github.com/teng-lin/notebooklm-py)) works well, but deploying Python across Windows, macOS, and Linux requires managing interpreters, virtual environments, and dependencies.

**nlm** is a single static binary. Download it, run it. No runtime dependencies, no setup friction — ideal for CI/CD pipelines, automation scripts, and cross-platform distribution.

## Installation

### Pre-built binaries

Download from [Releases](https://github.com/Dokkabei97/notebooklm-cli/releases) for your platform.

### Build from source

```bash
git clone https://github.com/Dokkabei97/notebooklm-cli.git
cd notebooklm-cli
make build        # produces ./nlm binary

# or install to $GOPATH/bin
go install github.com/Dokkabei97/notebooklm-cli@latest
```

### Cross-compile

```bash
GOOS=windows GOARCH=amd64 go build -o nlm.exe .
GOOS=linux   GOARCH=amd64 go build -o nlm-linux .
GOOS=darwin  GOARCH=arm64 go build -o nlm-darwin .
```

## Quick Start

```bash
# 1. Authenticate (reuses existing Chrome login — no browser window opened)
nlm auth login --reuse

# 2. List notebooks
nlm notebook list

# 3. Select a notebook
nlm use <notebook-id>

# 4. Ask a question
nlm chat ask "Summarize the key points"
```

## Authentication

### Option 1: Reuse Chrome cookies (recommended)

Extracts cookies directly from Chrome's local database. No browser is opened, and your existing login session is untouched.

```bash
nlm auth login --reuse
```

> On macOS, a Keychain access prompt may appear. Click "Allow".

### Option 2: Browser login

Opens a new Chrome window for Google sign-in.

```bash
nlm auth login
```

### Check / clear auth

```bash
nlm auth status
nlm auth clear
```

## Commands

### Notebooks

```bash
nlm notebook list                        # List notebooks
nlm notebook create "Research Notes"     # Create
nlm notebook rename <id> "New Title"     # Rename
nlm notebook delete <id>                 # Delete
nlm nb ls                                # Alias
```

### Active notebook

```bash
nlm use <notebook-id>
```

Subsequent commands use this notebook by default. Override with `-n <id>`.

### Sources

```bash
nlm source list                          # List sources
nlm source add https://example.com       # Add URL
nlm source add ./paper.pdf               # Upload file
nlm source add <url> --wait              # Add and wait for processing
nlm source get <id>                      # Source details
nlm source refresh <id>                  # Refresh
nlm source delete <id>                   # Delete
nlm src ls                               # Alias
```

### Chat

```bash
nlm chat ask "What are the main arguments?"
nlm chat ask "Focus on chapter 3" -s <source-id>
nlm chat history
```

### Generate content

```bash
nlm generate audio                       # Audio overview
nlm generate audio -i "Focus on X"       # With instructions
nlm generate report                      # Report
nlm generate quiz                        # Quiz
nlm generate video                       # Video
nlm generate mind-map                    # Mind map
nlm generate infographic                 # Infographic
nlm generate slide-deck                  # Slide deck
nlm gen audio --wait                     # Wait for completion
```

### Artifacts

```bash
nlm artifact list
nlm artifact export <id>                 # Print content
nlm artifact export <id> output.md       # Save to file
nlm artifact delete <id>
```

### Notes

```bash
nlm note list
nlm note create "Title" "Content"
nlm note update <id> "New Title" "New Content"
nlm note delete <id>
```

### Deep Research

```bash
nlm research start "Investigate X"
nlm research poll <research-id>
nlm research import <research-id>        # Import results as note
```

### Sharing

```bash
nlm share status
nlm share set viewer
nlm share set editor
nlm share set none
```

## Global Flags

| Flag | Description |
|------|-------------|
| `--json` | JSON output (useful for scripting) |
| `-v, --verbose` | Verbose logging |
| `--config <path>` | Config file path |
| `-n, --notebook <id>` | Specify notebook ID |

```bash
nlm notebook list --json
nlm source list --json | jq '.[].title'
```

## Shell Completion

```bash
# Bash
source <(nlm completion bash)

# Zsh
nlm completion zsh > "${fpath[1]}/_nlm"

# Fish
nlm completion fish | source
```

## Configuration

Auth and settings are stored in `~/.notebooklm/`.

| File | Description |
|------|-------------|
| `~/.notebooklm/storage_state.json` | Auth cookies |
| `~/.notebooklm/config.json` | Active notebook, preferences |

You can also pass auth via the `NOTEBOOKLM_AUTH_JSON` environment variable.

## Requirements

- **macOS**: Full support (Chrome cookie extraction via Keychain)
- **Windows/Linux**: Browser login (`nlm auth login`) or manual cookie import
- Google Chrome with a logged-in Google account
- Go 1.21+ (build only)

## Architecture

```
cmd/           CLI commands (Cobra)
internal/
  rpc/         batchexecute protocol (encoder, decoder, caller)
  auth/        Authentication (browser login, Chrome cookie extraction)
  api/         Business logic (notebooks, sources, chat, artifacts, ...)
  model/       Domain models
  config/      Config management (~/.notebooklm/)
  output/      Terminal output (lipgloss tables, JSON)
```

## References

- **[notebooklm-py](https://github.com/teng-lin/notebooklm-py)** — Python client by teng-lin. The protocol implementation (RPC method IDs, request/response encoding, auth flow) is based on this project.
- **[Google batchexecute](https://kovatch.medium.com/deciphering-google-batchexecute-74991e4e446c)** — Protocol documentation for Google's internal RPC mechanism.

## Contributing

Contributions are welcome! Please follow these guidelines:

1. **Fork** the repository and create a feature branch.
2. **Write tests** for new functionality. Run `make e2e` to verify.
3. **Follow the existing code style** — `go vet` and `go build` must pass cleanly.
4. **Submit a pull request** with a clear description of the change.

### Development

```bash
make build       # Build binary
make test        # Run unit tests
make e2e-basic   # Run CLI tests (no auth needed)
make e2e         # Run full E2E tests (requires auth)
```

### Reporting Issues

- Include the command you ran and the full error output.
- Specify your OS, Go version, and Chrome version.
- For auth issues, run `nlm auth status` and include the output (tokens are truncated).

## License

MIT License

Copyright (c) 2025

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.

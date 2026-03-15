---
name: notebooklm
description: |
  Research assistant powered by NotebookLM CLI (nlm).
  Add code/documents as sources to NotebookLM, run AI queries, deep research, and content generation.
  Triggered by keywords: 'NotebookLM', 'nlm', 'research'.
---

# NotebookLM CLI Integration

A skill for source-based research, AI queries, and content generation via Google NotebookLM.

## When to Apply

- User mentions "NotebookLM", "nlm", or asks to research with NotebookLM
- User wants to upload documents/code to NotebookLM for analysis
- User wants to ask source-based questions to NotebookLM AI
- User wants to run deep research and retrieve results
- User wants to generate content (audio, report, quiz, mind map, etc.)

## Prerequisites

1. **Install**: `cd <notebooklm-cli-project> && make install`
2. **Authenticate**: `nlm auth login` (browser login) or `nlm auth login --reuse` (reuse Chrome cookies)
3. **Verify**: `nlm auth status`

## Source Strategy (Hybrid)

Use a hybrid strategy when adding sources:

1. **User-specified sources exist** → Run `nlm source add` immediately
2. **No sources specified** → Explore current context for relevant files
   - Search for README, design docs, core source code via Glob/Read
   - Present candidate list to user and add after confirmation
3. **Always use `--wait` flag** → Ensure source processing completes

## Core Workflows

### 1. Research Workflow (Source → Ask → Research)

```bash
# Step 1: Create or select a notebook
nlm notebook create "Project Research" --json
nlm notebook list --json

# Step 2: Set active notebook
nlm use <notebook-id>

# Step 3: Add sources (hybrid strategy)
# When user specifies URL/file:
nlm source add "https://example.com/article" --wait --json
nlm source add ./document.pdf --wait --json
# When not specified: explore current project and suggest relevant files

# Step 4: Ask AI questions
nlm chat ask "Summarize the key points from these sources" --json

# Step 5: Run deep research (optional)
nlm research start "In-depth analysis on the topic" --json
nlm research poll <research-id> --json
nlm research import <research-id>
```

### 2. Content Generation Workflow

```bash
# Generate content from a notebook with sources
nlm generate audio --wait --json          # Audio overview
nlm generate report --wait --json         # Report
nlm generate quiz --wait --json           # Quiz
nlm generate mind-map --wait --json       # Mind map
nlm generate video --wait --json          # Video
nlm generate infographic --wait --json    # Infographic
nlm generate slide-deck --wait --json     # Slide deck

# Add custom instructions for audio
nlm generate audio --instructions "conversational style" --wait --json
```

### 3. Note Management Workflow

```bash
nlm note list --json
nlm note create "Title" "Content" --json
nlm note update <note-id> "New Title" "New Content" --json
nlm note delete <note-id>
```

## Command Reference

| Command | Description | Key Flags |
|---------|-------------|-----------|
| `nlm auth login` | Google login | `--reuse` |
| `nlm auth status` | Check auth status | |
| `nlm notebook list` | List notebooks | `--json` |
| `nlm notebook create <title>` | Create notebook | `--json` |
| `nlm use <id>` | Set active notebook | |
| `nlm source add <url\|file>` | Add source | `--wait`, `-n` |
| `nlm source list` | List sources | `-n`, `--json` |
| `nlm chat ask <question>` | AI query | `-n`, `-s`, `--json` |
| `nlm chat history` | Chat history | `-n`, `--json` |
| `nlm research start <query>` | Start deep research | `-n`, `--json` |
| `nlm research poll <id>` | Check research progress | `-n`, `--json` |
| `nlm generate <type>` | Generate content | `--wait`, `-n`, `--json` |
| `nlm note create <t> <c>` | Create note | `-n`, `--json` |
| `nlm share set <perm>` | Set sharing | `-n` |

## JSON Output

Add the `--json` flag to any command for structured JSON output.
Use this to extract IDs and chain them into subsequent commands.

```bash
# Example: Create notebook, extract ID, then add source
NB_ID=$(nlm notebook create "Research" --json | jq -r '.id')
nlm use "$NB_ID"
nlm source add "https://example.com" --wait --json
```

## Key Flags

- `--json`: JSON output (essential for parsing/automation)
- `--wait` / `-w`: Wait for source processing or content generation to complete
- `--notebook` / `-n <id>`: Specify notebook (alternative to `nlm use`)
- `--sources` / `-s <ids>`: Reference specific sources for chat
- `--verbose` / `-v`: Verbose logging

## Error Handling

| Error | Cause | Resolution |
|-------|-------|------------|
| `no saved credentials` | Not authenticated | `nlm auth login --reuse` |
| `authentication expired` | Token expired | `nlm auth login --reuse` |
| `please specify a notebook` | No active notebook set | `nlm use <id>` or `-n` flag |
| `source processing timeout` | Source processing delayed | Retry with `nlm source wait <id>` |

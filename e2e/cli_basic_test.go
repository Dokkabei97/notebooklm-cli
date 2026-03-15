package e2e

import (
	"strings"
	"testing"
)

// ============================================================
// Basic CLI behavior tests (no authentication required)
// ============================================================

func TestVersion(t *testing.T) {
	result := nlm(t, "version")
	assertSuccess(t, result)
	assertContains(t, result, "nlm")
}

func TestHelp(t *testing.T) {
	result := nlm(t, "--help")
	assertSuccess(t, result)
	assertContains(t, result, "NotebookLM")
	assertContains(t, result, "Available Commands")
}

func TestSubcommandHelp(t *testing.T) {
	subcommands := []string{"auth", "notebook", "source", "chat", "generate", "note", "research", "share", "artifact"}
	for _, sub := range subcommands {
		t.Run(sub, func(t *testing.T) {
			result := nlm(t, sub, "--help")
			assertSuccess(t, result)
			assertContains(t, result, "Available Commands")
		})
	}
}

func TestCompletionBash(t *testing.T) {
	result := nlm(t, "completion", "bash")
	assertSuccess(t, result)
	assertContains(t, result, "bash")
}

func TestCompletionZsh(t *testing.T) {
	result := nlm(t, "completion", "zsh")
	assertSuccess(t, result)
}

func TestUnknownCommand(t *testing.T) {
	result := nlm(t, "nonexistent")
	assertExitCode(t, result, 1)
}

func TestUseRequiresArg(t *testing.T) {
	result := nlm(t, "use")
	assertExitCode(t, result, 1)
}

func TestNotebookAliases(t *testing.T) {
	// Verify nb alias works
	result := nlm(t, "nb", "--help")
	assertSuccess(t, result)
	assertContains(t, result, "notebook")
}

func TestSourceAliases(t *testing.T) {
	result := nlm(t, "src", "--help")
	assertSuccess(t, result)
	assertContains(t, result, "source")
}

func TestGenerateAliases(t *testing.T) {
	result := nlm(t, "gen", "--help")
	assertSuccess(t, result)
	assertContains(t, result, "Generate content")
}

func TestGlobalFlags(t *testing.T) {
	result := nlm(t, "--help")
	assertSuccess(t, result)

	// Verify all global flags are present
	for _, flag := range []string{"--json", "--verbose", "--config"} {
		if !strings.Contains(result.stdout, flag) {
			t.Errorf("expected help to contain flag %q", flag)
		}
	}
}

package e2e

import (
	"testing"
)

// ============================================================
// Chat E2E tests (actual API calls)
// ============================================================

func TestChatAsk(t *testing.T) {
	requireAuth(t)

	nbID := getFirstNotebookID(t)
	if nbID == "" {
		t.Skip("no notebooks found")
	}

	// Simple question - verify a response is returned
	result := nlm(t, "chat", "ask", "Describe this notebook in one sentence", "-n", nbID)
	assertSuccess(t, result)

	// Verify the response is not empty
	if len(result.stdout) < 10 {
		t.Errorf("response too short: %q", result.stdout)
	}
}

func TestChatAskJSON(t *testing.T) {
	requireAuth(t)

	nbID := getFirstNotebookID(t)
	if nbID == "" {
		t.Skip("no notebooks found")
	}

	result := nlmJSON(t, "chat", "ask", "hello", "-n", nbID)
	assertSuccess(t, result)

	// Verify JSON response
	assertContains(t, result, "answer")
}

func TestChatRequiresQuestion(t *testing.T) {
	result := nlm(t, "chat", "ask")
	assertExitCode(t, result, 1)
}

package e2e

import (
	"bytes"
	"os"
	"os/exec"
	"strings"
	"testing"
)

// nlmBinary is the path to the compiled nlm binary.
var nlmBinary = "./nlm"

func init() {
	// Allow path override in CI, etc.
	if p := os.Getenv("NLM_BINARY"); p != "" {
		nlmBinary = p
	}
}

// nlmResult holds the output of a CLI invocation.
type nlmResult struct {
	stdout   string
	stderr   string
	exitCode int
}

// nlm runs the nlm binary with the given arguments and returns the result.
func nlm(t *testing.T, args ...string) nlmResult {
	t.Helper()
	cmd := exec.Command(nlmBinary, args...)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()

	exitCode := 0
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			exitCode = exitErr.ExitCode()
		}
	}

	return nlmResult{
		stdout:   stdout.String(),
		stderr:   stderr.String(),
		exitCode: exitCode,
	}
}

// nlmJSON runs nlm with --json flag and returns the result.
func nlmJSON(t *testing.T, args ...string) nlmResult {
	t.Helper()
	return nlm(t, append(args, "--json")...)
}

// requireAuth skips the test if auth is not configured.
func requireAuth(t *testing.T) {
	t.Helper()
	result := nlm(t, "auth", "status")
	if result.exitCode != 0 || !strings.Contains(result.stdout, "Authentication valid") {
		t.Skip("Authentication not configured. Run 'nlm auth login --reuse' first.")
	}
}

// assertSuccess checks that the command exited with code 0.
func assertSuccess(t *testing.T, result nlmResult) {
	t.Helper()
	if result.exitCode != 0 {
		t.Fatalf("command failed (exit %d)\nstdout: %s\nstderr: %s", result.exitCode, result.stdout, result.stderr)
	}
}

// assertContains checks that stdout contains the given substring.
func assertContains(t *testing.T, result nlmResult, substr string) {
	t.Helper()
	if !strings.Contains(result.stdout, substr) {
		t.Errorf("expected stdout to contain %q, got:\n%s", substr, result.stdout)
	}
}

// assertExitCode checks the exit code.
func assertExitCode(t *testing.T, result nlmResult, code int) {
	t.Helper()
	if result.exitCode != code {
		t.Errorf("expected exit code %d, got %d\nstdout: %s\nstderr: %s", code, result.exitCode, result.stdout, result.stderr)
	}
}

package e2e

import (
	"testing"
)

// ============================================================
// Authentication E2E tests
// ============================================================

func TestAuthStatusWhenAuthenticated(t *testing.T) {
	requireAuth(t)

	result := nlm(t, "auth", "status")
	assertSuccess(t, result)
	assertContains(t, result, "Authentication valid")
	assertContains(t, result, "Cookies")
	assertContains(t, result, "CSRF Token")
	assertContains(t, result, "Session ID")
}

func TestAuthLoginReuse(t *testing.T) {
	// --reuse extracts directly from Chrome cookie DB
	result := nlm(t, "auth", "login", "--reuse")
	assertSuccess(t, result)
	assertContains(t, result, "Authentication successful")
}

func TestAuthLoginHelp(t *testing.T) {
	result := nlm(t, "auth", "login", "--help")
	assertSuccess(t, result)
	assertContains(t, result, "--reuse")
}

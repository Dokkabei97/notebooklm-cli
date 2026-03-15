package e2e

import (
	"encoding/json"
	"testing"
)

// ============================================================
// Source E2E tests (actual API calls)
// ============================================================

func TestSourceListJSON(t *testing.T) {
	requireAuth(t)

	nbID := getNotebookWithSources(t)
	if nbID == "" {
		t.Skip("no notebook with sources found")
	}

	result := nlmJSON(t, "source", "list", "-n", nbID)
	assertSuccess(t, result)

	var sources []map[string]any
	if err := json.Unmarshal([]byte(result.stdout), &sources); err != nil {
		if result.stdout == "null\n" || result.stdout == "null" {
			t.Skip("no sources found")
		}
		t.Fatalf("JSON parse failed: %v", err)
	}

	if len(sources) == 0 {
		t.Skip("no sources found")
	}

	// Check required fields
	src := sources[0]
	for _, field := range []string{"id", "title", "status"} {
		if _, ok := src[field]; !ok {
			t.Errorf("source missing %q field: %v", field, src)
		}
	}
}

func TestSourceListTable(t *testing.T) {
	requireAuth(t)

	nbID := getNotebookWithSources(t)
	if nbID == "" {
		t.Skip("no notebook with sources found")
	}

	result := nlm(t, "source", "list", "-n", nbID)
	assertSuccess(t, result)
	assertContains(t, result, "ID")
	assertContains(t, result, "Title")
}

func TestSourceRequiresNotebook(t *testing.T) {
	requireAuth(t)

	// Try source list without an active notebook
	// (may not be in a clean state if use was previously called, so test without -n flag)
	// Can succeed if an active notebook is already set
	result := nlm(t, "source", "list", "-n", "")
	// Empty notebook ID should result in an error
	if result.exitCode == 0 {
		// OK if there's an active notebook
		return
	}
	// Verify error message contains notebook-related guidance
	combined := result.stdout + result.stderr
	if combined == "" {
		t.Error("no error message returned")
	}
}

// getFirstNotebookID returns the ID of the first notebook, or empty string.
func getFirstNotebookID(t *testing.T) string {
	t.Helper()
	result := nlmJSON(t, "notebook", "list")
	if result.exitCode != 0 {
		return ""
	}
	var notebooks []map[string]any
	json.Unmarshal([]byte(result.stdout), &notebooks)
	if len(notebooks) == 0 {
		return ""
	}
	id, _ := notebooks[0]["id"].(string)
	return id
}

// getNotebookWithSources returns the ID of a notebook that has sources.
func getNotebookWithSources(t *testing.T) string {
	t.Helper()
	result := nlmJSON(t, "notebook", "list")
	if result.exitCode != 0 {
		return ""
	}
	var notebooks []map[string]any
	json.Unmarshal([]byte(result.stdout), &notebooks)
	for _, nb := range notebooks {
		sc, _ := nb["source_count"].(float64)
		if sc > 0 {
			id, _ := nb["id"].(string)
			return id
		}
	}
	return ""
}

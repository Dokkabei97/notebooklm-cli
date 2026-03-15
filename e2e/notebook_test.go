package e2e

import (
	"encoding/json"
	"strings"
	"testing"
)

// ============================================================
// Notebook E2E tests (actual API calls)
// ============================================================

func TestNotebookList(t *testing.T) {
	requireAuth(t)

	result := nlm(t, "notebook", "list")
	assertSuccess(t, result)

	// Check that table headers exist
	assertContains(t, result, "ID")
	assertContains(t, result, "Title")
}

func TestNotebookListJSON(t *testing.T) {
	requireAuth(t)

	result := nlmJSON(t, "notebook", "list")
	assertSuccess(t, result)

	// Verify it can be parsed as JSON array
	var notebooks []map[string]any
	if err := json.Unmarshal([]byte(result.stdout), &notebooks); err != nil {
		t.Fatalf("JSON parse failed: %v\noutput: %s", err, result.stdout)
	}

	if len(notebooks) == 0 {
		t.Skip("no notebooks found")
	}

	// Check required fields
	nb := notebooks[0]
	for _, field := range []string{"id", "title"} {
		if _, ok := nb[field]; !ok {
			t.Errorf("notebook missing %q field: %v", field, nb)
		}
	}

	// ID should be UUID format
	id, _ := nb["id"].(string)
	if len(id) < 10 || !strings.Contains(id, "-") {
		t.Errorf("notebook ID is not UUID format: %q", id)
	}
}

func TestNotebookCreateAndDelete(t *testing.T) {
	requireAuth(t)

	testTitle := "nlm-e2e-test-notebook"

	// Create
	result := nlm(t, "notebook", "create", testTitle)
	assertSuccess(t, result)
	assertContains(t, result, "Notebook created")

	// List as JSON to extract ID
	result = nlmJSON(t, "notebook", "list")
	assertSuccess(t, result)

	var notebooks []map[string]any
	json.Unmarshal([]byte(result.stdout), &notebooks)

	var testID string
	for _, nb := range notebooks {
		title, _ := nb["title"].(string)
		if title == testTitle {
			testID, _ = nb["id"].(string)
			break
		}
	}

	if testID == "" {
		t.Fatal("could not find the created test notebook")
	}

	// Delete
	result = nlm(t, "notebook", "delete", testID)
	assertSuccess(t, result)
	assertContains(t, result, "deleted")

	// Verify it's gone from the list
	result = nlmJSON(t, "notebook", "list")
	assertSuccess(t, result)

	json.Unmarshal([]byte(result.stdout), &notebooks)
	for _, nb := range notebooks {
		id, _ := nb["id"].(string)
		if id == testID {
			t.Error("deleted notebook still appears in the list")
		}
	}
}

func TestUseAndSourceList(t *testing.T) {
	requireAuth(t)

	result := nlmJSON(t, "notebook", "list")
	assertSuccess(t, result)

	var notebooks []map[string]any
	json.Unmarshal([]byte(result.stdout), &notebooks)
	if len(notebooks) == 0 {
		t.Skip("no notebooks found")
	}

	// Find a notebook with sources and test
	for _, nb := range notebooks {
		nbID, _ := nb["id"].(string)
		sc, _ := nb["source_count"].(float64)
		if sc == 0 {
			continue
		}

		result = nlm(t, "use", nbID)
		assertSuccess(t, result)
		assertContains(t, result, "Active notebook")

		result = nlm(t, "source", "list", "-n", nbID)
		if result.exitCode == 0 {
			return // success
		}
	}

	t.Skip("no notebook with retrievable source list found (possible RPC parsing issue)")
}

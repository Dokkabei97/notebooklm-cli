package api

import (
	"github.com/Dokkabei97/notebooklm-cli/internal/rpc"
)

// ResearchResult holds the result of a deep research operation.
type ResearchResult struct {
	ID       string `json:"id"`
	Status   string `json:"status"`
	Progress int    `json:"progress"`
	Content  string `json:"content,omitempty"`
}

// StartResearch initiates a deep research task.
func (c *Client) StartResearch(notebookID, query string) (*ResearchResult, error) {
	params := []any{query}
	result, err := c.caller.Call(rpc.MethodStartDeepResearch, params, notebookPath(notebookID))
	if err != nil {
		return nil, wrapErr("StartResearch", "rpc call failed", err)
	}

	r := &ResearchResult{Status: "started"}
	if result != nil {
		r.ID = rpc.SafeString(result.Parsed, 0)
	}
	return r, nil
}

// PollResearch checks the progress of a research task.
func (c *Client) PollResearch(notebookID, researchID string) (*ResearchResult, error) {
	params := []any{researchID}
	result, err := c.caller.Call(rpc.MethodPollResearch, params, notebookPath(notebookID))
	if err != nil {
		return nil, wrapErr("PollResearch", "rpc call failed", err)
	}

	r := &ResearchResult{ID: researchID}
	if result == nil {
		r.Status = "unknown"
		return r, nil
	}

	statusVal := rpc.SafeFloat(result.Parsed, 0)
	switch int(statusVal) {
	case 1:
		r.Status = "in_progress"
		r.Progress = int(rpc.SafeFloat(result.Parsed, 2))
	case 2:
		r.Status = "completed"
		r.Progress = 100
		r.Content = rpc.SafeString(result.Parsed, 1)
	case 3:
		r.Status = "error"
	default:
		r.Status = "unknown"
	}

	return r, nil
}

// ImportResearch imports research results as a source.
func (c *Client) ImportResearch(notebookID, researchID string) error {
	params := []any{researchID}
	_, err := c.caller.Call(rpc.MethodImportResearch, params, notebookPath(notebookID))
	if err != nil {
		return wrapErr("ImportResearch", "rpc call failed", err)
	}
	return nil
}

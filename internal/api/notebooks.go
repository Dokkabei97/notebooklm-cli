package api

import (
	"time"

	"github.com/Dokkabei97/notebooklm-cli/internal/model"
	"github.com/Dokkabei97/notebooklm-cli/internal/rpc"
)

// ListNotebooks returns all notebooks for the authenticated user.
func (c *Client) ListNotebooks() ([]model.Notebook, error) {
	params := []any{nil, 1, nil, []any{2}}
	result, err := c.caller.Call(rpc.MethodListNotebooks, params, "/")
	if err != nil {
		return nil, wrapErr("ListNotebooks", "rpc call failed", err)
	}
	if result == nil {
		return nil, nil // empty
	}
	return parseNotebooks(result.Parsed)
}

// CreateNotebook creates a new notebook with the given title.
func (c *Client) CreateNotebook(title string) (*model.Notebook, error) {
	params := []any{title, nil, nil, []any{2}, []any{1}}
	result, err := c.caller.Call(rpc.MethodCreateNotebook, params, "/")
	if err != nil {
		return nil, wrapErr("CreateNotebook", "rpc call failed", err)
	}
	// null result = success but no response data, return minimal info
	if result == nil {
		return &model.Notebook{Title: title}, nil
	}
	return parseNotebook(result.Parsed)
}

// DeleteNotebook deletes a notebook by ID.
func (c *Client) DeleteNotebook(notebookID string) error {
	params := []any{[]any{notebookID}, []any{2}}
	_, err := c.caller.Call(rpc.MethodDeleteNotebook, params, "/")
	// null result = void success
	return err
}

// RenameNotebook renames a notebook.
func (c *Client) RenameNotebook(notebookID, newTitle string) error {
	params := []any{notebookID, []any{[]any{nil, nil, nil, []any{nil, newTitle}}}}
	_, err := c.caller.Call(rpc.MethodRenameNotebook, params, "/")
	return err
}

// GetNotebook gets details about a specific notebook.
func (c *Client) GetNotebook(notebookID string) (*model.Notebook, error) {
	params := []any{notebookID, nil, []any{2}, nil, 0}
	result, err := c.caller.Call(rpc.MethodGetNotebook, params, notebookPath(notebookID))
	if err != nil {
		return nil, wrapErr("GetNotebook", "rpc call failed", err)
	}
	if result == nil {
		return &model.Notebook{ID: notebookID}, nil
	}
	if result.Parsed != nil && len(result.Parsed) > 0 {
		nbInfo := rpc.SafeArray(result.Parsed, 0)
		if nbInfo != nil {
			return parseNotebookFromArray(nbInfo)
		}
	}
	return parseNotebook(result.Parsed)
}

func parseNotebooks(data []any) ([]model.Notebook, error) {
	if data == nil {
		return nil, nil
	}

	items := rpc.SafeArray(data, 0)
	if items == nil {
		return nil, nil
	}

	var notebooks []model.Notebook
	for _, item := range items {
		arr, ok := item.([]any)
		if !ok {
			continue
		}
		nb, err := parseNotebookFromArray(arr)
		if err != nil {
			continue
		}
		notebooks = append(notebooks, *nb)
	}

	return notebooks, nil
}

func parseNotebook(data []any) (*model.Notebook, error) {
	if data == nil {
		return nil, wrapErr("parseNotebook", "nil response", nil)
	}
	return parseNotebookFromArray(data)
}

func parseNotebookFromArray(arr []any) (*model.Notebook, error) {
	// Response structure: [title, sources_list, id, ...]
	nb := &model.Notebook{
		Title: rpc.SafeString(arr, 0),
		ID:    rpc.SafeString(arr, 2),
	}

	// Source count
	sourcesList := rpc.SafeArray(arr, 1)
	if sourcesList != nil {
		nb.SourceCount = len(sourcesList)
	}

	// Timestamps: try multiple possible positions
	for _, idx := range []int{3, 4, 5, 6} {
		if idx >= len(arr) {
			break
		}
		// [seconds, nanoseconds] format or microseconds
		if tsArr := rpc.SafeArray(arr, idx); tsArr != nil && len(tsArr) > 0 {
			if ts, ok := tsArr[0].(float64); ok && ts > 1e9 {
				if nb.CreatedAt.IsZero() {
					nb.CreatedAt = time.Unix(int64(ts), 0)
				} else if nb.UpdatedAt.IsZero() {
					nb.UpdatedAt = time.Unix(int64(ts), 0)
				}
			}
		}
		if ts := rpc.SafeFloat(arr, idx); ts > 1e15 {
			if nb.CreatedAt.IsZero() {
				nb.CreatedAt = time.UnixMicro(int64(ts))
			} else if nb.UpdatedAt.IsZero() {
				nb.UpdatedAt = time.UnixMicro(int64(ts))
			}
		}
	}

	return nb, nil
}

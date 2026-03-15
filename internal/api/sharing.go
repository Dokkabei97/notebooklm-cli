package api

import (
	"github.com/jmk/notebooklm-cli/internal/model"
	"github.com/jmk/notebooklm-cli/internal/rpc"
)

// ShareStatus holds the sharing state of a notebook.
type ShareStatus struct {
	IsShared   bool                  `json:"is_shared"`
	Permission model.SharePermission `json:"permission"`
	ShareURL   string                `json:"share_url,omitempty"`
}

// SetSharing sets the sharing permissions for a notebook.
func (c *Client) SetSharing(notebookID string, permission model.SharePermission) error {
	params := []any{int(permission)}
	_, err := c.caller.Call(rpc.MethodShareNotebook, params, notebookPath(notebookID))
	if err != nil {
		return wrapErr("SetSharing", "rpc call failed", err)
	}
	return nil
}

// GetSharing returns the current sharing status.
func (c *Client) GetSharing(notebookID string) (*ShareStatus, error) {
	result, err := c.caller.Call(rpc.MethodGetShareStatus, nil, notebookPath(notebookID))
	if err != nil {
		return nil, wrapErr("GetSharing", "rpc call failed", err)
	}

	status := &ShareStatus{}
	if result == nil {
		return status, nil
	}

	if rpc.SafeFloat(result.Parsed, 0) > 0 {
		status.IsShared = true
		status.Permission = model.SharePermission(int(rpc.SafeFloat(result.Parsed, 0)))
	}
	status.ShareURL = rpc.SafeString(result.Parsed, 1)

	return status, nil
}

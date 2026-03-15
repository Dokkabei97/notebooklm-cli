package rpc

import (
	"encoding/json"
	"fmt"
	"net/url"
)

// EncodeBatchRequest builds the form-encoded body for a batchexecute call.
// Format: f.req=[[[rpc_id, json_params, null, "generic"]]]&at=csrf_token&
func EncodeBatchRequest(rpcID string, params any, csrfToken string) (string, error) {
	// JSON-encode params compactly (no spaces, matching Chrome behavior)
	paramsJSON, err := json.Marshal(params)
	if err != nil {
		return "", fmt.Errorf("marshal params: %w", err)
	}

	// Inner request: [rpc_id, json_params_string, null, "generic"]
	inner := []any{rpcID, string(paramsJSON), nil, "generic"}

	// Triple-nest: [[inner]]
	outer := [][]any{{inner}}
	outerJSON, err := json.Marshal(outer)
	if err != nil {
		return "", fmt.Errorf("marshal outer: %w", err)
	}

	// URL-encode and build body with trailing &
	encoded := url.QueryEscape(string(outerJSON))
	body := "f.req=" + encoded

	if csrfToken != "" {
		body += "&at=" + url.QueryEscape(csrfToken)
	}

	body += "&"
	return body, nil
}

// BuildBatchURL constructs the full URL for a batchexecute call.
func BuildBatchURL(rpcID, sourcePath, sessionID string) string {
	params := url.Values{}
	params.Set("rpcids", rpcID)
	params.Set("source-path", sourcePath)
	if sessionID != "" {
		params.Set("f.sid", sessionID)
	}
	params.Set("rt", "c")
	return BatchExecuteURL + "?" + params.Encode()
}

// BuildChatURL constructs the URL for chat/query endpoint.
func BuildChatURL(sessionID string, reqID int) string {
	params := url.Values{}
	params.Set("bl", "boq_labs-tailwind-frontend_20260301.03_p0")
	params.Set("hl", "en")
	params.Set("_reqid", fmt.Sprintf("%d", reqID))
	params.Set("rt", "c")
	if sessionID != "" {
		params.Set("f.sid", sessionID)
	}
	return QueryURL + "?" + params.Encode()
}

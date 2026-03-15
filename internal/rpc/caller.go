package rpc

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

// Caller handles RPC calls to the NotebookLM batchexecute API.
type Caller struct {
	HTTPClient   *http.Client
	Cookies      []*http.Cookie
	CSRFToken    string
	SessionID    string
	cookieHeader string
}

// NewCaller creates a new RPC caller.
func NewCaller(cookies []*http.Cookie, csrfToken, sessionID string) *Caller {
	// Build Cookie header once (set directly instead of using AddCookie)
	var parts []string
	seen := make(map[string]bool)
	for _, c := range cookies {
		if c.Value == "" || seen[c.Name] {
			continue
		}
		seen[c.Name] = true
		parts = append(parts, c.Name+"="+c.Value)
	}

	return &Caller{
		HTTPClient:   &http.Client{},
		Cookies:      cookies,
		CSRFToken:    csrfToken,
		SessionID:    sessionID,
		cookieHeader: strings.Join(parts, "; "),
	}
}

// CallResult holds the parsed result of an RPC call.
type CallResult struct {
	Raw    json.RawMessage
	Parsed []any
}

// Call executes a single batchexecute RPC call.
// Returns (nil, nil) for void operations where Google returns null result_data.
func (c *Caller) Call(rpcID string, params any, notebookPath string) (*CallResult, error) {
	body, err := EncodeBatchRequest(rpcID, params, c.CSRFToken)
	if err != nil {
		return nil, fmt.Errorf("encode request: %w", err)
	}

	reqURL := BuildBatchURL(rpcID, notebookPath, c.SessionID)
	req, err := http.NewRequest("POST", reqURL, strings.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Origin", "https://notebooklm.google.com")
	req.Header.Set("Referer", "https://notebooklm.google.com/")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/131.0.0.0 Safari/537.36")
	req.Header.Set("X-Same-Domain", "1")
	req.Header.Set("Cookie", c.cookieHeader)

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("http request: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response: %w", err)
	}

	if err := checkHTTPStatus(resp.StatusCode, rpcID); err != nil {
		return nil, err
	}

	raw, err := DecodeBatchResponse(string(respBody), rpcID)
	if err != nil {
		return nil, err
	}

	// null result (void operation / empty list) → success with nil
	if raw == nil {
		return nil, nil
	}

	parsed, err := ParseResultArray(raw)
	if err != nil {
		return &CallResult{Raw: raw}, nil
	}

	return &CallResult{Raw: raw, Parsed: parsed}, nil
}

func checkHTTPStatus(statusCode int, rpcID string) error {
	switch {
	case statusCode == 200:
		return nil
	case statusCode == 401 || statusCode == 403:
		return &Error{Code: ErrorCode(statusCode), Message: "authentication required or expired", Method: rpcID}
	case statusCode == 404:
		return &Error{Code: ErrNotFound, Message: "resource not found", Method: rpcID}
	case statusCode == 429:
		return &Error{Code: ErrRateLimit, Message: "rate limited, please retry later", Method: rpcID}
	case statusCode >= 500:
		return &Error{Code: ErrServer, Message: fmt.Sprintf("server error (HTTP %d)", statusCode), Method: rpcID}
	default:
		return &Error{Code: ErrUnknown, Message: fmt.Sprintf("unexpected HTTP status %d", statusCode), Method: rpcID}
	}
}

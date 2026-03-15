package api

import (
	"github.com/Dokkabei97/notebooklm-cli/internal/auth"
	"github.com/Dokkabei97/notebooklm-cli/internal/rpc"
)

// Client is the main API client for NotebookLM.
type Client struct {
	caller *rpc.Caller
	tokens *auth.AuthTokens
}

// NewClient creates a new API client from auth tokens.
func NewClient(tokens *auth.AuthTokens) *Client {
	return &Client{
		caller: rpc.NewCaller(tokens.Cookies, tokens.CSRFToken, tokens.SessionID),
		tokens: tokens,
	}
}

// Authenticate loads saved auth state and creates a client.
func Authenticate() (*Client, error) {
	cookies, err := auth.LoadStorageState()
	if err != nil {
		return nil, wrapErr("auth", "load saved auth state", err)
	}

	tokens, err := auth.ExtractTokens(cookies)
	if err != nil {
		return nil, wrapErr("auth", "extract session tokens", err)
	}

	return NewClient(tokens), nil
}

func notebookPath(notebookID string) string {
	if notebookID == "" {
		return ""
	}
	return "/notebook/" + notebookID
}

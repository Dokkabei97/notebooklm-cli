package auth

import "net/http"

// AuthTokens holds all authentication data needed for API calls.
type AuthTokens struct {
	Cookies   []*http.Cookie `json:"-"`
	CSRFToken string         `json:"csrf_token"` // SNlM0e value
	SessionID string         `json:"session_id"`  // FdrFJe value
}

// IsValid checks if the tokens appear to be usable.
func (a *AuthTokens) IsValid() bool {
	return len(a.Cookies) > 0 && a.CSRFToken != ""
}

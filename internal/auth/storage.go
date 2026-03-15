package auth

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/Dokkabei97/notebooklm-cli/internal/config"
)

// StorageState mirrors Playwright's storage_state.json format.
type StorageState struct {
	Cookies []StorageCookie `json:"cookies"`
}

type StorageCookie struct {
	Name     string  `json:"name"`
	Value    string  `json:"value"`
	Domain   string  `json:"domain"`
	Path     string  `json:"path"`
	Expires  float64 `json:"expires"`
	HTTPOnly bool    `json:"httpOnly"`
	Secure   bool    `json:"secure"`
	SameSite string  `json:"sameSite"`
}

// LoadStorageState loads cookies from a Playwright-format storage_state.json file.
// It checks NOTEBOOKLM_AUTH_JSON env var first, then falls back to the default path.
func LoadStorageState() ([]*http.Cookie, error) {
	var data []byte
	var err error

	// Check env var first
	if envJSON := os.Getenv("NOTEBOOKLM_AUTH_JSON"); envJSON != "" {
		data = []byte(envJSON)
	} else {
		path := config.AuthFile()
		data, err = os.ReadFile(path)
		if err != nil {
			return nil, fmt.Errorf("read auth file %s: %w", path, err)
		}
	}

	var state StorageState
	if err := json.Unmarshal(data, &state); err != nil {
		return nil, fmt.Errorf("parse storage state: %w", err)
	}

	var cookies []*http.Cookie
	for _, sc := range state.Cookies {
		cookie := &http.Cookie{
			Name:     sc.Name,
			Value:    sc.Value,
			Domain:   sc.Domain,
			Path:     sc.Path,
			HttpOnly: sc.HTTPOnly,
			Secure:   sc.Secure,
		}
		if sc.Expires > 0 {
			cookie.Expires = time.Unix(int64(sc.Expires), 0)
		}
		switch sc.SameSite {
		case "Strict":
			cookie.SameSite = http.SameSiteStrictMode
		case "Lax":
			cookie.SameSite = http.SameSiteLaxMode
		case "None":
			cookie.SameSite = http.SameSiteNoneMode
		}
		cookies = append(cookies, cookie)
	}

	return cookies, nil
}

// SaveStorageState saves cookies in Playwright-compatible format.
func SaveStorageState(cookies []*http.Cookie) error {
	if err := config.EnsureDir(); err != nil {
		return err
	}

	var storageCookies []StorageCookie
	for _, c := range cookies {
		sc := StorageCookie{
			Name:     c.Name,
			Value:    c.Value,
			Domain:   c.Domain,
			Path:     c.Path,
			HTTPOnly: c.HttpOnly,
			Secure:   c.Secure,
		}
		if !c.Expires.IsZero() {
			sc.Expires = float64(c.Expires.Unix())
		}
		switch c.SameSite {
		case http.SameSiteStrictMode:
			sc.SameSite = "Strict"
		case http.SameSiteLaxMode:
			sc.SameSite = "Lax"
		case http.SameSiteNoneMode:
			sc.SameSite = "None"
		}
		storageCookies = append(storageCookies, sc)
	}

	state := StorageState{Cookies: storageCookies}
	data, err := json.MarshalIndent(state, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(config.AuthFile(), data, 0600)
}

// ClearStorageState removes the saved auth state.
func ClearStorageState() error {
	path := config.AuthFile()
	if err := os.Remove(path); err != nil && !os.IsNotExist(err) {
		return err
	}
	return nil
}

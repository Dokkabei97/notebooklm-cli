package auth

import (
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"
)

var (
	csrfRegex    = regexp.MustCompile(`"SNlM0e"\s*:\s*"([^"]+)"`)
	sessionRegex = regexp.MustCompile(`"FdrFJe"\s*:\s*"([^"]+)"`)
)

const notebookLMURL = "https://notebooklm.google.com/"

// ExtractTokens fetches the NotebookLM page and extracts CSRF token and session ID.
func ExtractTokens(cookies []*http.Cookie) (*AuthTokens, error) {
	// Build Cookie header directly, same as Python reference.
	// AddCookie checks domain/path/secure attributes and may omit some cookies.
	cookieHeader := BuildCookieHeader(cookies)

	client := &http.Client{
		// Follow redirects but keep Cookie header on each request
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			if len(via) >= 10 {
				return fmt.Errorf("too many redirects")
			}
			req.Header.Set("Cookie", cookieHeader)
			return nil
		},
	}

	req, err := http.NewRequest("GET", notebookLMURL, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Cookie", cookieHeader)
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/131.0.0.0 Safari/537.36")

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("fetch NotebookLM page: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("HTTP %d (final URL: %s)", resp.StatusCode, resp.Request.URL.String())
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read page body: %w", err)
	}

	html := string(body)
	finalURL := resp.Request.URL.String()

	// Detect redirect to login page
	if strings.Contains(finalURL, "accounts.google.com") {
		return nil, fmt.Errorf("authentication expired. Please log in to Google again in Chrome.")
	}

	csrfMatch := csrfRegex.FindStringSubmatch(html)
	if csrfMatch == nil {
		return nil, fmt.Errorf("CSRF token (SNlM0e) not found - cookies may have expired (final URL: %s)", finalURL)
	}

	tokens := &AuthTokens{
		Cookies:   cookies,
		CSRFToken: csrfMatch[1],
	}

	sessionMatch := sessionRegex.FindStringSubmatch(html)
	if sessionMatch != nil {
		tokens.SessionID = sessionMatch[1]
	}

	// Merge new cookies from response
	for _, rc := range resp.Cookies() {
		found := false
		for i, c := range tokens.Cookies {
			if c.Name == rc.Name {
				tokens.Cookies[i] = rc
				found = true
				break
			}
		}
		if !found {
			tokens.Cookies = append(tokens.Cookies, rc)
		}
	}

	return tokens, nil
}

// BuildCookieHeader creates a Cookie header string like Python does:
// "SID=xxx; HSID=yyy; ..."
func BuildCookieHeader(cookies []*http.Cookie) string {
	var parts []string
	seen := make(map[string]bool)
	for _, c := range cookies {
		if c.Value == "" || seen[c.Name] {
			continue
		}
		seen[c.Name] = true
		parts = append(parts, c.Name+"="+c.Value)
	}
	return strings.Join(parts, "; ")
}

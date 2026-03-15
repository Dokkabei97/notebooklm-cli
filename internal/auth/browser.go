package auth

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
	"github.com/go-rod/rod/lib/proto"
)

// BrowserLogin opens a new browser window for the user to log in to Google.
func BrowserLogin() ([]*http.Cookie, error) {
	fmt.Println("Opening a new browser for Google login...")

	path, _ := launcher.LookPath()
	u := launcher.New().
		Bin(path).
		Headless(false).
		Set("disable-gpu").
		MustLaunch()

	browser := rod.New().ControlURL(u).MustConnect()
	defer browser.MustClose()

	page := browser.MustPage("https://accounts.google.com/ServiceLogin?continue=https://notebooklm.google.com")
	fmt.Println("Please log in with your Google account. It will proceed automatically once login is complete...")

	return waitAndExtractCookies(page)
}

func waitAndExtractCookies(page *rod.Page) ([]*http.Cookie, error) {
	deadline := time.Now().Add(5 * time.Minute)
	for time.Now().Before(deadline) {
		info := page.MustInfo()
		if strings.HasPrefix(info.URL, "https://notebooklm.google.com") {
			time.Sleep(2 * time.Second)
			break
		}
		time.Sleep(1 * time.Second)
	}

	protoCookies, err := proto.NetworkGetAllCookies{}.Call(page)
	if err != nil {
		return nil, fmt.Errorf("failed to extract cookies: %w", err)
	}

	var cookies []*http.Cookie
	for _, pc := range protoCookies.Cookies {
		cookie := &http.Cookie{
			Name:     pc.Name,
			Value:    pc.Value,
			Domain:   pc.Domain,
			Path:     pc.Path,
			HttpOnly: pc.HTTPOnly,
			Secure:   pc.Secure,
		}
		if pc.Expires > 0 {
			cookie.Expires = time.Unix(int64(pc.Expires), 0)
		}
		cookies = append(cookies, cookie)
	}

	if len(cookies) == 0 {
		return nil, fmt.Errorf("could not extract cookies - login may not have completed")
	}

	fmt.Println("Login successful! Saving cookies.")
	return cookies, nil
}

package auth

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha1"
	"database/sql"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
	"unicode"

	"golang.org/x/crypto/pbkdf2"
	_ "modernc.org/sqlite"
)

// ExtractChromeCookies reads Google cookies directly from Chrome's SQLite database.
// Does not launch Chrome, so it does not affect the existing login state.
// Works even while Chrome is running (read-only).
func ExtractChromeCookies() ([]*http.Cookie, error) {
	// 1. Get Chrome encryption key from macOS Keychain
	password, err := getChromeKeychainPassword()
	if err != nil {
		return nil, fmt.Errorf("cannot get Chrome password from Keychain: %w\n"+
			"Please verify Chrome is installed.", err)
	}

	// 2. Derive decryption key with PBKDF2
	key := pbkdf2.Key([]byte(password), []byte("saltysalt"), 1003, 16, sha1.New)

	// 3. Read cookies from Default profile
	homeDir, _ := os.UserHomeDir()
	cookiesDB := filepath.Join(homeDir, "Library", "Application Support", "Google", "Chrome", "Default", "Cookies")

	if _, err := os.Stat(cookiesDB); os.IsNotExist(err) {
		return nil, fmt.Errorf("Chrome cookie DB not found: %s", cookiesDB)
	}

	cookies, err := readCookiesFromDB(cookiesDB, key)
	if err != nil {
		return nil, fmt.Errorf("failed to read cookies: %w", err)
	}

	if len(cookies) == 0 {
		return nil, fmt.Errorf("no Google cookies found. Please verify you are logged in to Google in Chrome.")
	}

	// Verify essential cookies
	hasSID := false
	for _, c := range cookies {
		if c.Name == "SID" {
			hasSID = true
			break
		}
	}
	if !hasSID {
		return nil, fmt.Errorf("SID cookie not found. Please verify you are logged in to Google in Chrome.")
	}

	fmt.Printf("Extracted %d cookies.\n", len(cookies))
	return cookies, nil
}

func findChromeProfiles(chromeDir string) []string {
	var profiles []string

	// Default profile
	if _, err := os.Stat(filepath.Join(chromeDir, "Default")); err == nil {
		profiles = append(profiles, "Default")
	}

	// Profile 1, Profile 2, ... profiles
	entries, err := os.ReadDir(chromeDir)
	if err != nil {
		return profiles
	}
	for _, e := range entries {
		if e.IsDir() && strings.HasPrefix(e.Name(), "Profile ") {
			profiles = append(profiles, e.Name())
		}
	}

	return profiles
}

func readCookiesFromDB(cookiesDB string, key []byte) ([]*http.Cookie, error) {
	// Chrome may have the DB locked, so use a temporary copy
	tmpFile, err := copyToTemp(cookiesDB)
	if err != nil {
		return nil, err
	}
	defer os.Remove(tmpFile)

	// Copy WAL/SHM files as well
	for _, ext := range []string{"-wal", "-shm"} {
		src := cookiesDB + ext
		if _, err := os.Stat(src); err == nil {
			dst := tmpFile + ext
			copyToPath(src, dst)
			defer os.Remove(dst)
		}
	}

	db, err := sql.Open("sqlite", tmpFile)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	// Extract only cookies that are actually sent to notebooklm.google.com:
	//   .google.com           - common across all Google services (SID, HSID, etc.)
	//   .notebooklm.google.com - NotebookLM-specific (OSID, etc.)
	//   notebooklm.google.com  - NotebookLM-specific (host-only)
	// Note: exclude regional domain cookies like .google.co.kr (may conflict with other account cookies)
	rows, err := db.Query(`
		SELECT name, value, encrypted_value, host_key, path,
		       expires_utc, is_httponly, is_secure
		FROM cookies
		WHERE host_key = '.google.com'
		   OR host_key = '.notebooklm.google.com'
		   OR host_key = 'notebooklm.google.com'
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var cookies []*http.Cookie
	for rows.Next() {
		var name, plainValue string
		var encValue []byte
		var domain, path string
		var expiresUTC int64
		var httpOnly, secure int

		if err := rows.Scan(&name, &plainValue, &encValue, &domain, &path, &expiresUTC, &httpOnly, &secure); err != nil {
			continue
		}

		// Use plaintext if available, otherwise decrypt
		value := plainValue
		if value == "" && len(encValue) > 0 {
			decrypted, err := decryptCookieValue(encValue, key)
			if err != nil || !isPrintableASCII(decrypted) {
				continue // decryption failed, skip
			}
			value = decrypted
		}

		if value == "" {
			continue
		}

		cookie := &http.Cookie{
			Name:     name,
			Value:    value,
			Domain:   domain,
			Path:     path,
			HttpOnly: httpOnly != 0,
			Secure:   secure != 0,
		}

		if expiresUTC > 0 {
			const chromeEpochDelta = 11644473600000000
			unixMicro := expiresUTC - chromeEpochDelta
			cookie.Expires = time.UnixMicro(unixMicro)
		}

		cookies = append(cookies, cookie)
	}

	return cookies, nil
}

// isPrintableASCII checks that the string contains only valid cookie characters.
func isPrintableASCII(s string) bool {
	if s == "" {
		return false
	}
	for _, r := range s {
		if r > unicode.MaxASCII || !unicode.IsPrint(r) {
			return false
		}
	}
	return true
}

func getChromeKeychainPassword() (string, error) {
	out, err := exec.Command("security", "find-generic-password",
		"-s", "Chrome Safe Storage", "-w").Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(out)), nil
}

func decryptCookieValue(encrypted []byte, key []byte) (string, error) {
	if len(encrypted) == 0 {
		return "", nil
	}

	// "v10" prefix = Chrome macOS encryption (AES-128-CBC)
	if len(encrypted) >= 3 && string(encrypted[:3]) == "v10" {
		encrypted = encrypted[3:]
	} else {
		return string(encrypted), nil
	}

	if len(encrypted) < aes.BlockSize || len(encrypted)%aes.BlockSize != 0 {
		return "", fmt.Errorf("invalid encrypted data length")
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	// macOS Chrome IV: 16 bytes of space (0x20)
	iv := make([]byte, aes.BlockSize)
	for i := range iv {
		iv[i] = ' '
	}

	mode := cipher.NewCBCDecrypter(block, iv)
	decrypted := make([]byte, len(encrypted))
	mode.CryptBlocks(decrypted, encrypted)

	// Remove PKCS#7 padding
	if len(decrypted) > 0 {
		padLen := int(decrypted[len(decrypted)-1])
		if padLen > 0 && padLen <= aes.BlockSize && padLen <= len(decrypted) {
			valid := true
			for i := len(decrypted) - padLen; i < len(decrypted); i++ {
				if decrypted[i] != byte(padLen) {
					valid = false
					break
				}
			}
			if valid {
				decrypted = decrypted[:len(decrypted)-padLen]
			}
		}
	}

	// Chrome prepends a 32-byte salt to the encrypted plaintext.
	// After decryption, skip the first 32 bytes to get the actual cookie value.
	const saltLen = 32
	if len(decrypted) > saltLen {
		decrypted = decrypted[saltLen:]
	}

	return string(decrypted), nil
}

func copyToTemp(src string) (string, error) {
	tmpFile, err := os.CreateTemp("", "nlm-cookies-*.db")
	if err != nil {
		return "", err
	}
	tmpPath := tmpFile.Name()

	srcFile, err := os.Open(src)
	if err != nil {
		tmpFile.Close()
		os.Remove(tmpPath)
		return "", err
	}
	defer srcFile.Close()

	_, err = io.Copy(tmpFile, srcFile)
	tmpFile.Close()
	if err != nil {
		os.Remove(tmpPath)
		return "", err
	}

	return tmpPath, nil
}

func copyToPath(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, in)
	return err
}

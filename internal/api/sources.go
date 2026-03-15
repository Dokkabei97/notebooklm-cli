package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/Dokkabei97/notebooklm-cli/internal/auth"
	"github.com/Dokkabei97/notebooklm-cli/internal/model"
	"github.com/Dokkabei97/notebooklm-cli/internal/rpc"
)

// ListSources returns all sources in a notebook.
func (c *Client) ListSources(notebookID string) ([]model.Source, error) {
	params := []any{notebookID, nil, []any{2}, nil, 0}
	result, err := c.caller.Call(rpc.MethodGetNotebook, params, notebookPath(notebookID))
	if err != nil {
		return nil, wrapErr("ListSources", "rpc call failed", err)
	}
	if result == nil {
		return nil, nil // empty notebook
	}
	return parseSourcesFromNotebook(result.Parsed)
}

// AddSourceURL adds a URL source to a notebook.
func (c *Client) AddSourceURL(notebookID, sourceURL string) (*model.Source, error) {
	params := []any{
		[]any{[]any{nil, nil, []any{sourceURL}, nil, nil, nil, nil, nil}},
		notebookID,
		[]any{2},
		nil,
		nil,
	}
	result, err := c.caller.Call(rpc.MethodAddSource, params, notebookPath(notebookID))
	if err != nil {
		return nil, wrapErr("AddSourceURL", "rpc call failed", err)
	}
	if result == nil {
		return &model.Source{URL: sourceURL, Status: "processing"}, nil
	}
	return parseSource(result.Parsed)
}

// AddSourceText adds a text/paste source to a notebook.
func (c *Client) AddSourceText(notebookID, title, content string) (*model.Source, error) {
	// params: [[[None, [title, content], None, None, None, None, None, None]], notebook_id, [2], None, None]
	params := []any{
		[]any{[]any{nil, []any{title, content}, nil, nil, nil, nil, nil, nil}},
		notebookID,
		[]any{2},
		nil,
		nil,
	}
	result, err := c.caller.Call(rpc.MethodAddSource, params, notebookPath(notebookID))
	if err != nil {
		return nil, wrapErr("AddSourceText", "rpc call failed", err)
	}
	return parseSource(result.Parsed)
}

// AddSourceFile uploads a file as a source to a notebook (3-step resumable upload).
func (c *Client) AddSourceFile(notebookID, filePath string) (*model.Source, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, wrapErr("AddSourceFile", "read file", err)
	}

	filename := filepath.Base(filePath)

	// Step 1: Register file source intent and get SOURCE_ID
	regParams := []any{
		[]any{[]any{filename}},
		notebookID,
		[]any{2},
		[]any{1, nil, nil, nil, nil, nil, nil, nil, nil, nil, []any{1}},
	}
	regResult, err := c.caller.Call(rpc.MethodAddSourceFile, regParams, notebookPath(notebookID))
	if err != nil {
		return nil, wrapErr("AddSourceFile", "register file", err)
	}

	sourceID := extractNestedString(regResult.Parsed)
	if sourceID == "" {
		return nil, wrapErr("AddSourceFile", "no source ID from registration", nil)
	}

	// Step 2: Start resumable upload
	uploadURL, err := c.startResumableUpload(notebookID, filename, len(data), sourceID)
	if err != nil {
		return nil, wrapErr("AddSourceFile", "start upload", err)
	}

	// Step 3: Upload file content
	if err := c.uploadFileContent(uploadURL, data); err != nil {
		return nil, wrapErr("AddSourceFile", "upload content", err)
	}

	return &model.Source{
		ID:    sourceID,
		Title: filename,
	}, nil
}

func (c *Client) startResumableUpload(notebookID, filename string, fileSize int, sourceID string) (string, error) {
	uploadInitURL := rpc.UploadURL + "?authuser=0"

	body, err := json.Marshal(map[string]string{
		"PROJECT_ID":  notebookID,
		"SOURCE_NAME": filename,
		"SOURCE_ID":   sourceID,
	})
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", uploadInitURL, bytes.NewReader(body))
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded;charset=UTF-8")
	req.Header.Set("Origin", "https://notebooklm.google.com")
	req.Header.Set("Referer", "https://notebooklm.google.com/")
	req.Header.Set("x-goog-authuser", "0")
	req.Header.Set("x-goog-upload-command", "start")
	req.Header.Set("x-goog-upload-header-content-length", fmt.Sprintf("%d", fileSize))
	req.Header.Set("x-goog-upload-protocol", "resumable")

	req.Header.Set("Cookie", auth.BuildCookieHeader(c.tokens.Cookies))

	resp, err := c.caller.HTTPClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return "", fmt.Errorf("upload init failed with status %d", resp.StatusCode)
	}

	uploadURL := resp.Header.Get("x-goog-upload-url")
	if uploadURL == "" {
		return "", fmt.Errorf("no upload URL in response headers")
	}
	return uploadURL, nil
}

func (c *Client) uploadFileContent(uploadURL string, data []byte) error {
	req, err := http.NewRequest("POST", uploadURL, bytes.NewReader(data))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded;charset=utf-8")
	req.Header.Set("Origin", "https://notebooklm.google.com")
	req.Header.Set("Referer", "https://notebooklm.google.com/")
	req.Header.Set("x-goog-authuser", "0")
	req.Header.Set("x-goog-upload-command", "upload, finalize")
	req.Header.Set("x-goog-upload-offset", "0")

	req.Header.Set("Cookie", auth.BuildCookieHeader(c.tokens.Cookies))

	resp, err := c.caller.HTTPClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return fmt.Errorf("file upload failed with status %d", resp.StatusCode)
	}
	return nil
}

func extractNestedString(data any) string {
	switch v := data.(type) {
	case string:
		return v
	case []any:
		if len(v) > 0 {
			return extractNestedString(v[0])
		}
	}
	return ""
}

// DeleteSource removes a source from a notebook.
func (c *Client) DeleteSource(notebookID, sourceID string) error {
	// params: [[[source_id]]]
	params := []any{[]any{[]any{sourceID}}}
	_, err := c.caller.Call(rpc.MethodDeleteSource, params, notebookPath(notebookID))
	if err != nil {
		return wrapErr("DeleteSource", "rpc call failed", err)
	}
	return nil
}

// RefreshSource refreshes a source.
func (c *Client) RefreshSource(notebookID, sourceID string) error {
	params := []any{sourceID}
	_, err := c.caller.Call(rpc.MethodRefreshSource, params, notebookPath(notebookID))
	if err != nil {
		return wrapErr("RefreshSource", "rpc call failed", err)
	}
	return nil
}

// GetSource gets details about a specific source (via list + filter).
func (c *Client) GetSource(notebookID, sourceID string) (*model.Source, error) {
	sources, err := c.ListSources(notebookID)
	if err != nil {
		return nil, err
	}
	for _, src := range sources {
		if src.ID == sourceID {
			return &src, nil
		}
	}
	return nil, wrapErr("GetSource", "source not found: "+sourceID, nil)
}

// WaitForSource polls until a source is ready.
func (c *Client) WaitForSource(notebookID, sourceID string, timeout time.Duration) (*model.Source, error) {
	deadline := time.Now().Add(timeout)
	interval := time.Second
	for time.Now().Before(deadline) {
		src, err := c.GetSource(notebookID, sourceID)
		if err != nil {
			return nil, err
		}
		if src.Status == "ready" || src.Status == "error" {
			return src, nil
		}
		time.Sleep(interval)
		if interval < 10*time.Second {
			interval = time.Duration(float64(interval) * 1.5)
		}
	}
	return nil, wrapErr("WaitForSource", "timeout waiting for source to be ready", nil)
}

// parseSourcesFromNotebook extracts sources from GET_NOTEBOOK response.
func parseSourcesFromNotebook(data []any) ([]model.Source, error) {
	if data == nil || len(data) == 0 {
		return nil, nil
	}

	// Response: [[nb_info_with_sources, ...]]
	nbInfo := rpc.SafeArray(data, 0)
	if nbInfo == nil || len(nbInfo) <= 1 {
		return nil, nil
	}

	// Sources at nbInfo[1]
	sourcesList := rpc.SafeArray(nbInfo, 1)
	if sourcesList == nil {
		return nil, nil
	}

	var sources []model.Source
	for _, item := range sourcesList {
		arr, ok := item.([]any)
		if !ok || len(arr) == 0 {
			continue
		}

		src := model.Source{}

		// ID at src[0][0] or src[0]
		idData := arr[0]
		if idArr, ok := idData.([]any); ok && len(idArr) > 0 {
			src.ID = fmt.Sprintf("%v", idArr[0])
		} else {
			src.ID = fmt.Sprintf("%v", idData)
		}

		// Title at src[1]
		if len(arr) > 1 {
			if title, ok := arr[1].(string); ok {
				src.Title = title
			}
		}

		// URL at src[2][7][0]
		if len(arr) > 2 {
			if meta, ok := arr[2].([]any); ok && len(meta) > 7 {
				if urlList, ok := meta[7].([]any); ok && len(urlList) > 0 {
					if u, ok := urlList[0].(string); ok {
						src.URL = u
					}
				}
				// Type code at src[2][4]
				if len(meta) > 4 {
					if tc, ok := meta[4].(float64); ok {
						src.Type = model.SourceType(int(tc))
					}
				}
				// Timestamp at src[2][2][0]
				if len(meta) > 2 {
					if tsList, ok := meta[2].([]any); ok && len(tsList) > 0 {
						if ts, ok := tsList[0].(float64); ok && ts > 0 {
							src.CreatedAt = time.Unix(int64(ts), 0)
						}
					}
				}
			}
		}

		// Status at src[3][1]
		if len(arr) > 3 {
			if statusArr, ok := arr[3].([]any); ok && len(statusArr) > 1 {
				if sc, ok := statusArr[1].(float64); ok {
					src.Status = rpc.SourceStatusCode(int(sc)).String()
				}
			}
		}
		if src.Status == "" {
			src.Status = "ready"
		}

		sources = append(sources, src)
	}

	return sources, nil
}

func parseSource(data []any) (*model.Source, error) {
	if data == nil {
		return nil, wrapErr("parseSource", "nil response", nil)
	}

	src := &model.Source{
		ID: extractNestedString(data),
	}

	if len(data) > 1 {
		if title, ok := data[1].(string); ok {
			src.Title = title
		}
	}

	return src, nil
}

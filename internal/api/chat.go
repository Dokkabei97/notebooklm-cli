package api

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/jmk/notebooklm-cli/internal/auth"
	"github.com/jmk/notebooklm-cli/internal/model"
	"github.com/jmk/notebooklm-cli/internal/rpc"
)

// Ask sends a question to the notebook and returns the AI response.
func (c *Client) Ask(notebookID, question string, sourceIDs []string) (*model.AskResult, error) {
	if len(sourceIDs) == 0 {
		sources, err := c.ListSources(notebookID)
		if err == nil {
			for _, s := range sources {
				sourceIDs = append(sourceIDs, s.ID)
			}
		}
	}

	var sourcesArray []any
	for _, sid := range sourceIDs {
		sourcesArray = append(sourcesArray, []any{[]any{sid}})
	}

	conversationID := generateUUID()

	params := []any{
		sourcesArray,
		question,
		nil,
		[]any{2, nil, []any{1}, []any{1}},
		conversationID,
		nil,
		nil,
		notebookID,
		1,
	}

	paramsJSON, err := json.Marshal(params)
	if err != nil {
		return nil, wrapErr("Ask", "marshal params", err)
	}

	fReq := []any{nil, string(paramsJSON)}
	fReqJSON, err := json.Marshal(fReq)
	if err != nil {
		return nil, wrapErr("Ask", "marshal f.req", err)
	}

	body := "f.req=" + url.QueryEscape(string(fReqJSON))
	if c.tokens.CSRFToken != "" {
		body += "&at=" + url.QueryEscape(c.tokens.CSRFToken)
	}
	body += "&"

	chatURL := rpc.BuildChatURL(c.tokens.SessionID, 100000)

	req, err := http.NewRequest("POST", chatURL, strings.NewReader(body))
	if err != nil {
		return nil, wrapErr("Ask", "create request", err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Origin", "https://notebooklm.google.com")
	req.Header.Set("Referer", "https://notebooklm.google.com/")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/131.0.0.0 Safari/537.36")
	req.Header.Set("Cookie", auth.BuildCookieHeader(c.tokens.Cookies))

	resp, err := c.caller.HTTPClient.Do(req)
	if err != nil {
		return nil, wrapErr("Ask", "http request", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, wrapErr("Ask", "read response", err)
	}

	if resp.StatusCode != 200 {
		return nil, wrapErr("Ask", fmt.Sprintf("HTTP %d", resp.StatusCode), nil)
	}

	return parseStreamingResponse(string(respBody))
}

// GetChatHistory returns recent conversation turns for a notebook.
func (c *Client) GetChatHistory(notebookID string) ([]model.ChatEntry, error) {
	convResult, err := c.caller.Call(rpc.MethodGetLastConversationID, nil, notebookPath(notebookID))
	if err != nil {
		return nil, wrapErr("GetChatHistory", "get conversation ID", err)
	}

	if convResult == nil {
		return nil, nil
	}

	convID := rpc.SafeString(convResult.Parsed, 0)
	if convID == "" {
		return nil, nil
	}

	params := []any{convID, 20}
	result, err := c.caller.Call(rpc.MethodGetConversationTurns, params, notebookPath(notebookID))
	if err != nil {
		return nil, wrapErr("GetChatHistory", "get turns", err)
	}

	return parseChatHistory(result.Parsed)
}

// parseStreamingResponse parses the NotebookLM streaming chat response.
//
// Response structure (each chunk):
//   [["wrb.fr", null, "<inner_json_string>"]]
//
// inner_json structure:
//   [
//     [answer_text, null, [conv_id, turn_id, ts], null, [[[formatting]], ..., TYPE]],
//     citations,
//     ...,
//     follow_ups
//   ]
//
// TYPE: 1 = answer, 2 = thinking
func parseStreamingResponse(body string) (*model.AskResult, error) {
	body = strings.TrimSpace(body)
	if idx := strings.Index(body, "\n"); idx > 0 && strings.HasPrefix(body, ")]}'") {
		body = body[idx+1:]
	}

	var lastAnswer string
	var followUps []string

	lines := strings.Split(body, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		// Skip byte count lines
		isNum := true
		for _, c := range line {
			if c < '0' || c > '9' {
				isNum = false
				break
			}
		}
		if isNum {
			continue
		}

		// Parse JSON
		var outer []any
		if err := json.Unmarshal([]byte(line), &outer); err != nil {
			continue
		}

		// Find wrb.fr entries
		for _, item := range outer {
			arr, ok := item.([]any)
			if !ok || len(arr) < 3 {
				continue
			}
			marker, _ := arr[0].(string)
			if marker != "wrb.fr" {
				continue
			}

			innerStr, ok := arr[2].(string)
			if !ok {
				continue
			}

			var inner []any
			if err := json.Unmarshal([]byte(innerStr), &inner); err != nil {
				continue
			}

			// inner[0] = [answer_text, null, [conv], null, [formatting..., TYPE]]
			entry, ok := inner[0].([]any)
			if !ok || len(entry) < 5 {
				continue
			}

			text, _ := entry[0].(string)
			if text == "" {
				continue
			}

			// Extract TYPE from inner[0][4] (last element)
			formatting, ok := entry[4].([]any)
			if !ok || len(formatting) == 0 {
				continue
			}
			chunkType, _ := formatting[len(formatting)-1].(float64)

			if int(chunkType) == 2 {
				continue // skip thinking chunks
			}

			// TYPE=1 or other values = answer (keep the longest)
			if len(text) > len(lastAnswer) {
				lastAnswer = text
			}

			// Extract follow-up questions (from the last chunk)
			if len(inner) >= 4 {
				if fuArr, ok := inner[len(inner)-2].([]any); ok {
					for _, fuItem := range fuArr {
						if fuList, ok := fuItem.([]any); ok {
							for _, fu := range fuList {
								if s, ok := fu.(string); ok {
									followUps = append(followUps, s)
								}
							}
						}
					}
				}
			}
		}
	}

	if lastAnswer == "" {
		return &model.AskResult{Answer: "(no response received)"}, nil
	}

	return &model.AskResult{
		Answer:    lastAnswer,
		FollowUps: followUps,
	}, nil
}

func parseChatHistory(data []any) ([]model.ChatEntry, error) {
	if data == nil {
		return nil, nil
	}

	items := rpc.SafeArray(data, 0)
	var entries []model.ChatEntry
	for _, item := range items {
		arr, ok := item.([]any)
		if !ok || len(arr) < 4 {
			continue
		}

		turnType := rpc.SafeFloat(arr, 2)
		text := rpc.SafeString(arr, 3)

		role := "user"
		if int(turnType) == 2 {
			role = "assistant"
		}

		entries = append(entries, model.ChatEntry{
			Role:    role,
			Content: text,
		})
	}

	return entries, nil
}

func generateUUID() string {
	b := make([]byte, 16)
	rand.Read(b)
	b[6] = (b[6] & 0x0f) | 0x40
	b[8] = (b[8] & 0x3f) | 0x80
	return fmt.Sprintf("%08x-%04x-%04x-%04x-%012x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:])
}

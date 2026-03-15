package rpc

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
)

const antiXSSIPrefix = ")]}'"

// DecodeBatchResponse parses the batchexecute response format.
// Format:
//
//	)]}'        <- anti-XSSI prefix (remove)
//	<byte_count>
//	<json_payload>
//	... (repeated chunks)
//
// Each payload is a JSON array. We look for wrb.fr entries matching the target RPC ID.
// null result_data is valid (empty list, void operation success).
func DecodeBatchResponse(body string, targetRPCID string) (json.RawMessage, error) {
	body = strings.TrimSpace(body)
	if strings.HasPrefix(body, antiXSSIPrefix) {
		body = body[len(antiXSSIPrefix):]
	}
	body = strings.TrimSpace(body)

	chunks := parseChunks(body)
	if len(chunks) == 0 {
		return nil, fmt.Errorf("no chunks found in response")
	}

	for _, chunk := range chunks {
		result, found, err := extractResult(chunk, targetRPCID)
		if err != nil {
			return nil, err // RPC returned an error entry
		}
		if found {
			return result, nil // result can be nil (null result = empty/void)
		}
	}

	return nil, fmt.Errorf("no result found for rpc %s", targetRPCID)
}

// parseChunks splits the response into JSON chunks.
// Format: alternating byte-count lines and JSON payload lines.
// We skip byte-count lines and collect JSON-parseable lines.
func parseChunks(body string) []string {
	var chunks []string
	lines := strings.Split(body, "\n")

	i := 0
	for i < len(lines) {
		line := strings.TrimSpace(lines[i])
		i++

		if line == "" {
			continue
		}

		// byte count line → skip it, take the next line as JSON
		if _, err := strconv.Atoi(line); err == nil {
			if i < len(lines) {
				jsonLine := strings.TrimSpace(lines[i])
				i++
				if jsonLine != "" && jsonLine[0] == '[' {
					chunks = append(chunks, jsonLine)
				}
			}
			continue
		}

		// Try as direct JSON
		if line[0] == '[' {
			chunks = append(chunks, line)
		}
	}

	return chunks
}

// extractResult looks for a wrb.fr entry matching the target RPC ID.
// Returns (result, found, error):
//   - (data, true, nil)  → matched, result_data has content
//   - (nil,  true, nil)  → matched, result_data is null (empty/void)
//   - (nil, false, nil)  → not matched in this chunk
//   - (nil, false, err)  → RPC error entry found
func extractResult(chunk string, targetRPCID string) (json.RawMessage, bool, error) {
	var outer []json.RawMessage
	if err := json.Unmarshal([]byte(chunk), &outer); err != nil {
		return nil, false, nil
	}

	for _, entry := range outer {
		var arr []json.RawMessage
		if err := json.Unmarshal(entry, &arr); err != nil {
			continue
		}

		if len(arr) < 3 {
			continue
		}

		var marker string
		if err := json.Unmarshal(arr[0], &marker); err != nil {
			continue
		}

		// "er" = error entry
		if marker == "er" {
			var rpcID string
			if err := json.Unmarshal(arr[1], &rpcID); err != nil || rpcID != targetRPCID {
				continue
			}
			var errCode any
			if len(arr) > 2 {
				json.Unmarshal(arr[2], &errCode)
			}
			return nil, false, &Error{
				Code:    ErrServer,
				Message: fmt.Sprintf("server returned error: %v", errCode),
				Method:  targetRPCID,
			}
		}

		if marker != "wrb.fr" {
			continue
		}

		var rpcID string
		if err := json.Unmarshal(arr[1], &rpcID); err != nil {
			continue
		}

		if rpcID != targetRPCID {
			continue
		}

		// result_data is at arr[2]
		// It can be: a JSON string (needs double-parse), a JSON object/array, or null
		raw := arr[2]

		// Check for null
		if string(raw) == "null" {
			return nil, true, nil // null result = empty/void success
		}

		// Try to parse as JSON string (double-encoded)
		var resultStr string
		if err := json.Unmarshal(raw, &resultStr); err != nil {
			// Not a string — return as-is (already JSON object/array)
			return raw, true, nil
		}

		return json.RawMessage(resultStr), true, nil
	}

	return nil, false, nil
}

// ParseResultArray parses the result JSON into a generic nested array.
func ParseResultArray(raw json.RawMessage) ([]any, error) {
	if raw == nil {
		return nil, nil
	}
	var result []any
	if err := json.Unmarshal(raw, &result); err != nil {
		return nil, fmt.Errorf("parse result array: %w", err)
	}
	return result, nil
}

// SafeString extracts a string from a nested array at the given indices.
func SafeString(data []any, indices ...int) string {
	current := any(data)
	for _, idx := range indices {
		arr, ok := current.([]any)
		if !ok || idx >= len(arr) || arr[idx] == nil {
			return ""
		}
		current = arr[idx]
	}
	s, ok := current.(string)
	if !ok {
		return fmt.Sprintf("%v", current)
	}
	return s
}

// SafeFloat extracts a float64 from a nested array at the given indices.
func SafeFloat(data []any, indices ...int) float64 {
	current := any(data)
	for _, idx := range indices {
		arr, ok := current.([]any)
		if !ok || idx >= len(arr) || arr[idx] == nil {
			return 0
		}
		current = arr[idx]
	}
	f, ok := current.(float64)
	if !ok {
		return 0
	}
	return f
}

// SafeArray extracts a sub-array from a nested array at the given indices.
func SafeArray(data []any, indices ...int) []any {
	current := any(data)
	for _, idx := range indices {
		arr, ok := current.([]any)
		if !ok || idx >= len(arr) || arr[idx] == nil {
			return nil
		}
		current = arr[idx]
	}
	arr, ok := current.([]any)
	if !ok {
		return nil
	}
	return arr
}

package api

import (
	"time"

	"github.com/Dokkabei97/notebooklm-cli/internal/model"
	"github.com/Dokkabei97/notebooklm-cli/internal/rpc"
)

// ListArtifacts returns all artifacts in a notebook.
func (c *Client) ListArtifacts(notebookID string) ([]model.Artifact, error) {
	result, err := c.caller.Call(rpc.MethodListArtifacts, nil, notebookPath(notebookID))
	if err != nil {
		return nil, wrapErr("ListArtifacts", "rpc call failed", err)
	}
	if result == nil {
		return nil, nil
	}
	return parseArtifacts(result.Parsed)
}

// CreateArtifact creates a new artifact (report, quiz, etc.).
func (c *Client) CreateArtifact(notebookID string, typeCode rpc.ArtifactTypeCode) (*model.Artifact, error) {
	params := []any{nil, int(typeCode)}
	result, err := c.caller.Call(rpc.MethodCreateArtifact, params, notebookPath(notebookID))
	if err != nil {
		return nil, wrapErr("CreateArtifact", "rpc call failed", err)
	}
	return parseArtifact(result.Parsed)
}

// GenerateAudio generates an audio overview for the notebook.
func (c *Client) GenerateAudio(notebookID string, instructions string) (*model.Artifact, error) {
	params := []any{nil, int(rpc.ArtifactCodeAudio)}
	if instructions != "" {
		params = append(params, instructions)
	}
	result, err := c.caller.Call(rpc.MethodCreateArtifact, params, notebookPath(notebookID))
	if err != nil {
		return nil, wrapErr("GenerateAudio", "rpc call failed", err)
	}
	return parseArtifact(result.Parsed)
}

// DeleteArtifact deletes an artifact.
func (c *Client) DeleteArtifact(notebookID, artifactID string) error {
	params := []any{[]any{artifactID}}
	_, err := c.caller.Call(rpc.MethodDeleteArtifact, params, notebookPath(notebookID))
	if err != nil {
		return wrapErr("DeleteArtifact", "rpc call failed", err)
	}
	return nil
}

// GetArtifact returns artifact details including content.
func (c *Client) GetArtifact(notebookID, artifactID string) (*model.Artifact, error) {
	// List artifacts and find by ID
	artifacts, err := c.ListArtifacts(notebookID)
	if err != nil {
		return nil, err
	}
	for _, art := range artifacts {
		if art.ID == artifactID {
			return &art, nil
		}
	}
	return nil, wrapErr("GetArtifact", "artifact not found: "+artifactID, nil)
}

// WaitForArtifact polls until an artifact is ready.
func (c *Client) WaitForArtifact(notebookID, artifactID string, timeout time.Duration) (*model.Artifact, error) {
	deadline := time.Now().Add(timeout)
	interval := 3 * time.Second
	for time.Now().Before(deadline) {
		art, err := c.GetArtifact(notebookID, artifactID)
		if err != nil {
			return nil, err
		}
		if art.Status == "completed" || art.Status == "failed" {
			return art, nil
		}
		time.Sleep(interval)
	}
	return nil, wrapErr("WaitForArtifact", "timeout waiting for artifact", nil)
}

func parseArtifacts(data []any) ([]model.Artifact, error) {
	if data == nil {
		return nil, nil
	}

	items := rpc.SafeArray(data, 0)
	if items == nil {
		// Try direct array
		items = data
	}

	var artifacts []model.Artifact
	for _, item := range items {
		arr, ok := item.([]any)
		if !ok {
			continue
		}
		art, err := parseArtifactFromArray(arr)
		if err != nil {
			continue
		}
		artifacts = append(artifacts, *art)
	}

	return artifacts, nil
}

func parseArtifact(data []any) (*model.Artifact, error) {
	if data == nil {
		return nil, wrapErr("parseArtifact", "nil response", nil)
	}
	return parseArtifactFromArray(data)
}

func parseArtifactFromArray(arr []any) (*model.Artifact, error) {
	art := &model.Artifact{
		ID:    rpc.SafeString(arr, 0),
		Title: rpc.SafeString(arr, 1),
	}

	// artifact_data[2] = type code
	if typeVal := rpc.SafeFloat(arr, 2); typeVal > 0 {
		art.Type = model.ArtifactType(int(typeVal))
	}

	// artifact_data[3] = content
	art.Content = rpc.SafeString(arr, 3)

	// artifact_data[4] = status code
	if statusVal := rpc.SafeFloat(arr, 4); statusVal > 0 {
		art.Status = rpc.ArtifactStatusCode(int(statusVal)).String()
	}
	if art.Status == "" {
		art.Status = "unknown"
	}

	// Audio URL if present
	art.AudioURL = rpc.SafeString(arr, 5)

	// Timestamp
	if ts := rpc.SafeFloat(arr, 6); ts > 0 {
		art.CreatedAt = time.UnixMicro(int64(ts))
	}

	return art, nil
}

package api

import (
	"time"

	"github.com/Dokkabei97/notebooklm-cli/internal/model"
	"github.com/Dokkabei97/notebooklm-cli/internal/rpc"
)

// ListNotes returns all notes in a notebook.
func (c *Client) ListNotes(notebookID string) ([]model.Note, error) {
	result, err := c.caller.Call(rpc.MethodGetNotesAndMindMaps, nil, notebookPath(notebookID))
	if err != nil {
		return nil, wrapErr("ListNotes", "rpc call failed", err)
	}
	if result == nil {
		return nil, nil
	}
	return parseNotes(result.Parsed)
}

// CreateNote creates a new note in the notebook.
func (c *Client) CreateNote(notebookID, title, content string) (*model.Note, error) {
	params := []any{nil, title, content}
	result, err := c.caller.Call(rpc.MethodCreateNote, params, notebookPath(notebookID))
	if err != nil {
		return nil, wrapErr("CreateNote", "rpc call failed", err)
	}
	if result == nil {
		return &model.Note{Title: title, Content: content}, nil
	}
	return parseNote(result.Parsed)
}

// UpdateNote updates an existing note.
func (c *Client) UpdateNote(notebookID, noteID, title, content string) (*model.Note, error) {
	params := []any{noteID, title, content}
	result, err := c.caller.Call(rpc.MethodUpdateNote, params, notebookPath(notebookID))
	if err != nil {
		return nil, wrapErr("UpdateNote", "rpc call failed", err)
	}
	if result == nil {
		return &model.Note{ID: noteID, Title: title, Content: content}, nil
	}
	return parseNote(result.Parsed)
}

// DeleteNote deletes a note.
func (c *Client) DeleteNote(notebookID, noteID string) error {
	params := []any{[]any{noteID}}
	_, err := c.caller.Call(rpc.MethodDeleteNote, params, notebookPath(notebookID))
	if err != nil {
		return wrapErr("DeleteNote", "rpc call failed", err)
	}
	return nil
}

func parseNotes(data []any) ([]model.Note, error) {
	if data == nil {
		return nil, nil
	}

	items := rpc.SafeArray(data, 0)
	if items == nil {
		return nil, nil
	}

	var notes []model.Note
	for _, item := range items {
		arr, ok := item.([]any)
		if !ok {
			continue
		}
		note, err := parseNoteFromArray(arr)
		if err != nil {
			continue
		}
		notes = append(notes, *note)
	}

	return notes, nil
}

func parseNote(data []any) (*model.Note, error) {
	if data == nil {
		return nil, wrapErr("parseNote", "nil response", nil)
	}
	return parseNoteFromArray(data)
}

func parseNoteFromArray(arr []any) (*model.Note, error) {
	note := &model.Note{
		ID:      rpc.SafeString(arr, 0),
		Title:   rpc.SafeString(arr, 1),
		Content: rpc.SafeString(arr, 2),
	}

	if ts := rpc.SafeFloat(arr, 3); ts > 0 {
		note.CreatedAt = time.UnixMicro(int64(ts))
	}
	if ts := rpc.SafeFloat(arr, 4); ts > 0 {
		note.UpdatedAt = time.UnixMicro(int64(ts))
	}

	return note, nil
}

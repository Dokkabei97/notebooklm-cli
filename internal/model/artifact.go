package model

import "time"

type Artifact struct {
	ID        string       `json:"id"`
	Title     string       `json:"title"`
	Type      ArtifactType `json:"type"`
	Status    string       `json:"status"`
	Content   string       `json:"content,omitempty"`
	AudioURL  string       `json:"audio_url,omitempty"`
	CreatedAt time.Time    `json:"created_at"`
}

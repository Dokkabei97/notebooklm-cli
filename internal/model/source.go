package model

import "time"

type Source struct {
	ID        string     `json:"id"`
	Title     string     `json:"title"`
	Type      SourceType `json:"type"`
	URL       string     `json:"url,omitempty"`
	Status    string     `json:"status"`
	CreatedAt time.Time  `json:"created_at"`
}

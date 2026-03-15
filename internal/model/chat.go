package model

type AskResult struct {
	Answer    string       `json:"answer"`
	Sources   []CitedChunk `json:"sources,omitempty"`
	FollowUps []string     `json:"follow_ups,omitempty"`
}

type CitedChunk struct {
	SourceID   string `json:"source_id"`
	SourceName string `json:"source_name"`
	Text       string `json:"text"`
}

type ChatEntry struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

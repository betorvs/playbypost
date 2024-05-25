package types

type Story struct {
	ID            int    `json:"id,omitempty"`
	Title         string `json:"title"`
	Announcement  string `json:"announcement"`
	Notes         string `json:"notes"`
	StorytellerID int    `json:"storyteller_id"`
}

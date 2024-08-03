package types

type Encounter struct {
	ID           int    `json:"id,omitempty"`
	Title        string `json:"title"`
	Announcement string `json:"announcement"`
	Notes        string `json:"notes"`
	StoryID      int    `json:"story_id"`
	WriterID     int    `json:"writer_id"`
}

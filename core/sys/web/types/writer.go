package types

/*
Writer

EncodingKeys story_id -> Encoding key
*/
type Writer struct {
	ID           int            `json:"id,omitempty"`
	Username     string         `json:"username"`
	Password     string         `json:"password,omitempty"`
	EncodingKeys map[int]string `json:"encoding_keys,omitempty"`
}

type Card struct {
	Username string   `json:"username"`
	UserID   string   `json:"user_id"`
	Stories  []string `json:"stories,omitempty"`
	Players  []string `json:"players,omitempty"`
}

type WriterStatus struct {
	UserID int64  `json:"user_id"`
	Status string `json:"status"`
}

type WriterUserAssociation struct {
	ID       int `json:"id,omitempty"`
	WriterID int `json:"writer_id"`
	UserID   int `json:"user_id"`
}

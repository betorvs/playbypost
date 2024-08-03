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

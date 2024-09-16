package types

const (
	HeaderUserID       string = "X-UserID"
	HeaderStory        string = "X-Story"
	HeaderPlayerID     string = "X-Player-ID"
	HeaderStoryChannel string = "X-Story-Channel"
)

type Command struct {
	// UserID string `json:"id,omitempty"`
	// Story  string `json:"story"`
	Text string `json:"text"`
}

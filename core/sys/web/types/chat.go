package types

const (
	Slack   string = "slack"
	Discord string = "discord"
)

type ChatInfo struct {
	ID       int    `json:"id,omitempty"`
	UserID   string `json:"user_id"`
	Username string `json:"username"`
	Channel  string `json:"channel"`
	Chat     string `json:"chat"`
}

type RunningChannels struct {
	Title   string `json:"title"`
	Channel string `json:"channel"`
	Kind    string `json:"kind"`
}

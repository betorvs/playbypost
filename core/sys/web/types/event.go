package types

const (
	EventAnnounce    = "announce"
	EventSuccess     = "success"
	EventFailure     = "failure"
	EventDead        = "dead"
	EventInformation = "information"
	EventEnd         = "end"
)

type Event struct {
	Channel  string `json:"channel"`
	UserID   string `json:"user_id"`
	Message  string `json:"message"`
	ImageURL string `json:"image_url"`
	Kind     string `json:"kind"`
}

type Options struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Value string `json:"value"`
}

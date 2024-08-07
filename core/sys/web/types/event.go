package types

const (
	EventAnnounce = "announce"
	EventSuccess  = "success"
	EventFailure  = "failure"
	EventDead     = "dead"
)

type Event struct {
	Channel string `json:"channel"`
	UserID  string `json:"user_id"`
	Message string `json:"message"`
	Kind    string `json:"kind"`
}

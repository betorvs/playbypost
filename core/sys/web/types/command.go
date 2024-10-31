package types

const (
	HeaderUserID       string = "X-UserID"
	HeaderStory        string = "X-Story"
	HeaderPlayerID     string = "X-Player-ID"
	HeaderStoryChannel string = "X-Story-Channel"

	// shared constants
	Solo     string = "solo"
	Choice   string = "choice"
	Decision string = "decision"
	Didatic  string = "didatic"

	// commands constants
	Cmd          string = "cmd"
	DidaticJoin  string = "didatic-join"
	DidaticNext  string = "didatic-next"
	DidaticStart string = "didatic-start"
	SoloNext     string = "solo-next"
	SoloStart    string = "solo-start"
	Opt          string = "opt"
)

type Command struct {
	// UserID string `json:"id,omitempty"`
	// Story  string `json:"story"`
	Text string `json:"text"`
}

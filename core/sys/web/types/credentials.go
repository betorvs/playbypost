package types

import "time"

const (
	HeaderToken    string = "X-Access-Token"
	HeaderUsername string = "X-Username"
)

type Credentials struct {
	Password string `json:"password"`
	Username string `json:"username"`
}

// each session contains the username of the user and the time at which it expires
type Session struct {
	Username     string
	Token        string
	Expiry       time.Time
	UserID       int
	EncodingKey  string
	ClientType   string
	ClientInfo   string
	IPAddress    string
	UserAgent    string
	CreatedAt    time.Time
	UpdatedAt    time.Time
	LastActivity time.Time
}

// we'll use this method later to determine if the session has expired
func (s Session) IsExpired() bool {
	return s.Expiry.Before(time.Now())
}

type Login struct {
	Status      string    `json:"status"`
	Message     string    `json:"message"`
	AccessToken string    `json:"access_token"`
	ExpireOn    time.Time `json:"expire_on"`
	UserID      int       `json:"user_id"`
}

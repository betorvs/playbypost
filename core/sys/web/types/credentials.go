package types

import (
	"time"
)

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
	ID           int64
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

type SessionEvent struct {
	ID        int       `json:"id"`
	SessionID int64     `json:"session_id"`
	EventType string    `json:"event_type"`
	Timestamp time.Time `json:"timestamp"`
	Data      string    `json:"data"`
}

type SessionEventRequest struct {
	SessionID int64  `json:"session_id"`
	EventType string `json:"event_type"`
	Data      string `json:"data"`
}

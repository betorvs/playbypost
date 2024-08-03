package types

type User struct {
	ID     int    `json:"id"`
	UserID string `json:"user_id"`
	Active bool   `json:"active"`
}

package usersservice

import (
	"encoding/json"
)

// User represents a user of the system.
// Owner of purchase orders and sell orders.
// Only the profile portion of a user is ever stored
// server side
type User interface {
	// JSON serializes the user record
	JSON() ([]byte, error)
}

// Profile represents the information about a user
// that the application stores.
type Profile interface {
	// JSON serializes the Profile record
	JSON() ([]byte, error)
}

type _User struct {
	IDToken     string   `json:"id_token"`
	AccessToken string   `json:"access_token"`
	Profile     _Profile `json:"profile"`
}

func (u *_User) JSON() ([]byte, error) {
	return json.Marshal(u)
}

type _Profile struct {
	UserID string  `bson:"user_id" json:"user_id"`
	Name   string  `bson:"name" json:"name"`
	Riches float32 `bson:"riches" json:"riches"`
}

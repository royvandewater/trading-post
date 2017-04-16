package usersservice

import (
	"encoding/json"
	"io"
	"io/ioutil"
)

// User represents a user of the system.
// Owner of purchase orders and sell orders.
type User interface {
	JSON() ([]byte, error)
}

// ParseUser creates a User instance from a ReadCloser
// (like an HTTP request body)
func ParseUser(data io.ReadCloser) (User, error) {
	dataBytes, err := ioutil.ReadAll(data)
	if err != nil {
		return nil, err
	}

	user := &user{}
	err = json.Unmarshal(dataBytes, user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

type user struct {
	Name string `json:"name"`
}

func (u *user) JSON() ([]byte, error) {
	return json.Marshal(u)
}

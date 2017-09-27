package usersservice

import (
	"encoding/json"
)

// User represents a user of the system.
// Owner of purchase orders and sell orders.
// Only the profile portion of a user is ever stored
// server side
type User interface {
	// GetAccessToken returns the user's access token
	GetAccessToken() string

	// GetIDToken returns the user's ID token
	GetIDToken() string

	// GetProfile returns the user profile
	GetProfile() Profile

	// GetRefreshToken returns the user's refresh token
	GetRefreshToken() string

	// JSON serializes the user record
	JSON() ([]byte, error)
}

// Profile represents the information about a user
// that the application stores.
type Profile interface {
	// GetName returns the profile's name
	GetName() string

	// GetRiches returns the profile's riches, counted in 1/1000th of a dollars
	// i.e, a value of 1000 would be $1
	GetRiches() int

	// GetStocks returns the profile's stocks
	GetStocks() []Stock

	// GetUserID returns the profile's ID in the system
	GetUserID() string
}

// Stock represents the quantity of stock
// owned by ticker
type Stock interface {
	// GetQuantity returns the quantity of the stock
	GetQuantity() int

	// GetTicker returns the ticker of the stock
	GetTicker() string
}

type _User struct {
	IDToken      string    `json:"id_token"`
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token"`
	Profile      *_Profile `json:"profile"`
}

func (u *_User) GetAccessToken() string {
	return u.AccessToken
}

func (u *_User) GetIDToken() string {
	return u.IDToken
}

func (u *_User) GetProfile() Profile {
	return u.Profile
}

func (u *_User) GetRefreshToken() string {
	return u.RefreshToken
}

func (u *_User) JSON() ([]byte, error) {
	return json.Marshal(u)
}

type _Profile struct {
	UserID string    `bson:"user_id" json:"sub"`
	Name   string    `bson:"name" json:"name"`
	Riches int       `bson:"riches,omitempty" json:"riches,omitempty"`
	Stocks []*_Stock `bson:"stocks"`
}

func (p *_Profile) GetName() string {
	return p.Name
}

func (p *_Profile) GetRiches() int {
	return p.Riches
}

func (p *_Profile) GetStocks() []Stock {
	stocks := make([]Stock, len(p.Stocks))

	for i, _stock := range p.Stocks {
		stocks[i] = _stock
	}

	return stocks
}

func (p *_Profile) GetUserID() string {
	return p.UserID
}

type _Stock struct {
	Quantity int    `bson:"quantity"`
	Ticker   string `bson:"ticker"`
}

func (s *_Stock) GetQuantity() int {
	return s.Quantity
}

func (s *_Stock) GetTicker() string {
	return s.Ticker
}

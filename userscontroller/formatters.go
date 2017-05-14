package userscontroller

import (
	"encoding/json"

	"github.com/royvandewater/trading-post/usersservice"
)

func formatUser(user usersservice.User) ([]byte, error) {
	profile := user.GetProfile()

	_stocks := make([]*_Stock, len(profile.GetStocks()))
	for i, stock := range profile.GetStocks() {
		_stocks[i] = &_Stock{
			Quantity: stock.GetQuantity(),
			Ticker:   stock.GetTicker(),
		}
	}

	_profile := &_Profile{
		Name:   profile.GetName(),
		Riches: float64(profile.GetRiches()) / 1000,
		Stocks: _stocks,
	}

	return json.MarshalIndent(_User{
		IDToken:     user.GetIDToken(),
		AccessToken: user.GetAccessToken(),
		Profile:     _profile,
	}, "", "  ")
}

type _User struct {
	IDToken     string    `json:"id_token"`
	AccessToken string    `json:"access_token"`
	Profile     *_Profile `json:"profile"`
}

type _Profile struct {
	UserID string    `json:"user_id"`
	Name   string    `json:"name"`
	Riches float64   `json:"riches"`
	Stocks []*_Stock `json:"stocks"`
}

type _Stock struct {
	Quantity int    `json:"quantity"`
	Ticker   string `json:"ticker"`
}

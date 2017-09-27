package profilescontroller

import (
	"encoding/json"

	"github.com/royvandewater/trading-post/usersservice"
)

func formatGetResponse(profile usersservice.Profile) ([]byte, error) {
	_stocks := make([]_Stock, len(profile.GetStocks()))

	for i, stock := range profile.GetStocks() {
		_stocks[i] = _Stock{
			Quantity: stock.GetQuantity(),
			Ticker:   stock.GetTicker(),
		}
	}

	return json.MarshalIndent(_GetProfile{
		Name:   profile.GetName(),
		Riches: float64(profile.GetRiches()) / 1000,
		Stocks: _stocks,
		UserID: profile.GetUserID(),
	}, "", "  ")
}

type _GetProfile struct {
	Name   string   `json:"name"`
	Riches float64  `json:"riches"`
	Stocks []_Stock `json:"stocks"`
	UserID string   `json:"user_id"`
}

type _Stock struct {
	Quantity int    `json:"quantity"`
	Ticker   string `json:"ticker"`
}

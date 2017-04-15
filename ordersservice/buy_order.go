package ordersservice

import (
	"encoding/json"
	"io"
	"io/ioutil"
)

// BuyOrder represents a purchase order for a
// specific stock
type BuyOrder interface {
	JSON() ([]byte, error)
}

// ParseBuyOrder instantiates a new BuyOrder instance from JSON
// data
func ParseBuyOrder(data io.ReadCloser) (BuyOrder, error) {
	dataBytes, err := ioutil.ReadAll(data)
	if err != nil {
		return nil, err
	}

	buyOrder := &buyOrder{}
	err = json.Unmarshal(dataBytes, buyOrder)
	if err != nil {
		return nil, err
	}

	return buyOrder, nil
}

type buyOrder struct {
	Ticker string `json:"ticker"`
}

func (order *buyOrder) JSON() ([]byte, error) {
	return make([]byte, 0), nil
}

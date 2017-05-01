package sellorderscontroller

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"strings"
)

func parseCreateBody(body io.ReadCloser) (*_CreateBody, error) {
	bodyBytes, err := ioutil.ReadAll(body)
	if err != nil {
		return nil, err
	}

	createBody := &_CreateBody{}
	err = json.Unmarshal(bodyBytes, createBody)
	if err != nil {
		return nil, err
	}

	if createBody.Quantity < 1 {
		return nil, fmt.Errorf("quantity must be greater than 0, was %v", createBody.Quantity)
	}

	createBody.Ticker = strings.ToLower(createBody.Ticker)
	return createBody, nil
}

type _CreateBody struct {
	Quantity int    `json:"quantity"`
	Ticker   string `json:"ticker"`
}

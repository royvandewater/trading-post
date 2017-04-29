package sellorderscontroller

import (
	"encoding/json"
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

	createBody.Ticker = strings.ToLower(createBody.Ticker)
	return createBody, nil
}

type _CreateBody struct {
	Ticker string `json:"ticker"`
}

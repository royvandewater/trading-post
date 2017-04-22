package sellorderscontroller

import (
	"encoding/json"
	"io"
	"io/ioutil"
)

func parseCreateBody(body io.ReadCloser) (_CreateBody, error) {
	createBody := _CreateBody{}

	bodyBytes, err := ioutil.ReadAll(body)
	if err != nil {
		return createBody, err
	}

	err = json.Unmarshal(bodyBytes, &createBody)
	if err != nil {
		return createBody, err
	}

	return createBody, nil
}

type _CreateBody struct {
	Ticker string `json:"ticker"`
}

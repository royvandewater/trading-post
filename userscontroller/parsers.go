package userscontroller

import (
	"encoding/json"
	"io"
	"io/ioutil"
)

func parseTokenBody(body io.ReadCloser) (*_TokenBody, error) {
	bodyBytes, err := ioutil.ReadAll(body)
	if err != nil {
		return nil, err
	}

	tokenBody := &_TokenBody{}
	err = json.Unmarshal(bodyBytes, tokenBody)
	if err != nil {
		return nil, err
	}

	return tokenBody, nil
}

type _TokenBody struct {
	GrantType    string `json:"grant_type"`
	RefreshToken string `json:"refresh_token"`
}

package oauthrefresh

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
)

// Refresh exchanges a refresh token for a new access token
func Refresh(refreshToken, clientID, clientSecret, tokenURL string) (*Token, error) {
	payload, err := json.MarshalIndent(_Payload{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		GrantType:    "refresh_token",
		RefreshToken: refreshToken,
	}, "", " ")

	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", tokenURL, bytes.NewReader(payload))
	if err != nil {
		return nil, err
	}

	req.Header.Add("content-type", "application/json")
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()
	bodyBytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	token := &Token{}
	err = json.Unmarshal(bodyBytes, token)
	if err != nil {
		return nil, err
	}

	return token, nil
}

type _Payload struct {
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	GrantType    string `json:"grant_type"`
	RefreshToken string `json:"refresh_token"`
}

// Token represents the response of a successful
// refresh token request.
type Token struct {
	AccessToken  string `json:"access_token"`
	ExpiresIn    int    `json:"expires_in"`
	IDToken      string `json:"id_token"`
	RefreshToken string `json:"refresh_token"`
	Scope        string `json:"scope"`
	TokenType    string `json:"token_type"`
}

package userscontroller

import "fmt"

func validateTokenBody(tokenBody *_TokenBody) error {
	if tokenBody.GrantType != "refresh_token" {
		return fmt.Errorf("Parameter \"grant_type\" must be \"refresh_token\", was \"%v\"", tokenBody.GrantType)
	}
	return nil
}

package profilescontroller

import (
	"encoding/json"

	"github.com/royvandewater/trading-post/usersservice"
)

func formatGetResponse(profile usersservice.Profile) ([]byte, error) {
	return json.MarshalIndent(struct {
		Name   string  `json:"name"`
		Riches float32 `json:"riches"`
	}{
		Name:   profile.GetName(),
		Riches: profile.GetRiches(),
	}, "", "  ")
}

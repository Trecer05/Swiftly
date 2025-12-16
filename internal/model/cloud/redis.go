package cloud

import (
	"encoding/json"
)

type Types string

type Envelope struct {
	TeamID int `json:"team_id"`
	Data   json.RawMessage `json:"data"`
	Type   Types `json:"type"`
}
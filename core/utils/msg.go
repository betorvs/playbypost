package utils

import (
	"encoding/json"

	"github.com/betorvs/playbypost/core/sys/web/types"
)

func ParseMsgBody(body []byte) (types.Msg, error) {
	var msg types.Msg
	err := json.Unmarshal(body, &msg)
	return msg, err
}

func ParseComposedBody(body []byte) (types.Composed, error) {
	var msg types.Composed
	err := json.Unmarshal(body, &msg)
	return msg, err
}

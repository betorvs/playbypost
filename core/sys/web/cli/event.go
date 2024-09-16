package cli

import (
	"encoding/json"

	"github.com/betorvs/playbypost/core/sys/web/types"
)

const (
	event string = "event"
)

func (c *Cli) PostEvent(channel, userid, message, kind string) ([]byte, error) {
	if kind == "" {
		kind = types.EventAnnounce
	}
	u := types.Event{
		Channel: channel,
		UserID:  userid,
		Message: message,
		Kind:    kind,
	}
	body, err := json.Marshal(u)
	if err != nil {
		return []byte{}, err
	}
	res, err := c.postGeneric(event, body)
	return res, err
}

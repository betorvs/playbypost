package cli

import (
	"encoding/json"

	"github.com/betorvs/playbypost/core/sys/web/types"
)

const (
	slack string = "info"
)

func (c *Cli) AddSlackInformation(userid, username, channel string) ([]byte, error) {
	u := types.SlackInfo{
		Username: username,
		UserID:   userid,
		Channel:  channel,
	}
	body, err := json.Marshal(u)
	if err != nil {
		return []byte{}, err
	}
	res, err := c.postGeneric(slack, body)
	return res, err
}

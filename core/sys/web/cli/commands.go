package cli

import (
	"encoding/json"

	"github.com/betorvs/playbypost/core/sys/web/types"
)

const (
	command string = "command"
)

func (c *Cli) PostCommand(userid, text, channel string) ([]byte, error) {
	cmd := types.Command{
		Text: text,
	}
	body, err := json.Marshal(cmd)
	if err != nil {
		return []byte{}, err
	}
	headers := makeHeaders("", "", userid, "", channel)
	res, err := c.postGenericWithHeaders(command, body, headers)
	if err != nil {
		return res, err
	}
	return res, nil
}

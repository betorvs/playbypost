package cli

import (
	"encoding/json"
	"fmt"

	"github.com/betorvs/playbypost/core/sys/web/types"
)

const (
	command string = "command"
)

func (c *Cli) postCommand(userid, text, channel string) ([]byte, error) {
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

func (c *Cli) PostCommandComposed(userid, text, channel string) (types.Composed, error) {
	var msg types.Composed
	body, err := c.postCommand(userid, text, channel)
	if err != nil {
		return msg, fmt.Errorf("PostCommand call error: %v", err)
	}
	err = json.Unmarshal(body, &msg)
	if err != nil {
		return msg, fmt.Errorf("json unmarshal error: %v", err)
	}
	return msg, nil
}

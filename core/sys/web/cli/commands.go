package cli

import (
	"encoding/json"

	"github.com/betorvs/playbypost/core/sys/web/types"
)

const (
	command string = "commands"
)

func (c *Cli) PostCommand(userid, story, text, storyChannel string, playerID int) ([]byte, error) {
	cmd := types.Command{
		Text: text,
	}
	body, err := json.Marshal(cmd)
	if err != nil {
		return []byte{}, err
	}
	headers := makeHeaders("", "", userid, story, storyChannel, playerID)
	res, err := c.postGenericWithHeaders(command, body, headers)
	if err != nil {
		return res, err
	}
	return res, nil
}

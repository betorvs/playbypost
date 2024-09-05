package cli

import (
	"encoding/json"
	"fmt"

	"github.com/betorvs/playbypost/core/sys/web/types"
)

const (
	info string = "info"
)

func (c *Cli) AddChatInformation(userid, username, channel, chat string) ([]byte, error) {
	u := types.ChatInfo{
		Username: username,
		UserID:   userid,
		Channel:  channel,
		Chat:     chat,
	}
	body, err := json.Marshal(u)
	if err != nil {
		return []byte{}, err
	}
	fmt.Println(chat)
	res, err := c.postGeneric(info, body)
	return res, err
}

func (c *Cli) GetChatInformation() ([]types.ChatInfo, error) {
	users := fmt.Sprintf("%s/%s", info, "users")
	body, err := c.getGeneric(users)
	if err != nil {
		return []types.ChatInfo{}, err
	}
	var u []types.ChatInfo
	err = json.Unmarshal(body, &u)
	if err != nil {
		return []types.ChatInfo{}, err
	}
	return u, nil
}

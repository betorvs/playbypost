package cli

import (
	"encoding/json"

	"github.com/betorvs/playbypost/core/sys/web/types"
)

const (
	user string = "user"
)

func (c *Cli) GetStoryteller() ([]types.Storyteller, error) {
	var users []types.Storyteller
	body, err := c.getGeneric(user)
	if err != nil {
		return users, err
	}
	err = json.Unmarshal(body, &users)
	if err != nil {
		return users, err
	}
	return users, nil
}

func (c *Cli) CreateStoryteller(username, userid, password string) ([]byte, error) {
	u := types.Storyteller{
		Username: username,
		Password: password,
	}
	body, err := json.Marshal(u)
	if err != nil {
		return []byte{}, err
	}
	res, err := c.postGeneric(user, body)
	return res, err
}

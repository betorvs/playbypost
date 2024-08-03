package cli

import (
	"encoding/json"

	"github.com/betorvs/playbypost/core/sys/web/types"
)

const (
	storyteller string = "storyteller"
)

func (c *Cli) GetStoryteller() ([]types.Writer, error) {
	var storytellers []types.Writer
	body, err := c.getGeneric(storyteller)
	if err != nil {
		return storytellers, err
	}
	err = json.Unmarshal(body, &storytellers)
	if err != nil {
		return storytellers, err
	}
	return storytellers, nil
}

func (c *Cli) CreateStoryteller(username, password string) ([]byte, error) {
	u := types.Writer{
		Username: username,
		Password: password,
	}
	body, err := json.Marshal(u)
	if err != nil {
		return []byte{}, err
	}
	res, err := c.postGeneric(storyteller, body)
	return res, err
}

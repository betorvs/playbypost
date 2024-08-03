package cli

import (
	"encoding/json"
	"fmt"

	"github.com/betorvs/playbypost/core/rules"
	"github.com/betorvs/playbypost/core/sys/web/types"
)

const (
	player string = "player"
)

func (c *Cli) GeneratePlayer(name string, playerid, stageid int) ([]byte, error) {
	u := types.GeneratePlayer{
		Name:     name,
		PlayerID: playerid,
		StageID:  stageid,
	}
	body, err := json.Marshal(u)
	if err != nil {
		return []byte{}, err
	}
	res, err := c.postGeneric(player, body)
	return res, err
}

func (c *Cli) GetPlayersByStoryID(id int) (map[int]rules.Creature, error) {
	var list map[int]rules.Creature
	play := fmt.Sprintf("%s/story/%d", player, id)
	body, err := c.getGeneric(play)
	if err != nil {
		return list, err
	}
	err = json.Unmarshal(body, &list)
	if err != nil {
		return list, err
	}
	return list, nil
}

func (c *Cli) GetPlayersByID(id int) (rules.Creature, error) {
	var list rules.Creature
	play := fmt.Sprintf("%s/%d", player, id)
	body, err := c.getGeneric(play)
	if err != nil {
		return list, err
	}
	err = json.Unmarshal(body, &list)
	if err != nil {
		return list, err
	}
	return list, nil
}

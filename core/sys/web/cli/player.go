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

func (c *Cli) GeneratePlayer(name, userID string, playerid, stageid int) ([]byte, error) {
	u := types.GeneratePlayer{
		Name:    name,
		StageID: stageid,
	}
	if userID != "" {
		u.UserID = userID
	}
	if playerid != 0 {
		u.PlayerID = playerid
	}
	body, err := json.Marshal(u)
	if err != nil {
		return []byte{}, err
	}
	res, err := c.postGeneric(player, body)
	return res, err
}

func (c *Cli) GetPlayersByStageID(id int) (map[int]rules.Creature, error) {
	var list map[int]rules.Creature
	play := fmt.Sprintf("stage/%s/%d", player, id)
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

func (c *Cli) GetPlayers() ([]types.Players, error) {
	var list []types.Players
	// play := fmt.Sprintf("%s", player)
	body, err := c.getGeneric(player)
	if err != nil {
		return list, err
	}
	err = json.Unmarshal(body, &list)
	if err != nil {
		return list, err
	}
	return list, nil
}

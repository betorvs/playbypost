package parser

import (
	"errors"
	"strconv"
	"strings"
)

// NF is reference to awk NF to get last field
type CommandAction struct {
	ID   int
	Act  string
	Text string
	NF   int
}

func TextToCommand(s string) (CommandAction, error) {
	player := CommandAction{
		Text: s,
	}
	if s != "" {
		p := strings.Split(s, ";")
		length := len(p)
		if length > 0 {
			player.Act = p[1]
			if strings.Contains(p[1], ":") {
				cmds := strings.Split(p[1], ":")
				if len(cmds) > 1 {
					id, err := strconv.Atoi(cmds[1])
					if err == nil {
						player.ID = id
					}
					player.Act = cmds[0]
				}
			}
			if length == 3 {
				player.NF, _ = strconv.Atoi(p[2])
			}
		}
		return player, nil
	}
	return player, errors.New("text cannot be empty")
}

func TextToTaskID(s string) (int, error) {
	if s != "" {
		p := strings.Split(s, ";")
		length := len(p)
		if length > 1 {
			text := strings.Split(p[1], ":")
			if len(text) > 1 {
				id, err := strconv.Atoi(text[1])
				if err == nil {
					return id, nil
				}
			}
		}
	}
	return 0, errors.New("text cannot be empty")
}

package parser

import (
	"errors"
	"fmt"
	"strings"

	"github.com/betorvs/playbypost/core/types"
)

type PlayerAction struct {
	Act    types.Actions
	Text   string
	NotAct string
}

func TextToCommand(s string) (PlayerAction, error) {
	player := PlayerAction{
		Text: s,
	}
	if s != "" {
		p := strings.Split(s, " ")
		length := len(p)
		var command string
		if length > 0 {
			for i := 0; i < length; i++ {
				command = fmt.Sprintf("%s%s ", command, p[i])
				opt, err := types.ActAtoi(strings.TrimSpace(command))
				if err == nil {
					player.Act = opt
					player.NotAct = strings.Join(p[i+1:], " ")
					break
				}
			}
		}
		// var firstCommand, secondCommand, thirdCommand string
		// switch {
		// case length >= 4:
		// 	if p[2] != "" {
		// 		thirdCommand = p[2]
		// 		secondCommand = p[1]
		// 		firstCommand = p[0]
		// 		player.Act = types.Atoi(fmt.Sprintf("%s %s %s", firstCommand, secondCommand, thirdCommand))
		// 	}
		// case length == 3:
		// 	if p[1] != "" {
		// 		secondCommand = p[1]
		// 		firstCommand = p[0]
		// 		player.Act = types.Atoi(fmt.Sprintf("%s %s", firstCommand, secondCommand))
		// 	}
		// case length <= 2:
		// 	if p[0] != "" {
		// 		firstCommand = p[0]
		// 		player.Act = types.Atoi(firstCommand)
		// 	}
		// }
		// // fmt.Println("command", firstCommand, secondCommand, thirdCommand)
		if player.Act == types.NoAction {
			return player, errors.New("invalid")
		}
		return player, nil
	}
	return player, errors.New("text cannot be empty")
}

// func notAct(start int, words []string) string {
// 	// a := strings.Join(words[start:], " ")
// 	// fmt.Println("func notAct ", a)
// 	return strings.Join(words[start:], " ")
// }

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
		// var command string
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
			// for i := 0; i < length; i++ {
			// 	command = fmt.Sprintf("%s%s ", command, p[i])
			// 	opt, err := types.ActAtoi(strings.TrimSpace(command))
			// 	if err == nil {
			// 		player.Act = opt
			// 		player.NotAct = strings.Join(p[i+1:], " ")
			// 		break
			// 	}
			// }
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

// func notAct(start int, words []string) string {
// 	// a := strings.Join(words[start:], " ")
// 	// fmt.Println("func notAct ", a)
// 	return strings.Join(words[start:], " ")
// }

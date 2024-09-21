package cli

import (
	"fmt"
	"strings"
)

func ParserValues(option, startOpt string) (channel, userid, text, display string, err error) {
	splitted := strings.Split(option, ";")
	if len(splitted) == 5 {
		channel = splitted[1]
		userid = splitted[2]
		text = fmt.Sprintf("%s;%s;%s", startOpt, splitted[3], splitted[4])
		display = splitted[3]
		return channel, userid, text, display, nil
	}
	return "", "", "", "", fmt.Errorf("invalid number of fields, expected 5: %s", option)
}

//go:build unit

package parser

import (
	"testing"
)

func TestParser(t *testing.T) {
	options := []string{
		"cmd;change-encounter-to-started:1;2",
	}
	t.Log("table test with command from storyteller")
	{
		for _, v := range options {
			a1, err1 := TextToCommand(v)
			if err1 != nil {
				t.Error("error", err1.Error())
			}
			if a1.ID == 0 {
				t.Error("issue a1 cannot find id")
			}
			t.Log("    command", a1.Text, ", act: ", a1.Act, ", id: ", a1.ID)
		}
	}
}

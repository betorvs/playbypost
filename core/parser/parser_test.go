package parser

import (
	"testing"
)

func TestParser(t *testing.T) {
	// single := "attack orc"
	// twocommands := "ability check strenght"
	// longcommand := "use magic item portion"

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

	// t.Log("should fail tests")
	// {
	// 	c1, err2 := TextToCommand("cast")
	// 	if err2 == nil {
	// 		t.Error("error cannot be nil")
	// 	}
	// 	t.Log("    command", c1.Act.String(), "text ", c1.Text)
	// 	e1, err3 := TextToCommand("")
	// 	if err3 == nil {
	// 		t.Error("error cannot be nil")
	// 	}
	// 	t.Log("    command", e1.Act.String(), "text ", e1.Text)
	// }

	// c1, err3 := TextToCommand(longcommand)
	// if err3 != nil {
	// 	t.Error("error", err3.Error())
	// }
	// if c1.Act == types.NoAction {
	// 	t.Error("issue c1")
	// }

}

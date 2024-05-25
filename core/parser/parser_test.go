package parser

import (
	"testing"

	"github.com/betorvs/playbypost/core/types"
)

func TestParser(t *testing.T) {
	// single := "attack orc"
	// twocommands := "ability check strenght"
	// longcommand := "use magic item portion"

	options := []string{
		"attack orc",
		"ability check strenght",
		"use magic item portion",
	}
	t.Log("table test with possible commands from players")
	{
		for _, v := range options {
			a1, err1 := TextToCommand(v)
			if err1 != nil {
				t.Error("error", err1.Error())
			}
			if a1.Act == types.NoAction {
				t.Errorf("issue %v - %s", a1.Act, a1.Text)
			}
			t.Log("    command", a1.Act.String(), ", not act: ", a1.NotAct)
		}
	}

	t.Log("should fail tests")
	{
		c1, err2 := TextToCommand("cast")
		if err2 == nil {
			t.Error("error cannot be nil")
		}
		t.Log("    command", c1.Act.String(), "text ", c1.Text)
		e1, err3 := TextToCommand("")
		if err3 == nil {
			t.Error("error cannot be nil")
		}
		t.Log("    command", e1.Act.String(), "text ", e1.Text)
	}

	// c1, err3 := TextToCommand(longcommand)
	// if err3 != nil {
	// 	t.Error("error", err3.Error())
	// }
	// if c1.Act == types.NoAction {
	// 	t.Error("issue c1")
	// }

}

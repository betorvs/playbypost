package mechanism

import (
	"testing"
)

func TestRoundControl(t *testing.T) {
	// new pc1 rolls 8
	round1 := NewRoundControl("pc1", 8)
	// pc2 rolls 3
	round1.Append("pc2", 3)
	// pc3 rolls 14
	round1.Append("pc3", 14)
	// pc4 rolls 18
	round1.Append("pc4", 18)
	// npc1 rolls 4
	round1.Append("npc1", 4)
	// npc2 rolls 10
	round1.Append("npc2", 10)
	//
	t.Log("Given a encounter with 6 participants")
	{
		expectedOrderString := "0 - pc4, 1 - pc3, 2 - npc2, 3 - pc1, 4 - npc1, 5 - pc2"
		s := round1.String()
		if expectedOrderString != s {
			t.Errorf("\t Order not expected found %s", s)
		}
		// first to play
		first := round1.Next()
		if first != "pc4" {
			t.Errorf("\t Not expected first player %s", first)
		}
		// last to play
		second := round1.Next()
		if second != "pc3" {
			t.Errorf("\t Not expected second player %s", second)
		}
		// append npc2 again
		err := round1.Append("npc2", 10)
		if err == nil {
			t.Errorf("\t It should not add npc2 again %s", err.Error())
		}
		// loop calling next turns to test reaching the last participant and returning to first
		t.Log("Looping over turns to test cycle")
		for i := 0; i < 4; i++ {
			_ = round1.Next()
			// t.Log("turn of", discard)
		}
		// first to play again
		first = round1.Next()
		if first != "pc4" {
			t.Errorf("\t After full rotation, not expected first player %s", first)
		}
	}
	t.Log("Given a encounter with 5 participants because npc2 dies")
	{
		round1.Remove("npc2")
		expectedOrderString := "0 - pc4, 1 - pc3, 2 - pc1, 3 - npc1, 4 - pc2"
		s := round1.String()
		if expectedOrderString != s {
			t.Errorf("\t Order not expected found %s", s)
		}
	}
	// new encounter with 2 participants with the same rolled value
	round2 := NewRoundControl("pc1", 8)
	// pc2 rolls 3
	round2.Append("pc2", 3)
	// pc3 rolls 14
	round2.Append("pc3", 14)
	// npc1 rolls 4
	round2.Append("npc1", 8)
	// npc2 rolls 10
	round2.Append("npc2", 10)
	t.Log("Given a encounter with 5 participants and two initiative rolls equal")
	{
		// Cannot trust on String method because it can change pc1 and npc1 because they have the same initiative roll
		expectedOrderString := "0 - pc3, 1 - npc2, 2 - pc1, 3 - npc1, 4 - pc2"
		s := round2.String()
		if expectedOrderString != s {
			t.Errorf("\t Order not expected found %s", s)
		}
		// first to play
		first := round2.Next()
		if first != "pc3" {
			t.Errorf("\t Not expected first player %s", first)
		}
		t.Log("Looping over turns to test cycle")
		for i := 0; i < 4; i++ {
			_ = round2.Next()
			// t.Log("turn of", discard)
		}
		// first to play again
		first = round2.Next()
		if first != "pc3" {
			t.Errorf("\t After full rotation, not expected first player %s", first)
		}
	}
}

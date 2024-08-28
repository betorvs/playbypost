//go:build integration

package integration

import (
	"fmt"
	"testing"

	"github.com/betorvs/playbypost/core/sys/web/cli"
	"github.com/betorvs/playbypost/core/sys/web/types"
	"github.com/betorvs/playbypost/core/utils"
)

func TestIntegration(t *testing.T) {
	// Test code here
	t.Log("Integration test")
	// should read from root folder
	creds, err1 := utils.Read("../../../creds")
	if err1 != nil {
		t.Error("error loading creds")
	}
	if creds == "" {
		t.Error("error creds emtpy")
	}
	server := "http://localhost:3000"
	h := cli.NewHeaders(server, "admin", creds)
	if h == nil {
		t.Error("error creating client")
	}
	random := utils.RandomString(6)
	writerUsername := fmt.Sprintf("writer-%s", random)
	_, err2 := h.CreateWriter(writerUsername, "asdQWE123")
	if err2 != nil {
		t.Error("error creating writer")
	}
	writers, err3 := h.GetWriter()
	if err3 != nil {
		t.Error("error getting writers")
	}
	if len(writers) == 0 {
		t.Error("error writers empty")
	}
	writer1 := types.Writer{}
	for _, w := range writers {
		if w.Username == writerUsername {
			writer1 = w
		}
	}
	if writer1.ID == 0 {
		t.Error("error writer1 not found")
	}
	storyTitle := fmt.Sprintf("story-%s", random)
	annouce := fmt.Sprintf("annouce-%s", random)
	note := fmt.Sprintf("note-%s", random)
	_, err4 := h.CreateStory(storyTitle, annouce, note, writer1.ID)
	if err4 != nil {
		t.Error("error creating story")
	}
	stories, err5 := h.GetStory()
	if err5 != nil {
		t.Error("error getting stories")
	}
	if len(stories) == 0 {
		t.Error("error stories empty")
	}
	story1 := types.Story{}
	for _, s := range stories {
		if s.Title == storyTitle {
			story1 = s
		}
	}
	if story1.Announcement == annouce {
		t.Error("error story1 announcement not encrypted")
	}
	encounter1Title := fmt.Sprintf("encounter-1-%s", random)
	encounter1Note := fmt.Sprintf("encounter-1-note-%s", random)
	encounter1Announce := fmt.Sprintf("encounter-1-announce-%s", random)
	_, err6 := h.CreateEncounter(encounter1Title, encounter1Announce, encounter1Note, story1.ID, writer1.ID, true, false)
	if err6 != nil {
		t.Error("error creating encounter 1")
	}
	encounter2Title := fmt.Sprintf("encounter-2-%s", random)
	encounter2Note := fmt.Sprintf("encounter-2-note-%s", random)
	encounter2Announce := fmt.Sprintf("encounter-2-announce-%s", random)
	_, err7 := h.CreateEncounter(encounter2Title, encounter2Announce, encounter2Note, story1.ID, writer1.ID, false, true)
	if err7 != nil {
		t.Error("error creating encounter 2")
	}
	encounter1 := types.Encounter{}
	encounter2 := types.Encounter{}
	encounters, err8 := h.GetEncounters()
	if err8 != nil {
		t.Error("error getting encounters")
	}
	if len(encounters) == 0 {
		t.Error("error encounters empty")
	}
	for _, e := range encounters {
		if e.Title == encounter1Title {
			encounter1 = e
		}
		if e.Title == encounter2Title {
			encounter2 = e
		}
	}
	if encounter1.Announcement == encounter1Announce {
		t.Error("error encounter1 announcement not encrypted")
	}
	if encounter1.FirstEncounter == false {
		t.Error("error encounter1 first encounter not true")
	}
	if encounter2.Announcement == encounter2Announce {
		t.Error("error encounter2 announcement not encrypted")
	}
	if encounter2.LastEncounter == false {
		t.Error("error encounter2 last encounter not true")
	}

}

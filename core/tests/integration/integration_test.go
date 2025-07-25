//go:build integration

package integration_test

//

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/betorvs/playbypost/app/server/handlers/v1/validator"
	"github.com/betorvs/playbypost/core/sys/web/cli"
	"github.com/betorvs/playbypost/core/sys/web/types"
	"github.com/betorvs/playbypost/core/utils"
)

var (
	autoPlayEncounters = map[string][]string{
		"first encounter on story 2":   {"AB, you find it. Choose A or B", "secret note from"},
		"2 encounter on story 2":       {"A. Keep following it.", "secret note from"},
		"3 encounter on story 2":       {"B. Keep following it.", "secret note from"},
		"4 encounter on story 2":       {"After finding A, you got it. Go to end notes.", "secret note from"},
		"bad end encounter on story 2": {"You have got into trouble and this is the end of it.", "secret note from"},
		"last encounter on story 2":    {"You discover it. Thank you, hero! We're proud of you!", "secret note from"},
	}
)

func createBaseStoryEncounter(t *testing.T, h *cli.Cli, name, random string) (types.Writer, types.Story) {
	writerUsername := fmt.Sprintf("writer-%s-%s", name, random)
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
	storyTitle := fmt.Sprintf("story-%s-%s", name, random)
	annouce := fmt.Sprintf("annouce-%s-%s", name, random)
	note := fmt.Sprintf("note-%s-%s", name, random)
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

	return writer1, story1
}

func createBaseStage(t *testing.T, h *cli.Cli, name, random string, storyStage types.Story, writerStage types.Writer) (string, types.Stage) {
	// create stage
	stageText := fmt.Sprintf("stage text %s-%s", name, random)
	storyteller := fmt.Sprintf("storyteller-%s-%s", name, random)
	_, err8 := h.CreateStage(stageText, storyteller, storyStage.ID, writerStage.ID)
	if err8 != nil {
		t.Error("error creating stage")
	}
	stages, err9 := h.GetStage()
	if err9 != nil {
		t.Error("error getting stages")
	}
	stage1 := types.Stage{}
	for _, s := range stages {
		if s.Text == stageText {
			stage1 = s
		}
	}
	if stage1.ID == 0 {
		t.Error("error stage1 not found")
	}
	return storyteller, stage1
}

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

	// create writer, story and encounters and check text encryption
	{
		t.Log("Test create writer, story and encounters and check text encryption - 01")
		random := utils.RandomString(6)
		writer1, story1 := createBaseStoryEncounter(t, h, "test01", random)
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
		t.Log("Test create writer, story and encounters and check text encryption - 01 - ok")
	}

	// auto play
	{
		t.Log("Test auto play - 02")
		random := utils.RandomString(6)
		writerAutoPlay, storyAutoPlay := createBaseStoryEncounter(t, h, "test02", random)
		for k, v := range autoPlayEncounters {
			first := false
			end := false
			if k == "first encounter on story 2" {
				first = true
			}
			if k == "last encounter on story 2" || k == "bad end encounter on story 2" {
				end = true
			}
			title := fmt.Sprintf("%s-%s", k, random)
			_, err6 := h.CreateEncounter(title, v[0], v[1], storyAutoPlay.ID, writerAutoPlay.ID, first, end)
			if err6 != nil {
				t.Error("error creating auto play encounters")
			}
		}
		// create auto play
		autoPlayText := fmt.Sprintf("auto play text %s", random)
		_, err7 := h.CreateAutoPlay(autoPlayText, storyAutoPlay.ID, writerAutoPlay.ID, true)
		if err7 != nil {
			t.Error("error creating auto play")
		}
		autoPlays, err8 := h.GetAutoPlay()
		if err8 != nil {
			t.Error("error getting auto plays")
		}
		if len(autoPlays) == 0 {
			t.Error("error auto plays empty")
		}
		autoPlay1 := types.AutoPlay{}
		for _, a := range autoPlays {
			if a.Text == autoPlayText {
				autoPlay1 = a
			}
		}
		if autoPlay1.ID == 0 {
			t.Error("error auto play not found")
		}
		t.Log("auto play created", "autoPlay1", autoPlay1)
		encounters, err9 := h.GetEncounters()
		if err9 != nil {
			t.Error("error getting encounters")
		}
		autoPlayEncounters := make(map[string]types.Encounter)
		for _, e := range encounters {
			if e.StoryID == storyAutoPlay.ID {
				autoPlayEncounters[e.Title] = e
			}
		}

		// add next encounters
		for k, v := range autoPlayEncounters {
			// title := fmt.Sprintf("%s-%s", k, random)
			if k == fmt.Sprintf("%s-%s", "first encounter on story 2", random) {
				// encounter 2 and 3
				titleNext := fmt.Sprintf("%s-%s", "2 encounter on story 2", random)
				next := types.Next{
					UpstreamID:      autoPlay1.ID,
					EncounterID:     v.ID,
					NextEncounterID: autoPlayEncounters[titleNext].ID,
					Text:            "If you want A",
				}
				_, err9 := h.AddNextEncounter(next)
				if err9 != nil {
					t.Error("error adding next encounter")
				}
				titleNext2 := fmt.Sprintf("%s-%s", "3 encounter on story 2", random)
				next2 := types.Next{
					UpstreamID:      autoPlay1.ID,
					EncounterID:     v.ID,
					NextEncounterID: autoPlayEncounters[titleNext2].ID,
					Text:            "If you want B",
				}
				_, err10 := h.AddNextEncounter(next2)
				if err10 != nil {
					t.Error("error adding next encounter")
				}
			}
			if k == fmt.Sprintf("%s-%s", "2 encounter on story 2", random) {
				// encouter 4
				titleNext := fmt.Sprintf("%s-%s", "4 encounter on story 2", random)
				next := types.Next{
					UpstreamID:      autoPlay1.ID,
					EncounterID:     v.ID,
					NextEncounterID: autoPlayEncounters[titleNext].ID,
					Text:            "moving forward with A",
				}
				_, err11 := h.AddNextEncounter(next)
				if err11 != nil {
					t.Error("error adding next encounter")
				}
			}
			if k == fmt.Sprintf("%s-%s", "3 encounter on story 2", random) {
				// encounter 5
				titleNext := fmt.Sprintf("%s-%s", "bad end encounter on story 2", random)
				next := types.Next{
					UpstreamID:      autoPlay1.ID,
					EncounterID:     v.ID,
					NextEncounterID: autoPlayEncounters[titleNext].ID,
					Text:            "moving forward with B",
				}
				_, err12 := h.AddNextEncounter(next)
				if err12 != nil {
					t.Error("error adding next encounter")
				}
			}
			if k == fmt.Sprintf("%s-%s", "4 encounter on story 2", random) {
				// encounter 6
				titleNext := fmt.Sprintf("%s-%s", "last encounter on story 2", random)
				next := types.Next{
					UpstreamID:      autoPlay1.ID,
					EncounterID:     v.ID,
					NextEncounterID: autoPlayEncounters[titleNext].ID,
					Text:            "go to end notes",
				}
				_, err13 := h.AddNextEncounter(next)
				if err13 != nil {
					t.Error("error adding next encounter")
				}
			}
		}
		// publish auto play
		_, err14 := h.PublishAutoPlay(autoPlay1.ID)
		if err14 != nil {
			t.Error("error publishing auto play")
		}
		playerAutoPlay := fmt.Sprintf("player-auto-play-%s", random)
		channelAutoPlay := fmt.Sprintf("channel-auto-play-%s", random)
		// a.postCommand(userid, "solo-start", i.ChannelID)
		msgSoloStart, err14 := h.PostCommandComposed(playerAutoPlay, types.SoloStart, channelAutoPlay)
		if err14 != nil {
			t.Error("error post solo-start command")
		}
		// parse solo-start response
		for _, m := range msgSoloStart.Opts {
			if m.Name == autoPlay1.Text {
				text := fmt.Sprintf("%s;%s;%d", types.Choice, m.Value, m.ID)
				_, err16 := h.PostCommandComposed(playerAutoPlay, text, channelAutoPlay)
				if err16 != nil {
					t.Error("error post cmd to command")
				}
				break
			}
		}
		// solo-next
		// wait it to be processed
		sleep := 11
		t.Logf("waiting %d seconds to process solo-start command", sleep)
		time.Sleep(time.Duration(sleep) * time.Second)
		msgSoloNext, err16 := h.PostCommandComposed(playerAutoPlay, types.SoloNext, channelAutoPlay)
		if err16 != nil {
			t.Error("error post solo-next command")
		}
		// parse solo-next response
		if len(msgSoloNext.Msg) == 0 {
			t.Log("error solo-next response empty", "msg", msgSoloNext.Msg)
			t.Error("error solo-next response empty")
		}
		t.Log("Test auto play - 02 - ok")
	}

	// stage
	{
		t.Log("Test stage - 03")
		random := utils.RandomString(6)
		writerStage, storyStage := createBaseStoryEncounter(t, h, "test03", random)
		// create encounter
		encounter1Title := fmt.Sprintf("encounter-1-%s", random)
		encounter1Note := fmt.Sprintf("encounter-1-note-%s", random)
		encounter1Announce := fmt.Sprintf("encounter-1-announce-%s", random)
		_, err6 := h.CreateEncounter(encounter1Title, encounter1Announce, encounter1Note, storyStage.ID, writerStage.ID, true, false)
		if err6 != nil {
			t.Error("error creating encounter 1")
		}
		// add task
		taskDescription := fmt.Sprintf("task-investigation-%s", random)
		taskAbility := "wits"
		taskSkill := "investigation"
		taskDifficulty := 2
		taskKind := 2
		_, err7 := h.CreateTask(taskDescription, taskAbility, taskSkill, taskKind, taskDifficulty)
		if err7 != nil {
			t.Error("error creating task")
		}

		// create storyteller, stage
		storyteller, stage1 := createBaseStage(t, h, "test03", random, storyStage, writerStage)

		// add encounter to stage
		encounters, err10 := h.GetEncounters()
		if err10 != nil {
			t.Error("error getting encounters")
		}
		encounterStage := types.Encounter{}
		for _, e := range encounters {
			if e.Title == encounter1Title {
				encounterStage = e
			}
		}
		_, err11 := h.AddEncounterToStage(taskDescription, storyStage.ID, stage1.ID, encounterStage.ID)
		if err11 != nil {
			t.Error("error adding encounter to stage")
		}
		// add participants
		player := fmt.Sprintf("player-%s", random)
		playerID := fmt.Sprintf("player-id-%s", random)
		_, err12 := h.GeneratePlayer(player, playerID, 0, stage1.ID)
		if err12 != nil {
			t.Error("error generating player", "error", err12.Error())
		}
		// start stage
		channelStage := fmt.Sprintf("channel-stage-%s", random)
		_, err13 := h.StartStage(stage1.ID, channelStage)
		if err13 != nil {
			t.Error("error starting stage")
		}
		storytellerMsg, err14 := h.PostCommandComposed(storyteller, "opt", channelStage)
		if err14 != nil {
			t.Error("error post opt command")
		}
		if len(storytellerMsg.Opts) == 0 {
			t.Error("error storyteller response empty")
		}
		for _, m := range storytellerMsg.Opts {
			if m.Name == encounter1Title {
				text := fmt.Sprintf("cmd;%s;%d", m.Value, m.ID)
				_, err16 := h.PostCommandComposed(player, text, channelStage)
				if err16 != nil {
					t.Error("error post storyteller cmd to command")
				}
				break
			}
		}
		// wait it to be processed
		sleep := 11
		t.Logf("waiting %d seconds to process cmd command", sleep)
		time.Sleep(time.Duration(sleep) * time.Second)
		storytellerMsg2, err17 := h.PostCommandComposed(storyteller, "opt", channelStage)
		if err17 != nil {
			t.Error("error post storyteller 2 opt command")
		}
		if len(storytellerMsg2.Opts) == 0 {
			t.Error("error storyteller 2 response empty")
		}
		t.Log("Test stage - 03 - ok")
	}
	// validator
	{
		t.Log("Test validator - 04")
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
		writerValidator := types.Writer{}
		for _, w := range writers {
			if w.Username == writerUsername {
				writerValidator = w
			}
		}
		if writerValidator.ID == 0 {
			t.Error("error writerValidator not found")
		}
		storyEmptyTitle := fmt.Sprintf("story-empty-%s", random)
		annouceEmpty := fmt.Sprintf("annouce-empty-%s", random)
		noteEmpty := fmt.Sprintf("note-empty-%s", random)
		_, err4 := h.CreateStory(storyEmptyTitle, annouceEmpty, noteEmpty, writerValidator.ID)
		if err4 != nil {
			t.Error("error creating empty story")
		}
		stories, err5 := h.GetStory()
		if err5 != nil {
			t.Error("error getting stories")
		}
		storyEmpty := types.Story{}
		for _, s := range stories {
			if s.Title == storyEmptyTitle {
				storyEmpty = s
			}
		}
		if storyEmpty.ID == 0 {
			t.Error("error storyEmpty not found")
		}
		// call validator
		_, err6 := h.ValidatorPut("story", storyEmpty.ID)
		if err6 != nil {
			t.Error("error calling validator")
		}
		time.Sleep(5 * time.Second)
		body1, err7 := h.GetValidator()
		if err7 != nil {
			t.Error("error getting validator")
		}
		var obj1 []validator.Request
		err8 := json.Unmarshal(body1, &obj1)
		if err8 != nil {
			t.Error("error unmarshal validator")
		}
		for _, v := range obj1 {
			if v.Kind == "story" {
				if v.ID == storyEmpty.ID {
					if v.Valid != false {
						t.Error("error validator story ID")
					}
					if len(v.Analise.Results) == 0 {
						t.Error("error validator story results")
					}
				}
			}
		}
		t.Log("Test validator - 04 - ok")
	}
	// delete writer user association
	{
		t.Log("Test delete writer user association - 05")
		random := utils.RandomString(6)
		writerDelete, storyDelete := createBaseStoryEncounter(t, h, "test05", random)
		_, stageDelete := createBaseStage(t, h, "test05", random, storyDelete, writerDelete)

		playerDeleteUsername := fmt.Sprintf("player-delete-%s", random)
		playerDeleteID := fmt.Sprintf("player-delete-id-%s", random)
		_, err1 := h.GeneratePlayer(playerDeleteUsername, playerDeleteID, 0, stageDelete.ID)
		if err1 != nil {
			t.Errorf("error creating user: %v", err1)
		}
		users, err := h.GetPlayers()
		if err != nil {
			t.Errorf("error getting users: %v", err)
		}
		user := types.Players{}
		for _, u := range users {
			if u.Name == playerDeleteUsername {
				user = u
			}
		}
		if user.ID == 0 {
			t.Error("error user not found")
		}

		_, err = h.CreateWriterUserAssociation(writerDelete.ID, user.ID)
		if err != nil {
			t.Errorf("error creating writer user association: %v", err)
		}

		// Verify association exists before deletion
		associationExists, err := h.CheckWriterUserAssociationExists(writerDelete.ID, user.ID)
		if err != nil {
			t.Errorf("error checking writer user association: %v", err)
		}
		if !associationExists {
			t.Error("writer user association should exist before deletion")
		}
		// get all associations
		associations, err := h.GetWriterUsersAssociation()
		if err != nil {
			t.Errorf("error getting writer user association: %v", err)
		}
		association := types.WriterUserAssociation{}
		for _, a := range associations {
			t.Log("association", "id", a.ID, "writer_id", a.WriterID, "user_id", a.UserID)
			if a.WriterID == writerDelete.ID && a.UserID == user.ID {
				association = a
			}
		}
		if association.ID == 0 {
			t.Error("error association not found")
		}

		// Delete the association
		err = h.DeleteWriterUserAssociation(association.ID)
		if err != nil {
			t.Errorf("error deleting writer user association: %v", err)
		}
		// Verify association does not exist after deletion
		associationExists, err = h.CheckWriterUserAssociationExists(writerDelete.ID, user.ID)
		if err != nil {
			t.Errorf("error checking writer user association after deletion: %v", err)
		}
		if associationExists {
			t.Error("writer user association should not exist after deletion")
		}
		t.Log("Test delete writer user association - 05 - ok")
	}
}

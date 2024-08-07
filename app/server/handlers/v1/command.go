package v1

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"encoding/json"

	"github.com/betorvs/playbypost/core/parser"
	"github.com/betorvs/playbypost/core/sys/web/types"
)

func (a MainApi) ExecuteCommand(w http.ResponseWriter, r *http.Request) {
	headerUserID := r.Header.Get(types.HeaderUserID)
	headerStoryChannel := r.Header.Get(types.HeaderStoryChannel)
	if headerUserID == "" || headerStoryChannel == "" {
		a.s.ErrJSON(w, http.StatusForbidden, "required authentication headers")
		return
	}
	obj := types.Command{}
	err := json.NewDecoder(r.Body).Decode(&obj)
	if err != nil {
		a.s.ErrJSON(w, http.StatusBadRequest, "json decode error")
		return
	}
	a.logger.Info("command received", "command", obj.Text, "userid", headerUserID, "channel", headerStoryChannel)
	runningStage, err := a.db.GetRunningStageByChannelID(a.ctx, headerStoryChannel, headerUserID)
	if err != nil {
		a.s.ErrJSON(w, http.StatusBadRequest, "invalid userid")
		return
	}
	a.logger.Info("running stage", "runningStage", runningStage)
	storyteller := false
	if runningStage.Stage.UserID == headerUserID {
		storyteller = true
	}
	switch {
	case strings.Contains(strings.ToLower(obj.Text), "opt"):
		msg := fmt.Sprintf("Player '%s' ", runningStage.Players.Name)
		if storyteller {
			msg = fmt.Sprintf("Storyteller '%s' ", runningStage.Stage.UserID)
		}
		// parse options
		encOptions := parser.ParserOptions(storyteller, runningStage)

		composed := types.Composed{Msg: msg, Opt: encOptions}
		a.logger.Info("msg back", "composed", composed)
		a.s.JSON(w, composed)
		return
	case strings.HasPrefix(strings.ToLower(obj.Text), "cmd"):
		actions := types.NewActions()
		cmd, err := parser.TextToCommand(obj.Text)
		if err != nil {
			a.s.ErrJSON(w, http.StatusBadRequest, "command error")
			return
		}
		actions["command"] = cmd.Act
		actions["text"] = cmd.Text
		actions["channel"] = headerStoryChannel
		actions["userid"] = headerUserID
		foundID, err := parser.TextToTaskID(obj.Text)
		if err != nil {
			a.s.ErrJSON(w, http.StatusBadRequest, "cannot find id in command text")
			return
		}
		if strings.HasPrefix(cmd.Act, parser.Task) {
			actions["task_id"] = strconv.Itoa(foundID)
		}
		if strings.HasPrefix(cmd.Act, parser.AttackPlayer) && cmd.NF != 0 {
			actions["npc_id"] = strconv.Itoa(cmd.NF)
			actions["player_id"] = strconv.Itoa(foundID)
		}
		if strings.Contains(cmd.Act, fmt.Sprintf("%s-npc", parser.HealthStatus)) {
			actions["npc_id"] = strconv.Itoa(foundID)
		}
		if strings.HasPrefix(cmd.Act, parser.AttackNPC) {
			actions["npc_id"] = strconv.Itoa(foundID)
		}
		if runningStage.Encounter.InitiativeID != 0 {
			actions["initiative_id"] = strconv.Itoa(runningStage.Encounter.InitiativeID)
		}

		encounterID := runningStage.Encounter.ID
		// it should recover encounter id from command in case of change encounter phase
		if strings.HasPrefix(cmd.Act, parser.ChangeEncounter) && cmd.ID > 0 {
			encounterID = cmd.ID
		}
		actions["encounter_id"] = strconv.Itoa(encounterID)
		if !storyteller {
			actions["player_id"] = strconv.Itoa(runningStage.Players.ID)
		}

		err = a.db.RegisterActivities(a.ctx, runningStage.Stage.ID, encounterID, actions)
		if err != nil {
			a.logger.Error("register activities error", "error", err.Error(), "encounterID", encounterID, "actions", actions)
			a.s.ErrJSON(w, http.StatusBadRequest, "register activities")
			return
		}
		msg := "command accepted"
		if runningStage.Encounter.InitiativeID != 0 {
			msg = fmt.Sprintf("command accepted, initiative id %d found. It will check initiative order before rolling your action.", runningStage.Encounter.InitiativeID)
		}
		a.s.JSON(w, types.Msg{Msg: msg})
		return
	}
	// msg := fmt.Sprintf("player found '%s' and story id found '%d' ", player.Name, scene.Story.ID)
	msg := "no options for you"
	a.s.JSON(w, types.Msg{Msg: msg})
}

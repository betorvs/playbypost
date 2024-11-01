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
	text := strings.ToLower(obj.Text)
	a.logger.Debug("command received", "command", text, "userid", headerUserID, "channel", headerStoryChannel)

	switch {
	case strings.HasPrefix(text, types.SoloDescribe), strings.HasPrefix(text, types.DidaticDescribe): // describe
		solo := true
		if strings.HasPrefix(text, types.DidaticDescribe) {
			solo = false
		}

		options, err := a.db.DescribeAutoPlayPublished(a.ctx, solo)
		if err != nil {
			a.s.ErrJSON(w, http.StatusBadRequest, "describe auto play")
			return
		}
		opts := parser.ParseAutoPlayDescribe(options)
		composed := types.Composed{Msg: "options", Opts: opts}
		a.logger.Info("msg back", "composed", composed)
		a.s.JSON(w, composed)
		return

	case strings.HasPrefix(text, types.SoloStart), strings.HasPrefix(text, types.DidaticStart): // solo-start
		// solo mode: decide workflow
		// list all available solo modes
		// will return a list of auto_play entries
		options, err := a.db.GetAutoPlay(a.ctx)
		if err != nil {
			a.s.ErrJSON(w, http.StatusBadRequest, "get auto play")
			return
		}
		opts := parser.ParserAutoPlays(options, text)
		composed := types.Composed{Msg: "options", Opts: opts}
		a.logger.Info("msg back", "composed", composed)
		a.s.JSON(w, composed)
		return

	case strings.HasPrefix(text, types.DidaticJoin): // FIX
		// didatic mode: decide workflow
		// list all available didatic modes
		// will return a list of auto_play entries
		// headerStoryChannel, headerUserID
		err := a.db.AddAutoPlayGroup(a.ctx, headerStoryChannel, headerUserID)
		if err != nil {
			a.s.ErrJSON(w, http.StatusBadRequest, "add auto play group")
			return
		}
		composed := types.Composed{Msg: "Didatic join message received"}
		a.s.JSON(w, composed)
		return

	case strings.HasPrefix(text, types.SoloNext), strings.HasPrefix(text, types.DidaticNext): // solo-next
		// solo next mode: requires channel id and user id
		opt, err := a.getAutoPlayOptByChannelID(headerStoryChannel, headerUserID)
		if err != nil {
			a.logger.Error("no auto play found", "error", err.Error())
			a.s.ErrJSON(w, http.StatusBadRequest, "no auto play found")
			return
		}
		a.logger.Debug("auto play found", "opt", opt)

		composed := types.Composed{Msg: "next options"}
		if len(opt.NextEncounters) > 0 {
			opts, _ := parser.ParserAutoPlaysNext(opt.NextEncounters, opt.AutoPlay.Solo)
			a.logger.Info("auto play found", "opts", opts)
			composed.Opts = opts
		}

		a.s.JSON(w, composed)
		return

	case strings.Contains(text, types.Opt): // opt
		runningStage, storyteller, err := a.getRunningStageByChannelID(headerStoryChannel, headerUserID)
		if err != nil {
			a.s.ErrJSON(w, http.StatusBadRequest, "invalid userid")
			return
		}
		msg := fmt.Sprintf("Player '%s' ", runningStage.Players.Name)
		if storyteller {
			msg = fmt.Sprintf("Storyteller '%s' ", runningStage.Stage.UserID)
		}
		// parse options
		opts := parser.ParserOptions(storyteller, runningStage)

		composed := types.Composed{Msg: msg, Opts: opts}
		a.logger.Debug("msg back", "composed", composed)
		a.s.JSON(w, composed)
		return
	case strings.HasPrefix(text, types.Cmd): //cmd
		runningStage, storyteller, err := a.getRunningStageByChannelID(headerStoryChannel, headerUserID)
		if err != nil {
			a.s.ErrJSON(w, http.StatusBadRequest, "invalid userid")
			return
		}
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

	case strings.HasPrefix(text, types.Choice), strings.HasPrefix(text, types.Decision): // choice, decision
		// process solo choices for player
		// start solo mode: requires channel id and user id
		// next solo mode: requires channel id and user id
		a.logger.Info("command received", "text", text)
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
		a.logger.Debug("command parsed", "cmd", cmd)
		// command=choice;start-solo-solo-adventure-1:1;1 userid=1272952428379242611 channel=1275626175678517289
		// err = a.db.RegisterActivitiesAutoPlay(a.ctx, cmd.ID, cmd.NF, actions)
		switch {
		case strings.HasPrefix(cmd.Act, parser.JoinDidatic): // join didatic
			actions["auto_play_id"] = strconv.Itoa(cmd.ID)
			// add auto play
			encounterID, err := a.db.CreateAutoPlayChannelTx(a.ctx, headerStoryChannel, headerUserID, cmd.ID)
			if err != nil {
				a.s.ErrJSON(w, http.StatusBadRequest, "create auto play channel")
				return
			}
			actions["encounter_id"] = strconv.Itoa(encounterID)
			// create registry
			err = a.db.RegisterActivitiesAutoPlay(a.ctx, cmd.ID, encounterID, actions)
			if err != nil {
				a.s.ErrJSON(w, http.StatusBadRequest, "register activities auto play")
				return
			}

		case strings.HasPrefix(cmd.Act, parser.StartSolo), strings.HasPrefix(cmd.Act, parser.StartDidatic): // solo-start
			actions["auto_play_id"] = strconv.Itoa(cmd.ID)
			// add auto play
			encounterID, err := a.db.CreateAutoPlayChannelTx(a.ctx, headerStoryChannel, headerUserID, cmd.ID)
			if err != nil {
				a.s.ErrJSON(w, http.StatusBadRequest, "create auto play channel")
				return
			}
			actions["encounter_id"] = strconv.Itoa(encounterID)
			// create registry
			err = a.db.RegisterActivitiesAutoPlay(a.ctx, cmd.ID, encounterID, actions)
			if err != nil {
				a.s.ErrJSON(w, http.StatusBadRequest, "register activities auto play")
				return
			}

		case strings.HasPrefix(cmd.Act, parser.NextSolo), strings.HasPrefix(cmd.Act, parser.NextDidatic):
			// {ID:8 Act:next-solo-for-A-go-enc-2 Text:choice;next-solo-for-A-go-enc-2:8;1 NF:1}"
			actions["auto_play_id"] = strconv.Itoa(cmd.NF)
			actions["encounter_id"] = strconv.Itoa(cmd.ID)
			err = a.db.RegisterActivitiesAutoPlay(a.ctx, cmd.NF, cmd.ID, actions)
			if err != nil {
				a.s.ErrJSON(w, http.StatusBadRequest, "register activities auto play")
				return
			}

		case strings.HasPrefix(cmd.Act, parser.DiceRollSolo):
			//
			actions["auto_play_id"] = strconv.Itoa(cmd.NF)
			actions["encounter_id"] = strconv.Itoa(cmd.ID)
			err = a.db.RegisterActivitiesAutoPlay(a.ctx, cmd.NF, cmd.ID, actions)
			if err != nil {
				a.s.ErrJSON(w, http.StatusBadRequest, "register activities auto play")
				return
			}
			// end switch
		}

		msg := "command accepted"
		a.s.JSON(w, types.Msg{Msg: msg})
		return

		// case strings.HasPrefix(text, types.Decision): // decision
		// 	// process didatic commands
		// 	a.logger.Info("decision command received", "text", text)
		// 	msg := "command accepted"
		// 	a.s.JSON(w, types.Msg{Msg: msg})
		// 	return
		// end switch case
	}
	// msg := fmt.Sprintf("player found '%s' and story id found '%d' ", player.Name, scene.Story.ID)
	msg := "no options for you"
	a.s.JSON(w, types.Msg{Msg: msg})
}

func (a MainApi) getRunningStageByChannelID(channelID, userID string) (types.RunningStage, bool, error) {
	storyteller := false
	runningStage, err := a.db.GetRunningStageByChannelID(a.ctx, channelID, userID, a.rpg)
	if err != nil {
		return types.RunningStage{}, storyteller, err
	}
	if runningStage.Stage.UserID == userID {
		storyteller = true
	}
	return runningStage, storyteller, nil
}

func (a MainApi) getAutoPlayOptByChannelID(channel, userID string) (types.AutoPlayOptions, error) {
	opt, err := a.db.GetAutoPlayOptionsByChannelID(a.ctx, channel, userID)
	if err != nil {
		return types.AutoPlayOptions{}, err
	}
	if opt.ID == 0 {
		return types.AutoPlayOptions{}, fmt.Errorf("no auto play found")
	}
	return opt, nil
}

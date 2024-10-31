package worker

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/betorvs/playbypost/core/parser"
	"github.com/betorvs/playbypost/core/sys/web/types"
	"golang.org/x/exp/slices"
)

func (a *WorkerAPI) parseAutoPlayCommand(cmd types.Activity) error {
	// Add your code here to handle the auto play command
	announce, last, err := a.db.GetAnnounceByEncounterID(a.ctx, cmd.EncounterID, cmd.UpstreamID)
	if err != nil {
		a.logger.Error("error getting announce by encounter id", "error", err.Error())
		return err
	}
	a.logger.Info("announce found", "announce", announce)
	processed := false
	switch {
	// case strings.HasPrefix(cmd.Actions["command"], parser.JoinDidatic):
	// 	// join didatic mode
	// 	a.logger.Info("join mode", "command", cmd.Actions["command"])
	// 	// send message to chat
	// 	joinMessage := fmt.Sprintf("Joining didatic: %s. Wait all collegues to join and then call: didatic-start", cmd.Actions["userid"])
	// 	body, err := a.client.PostEvent(cmd.Actions["channel"], "ALL", joinMessage, types.EventInformation)
	// 	if err != nil {
	// 		a.logger.Error("error posting event join didatic mode", "error", err.Error(), "body", string(body))
	// 		return err
	// 	}
	// 	cmd.Actions["result"] = joinMessage
	// 	processed = true

	case strings.HasPrefix(cmd.Actions["command"], parser.StartSolo), strings.HasPrefix(cmd.Actions["command"], parser.StartDidatic):
		// start solo mode
		a.logger.Info("start mode", "command", cmd.Actions["command"])
		// send announce from story too
		title, storyAnnouncement, err := a.db.GetStoryAnnouncementByAutoPlayID(a.ctx, cmd.UpstreamID)
		if err != nil {
			a.logger.Error("error getting story announcement by auto play id", "error", err.Error())
			return err
		}
		// send message to chat
		user := cmd.Actions["userid"]
		if cmd.Actions["command"] == parser.StartDidatic {
			user = "ALL"
		}
		announceStory := fmt.Sprintf("Starting: %s\n%s", title, storyAnnouncement)
		body, err := a.client.PostEvent(cmd.Actions["channel"], user, announceStory, types.EventInformation)
		if err != nil {
			a.logger.Error("error posting event start auto play mode", "error", err.Error(), "body", string(body))
			return err
		}
		cmd.Actions["story_announcement"] = announceStory
		// send message to chat
		body, err = a.client.PostEvent(cmd.Actions["channel"], user, announce, types.EventInformation)
		if err != nil {
			a.logger.Error("error posting event start auto play mode", "error", err.Error(), "body", string(body))
			return err
		}

		// register encounter used
		err = a.db.UpdateAutoPlayState(a.ctx, cmd.Actions["channel"], cmd.EncounterID)
		if err != nil {
			a.logger.Error("error updating auto play state", "error", err.Error())
			return err
		}
		cmd.Actions["result"] = announce

		processed = true

	case strings.HasPrefix(cmd.Actions["command"], parser.NextSolo), strings.HasPrefix(cmd.Actions["command"], parser.NextDidatic):
		// next solo mode
		a.logger.Info("next mode", "command", cmd.Actions["command"])
		user := cmd.Actions["userid"]
		if strings.HasPrefix(cmd.Actions["command"], parser.NextDidatic) {
			user = "ALL"
		}
		body, err := a.client.PostEvent(cmd.Actions["channel"], user, announce, types.EventInformation)
		if err != nil {
			a.logger.Error("error posting event start solo mode", "error", err.Error(), "body", string(body))
			return err
		}

		// register encounter used
		err = a.db.UpdateAutoPlayState(a.ctx, cmd.Actions["channel"], cmd.EncounterID)
		if err != nil {
			a.logger.Error("error updating auto play state", "error", err.Error())
			return err
		}
		cmd.Actions["result"] = announce
		if last {
			// finish auto play
			a.logger.Info("finish auto play")
			// send message to chat
			body, err := a.client.PostEvent(cmd.Actions["channel"], user, "End of this game. Congratulations", types.EventEnd)
			if err != nil {
				a.logger.Error("error posting event end solo mode", "error", err.Error(), "body", string(body))
				return err
			}
			// close channel
			err = a.db.CloseAutoPlayChannel(a.ctx, cmd.Actions["channel"], cmd.UpstreamID)
			if err != nil {
				a.logger.Error("error closing auto play channel", "error", err.Error())
				return err
			}
			cmd.Actions["end"] = "true"
		}
		processed = true

	case strings.HasPrefix(cmd.Actions["command"], parser.DiceRollSolo):
		// dice roll solo mode
		a.logger.Info("dice roll solo mode")
		// roll a dice
		rolled, err := a.dice.FreeRoll("free-dice-roll", a.autoPlay.BaseDice)
		if err != nil {
			a.logger.Error("error rolling dice", "error", err.Error())
			return err
		}
		// send message to chat
		msg := fmt.Sprintf("Solo Dice Result (%s) rolled: %d", a.autoPlay.BaseDice, rolled.Result)
		body, err := a.client.PostEvent(cmd.Actions["channel"], cmd.Actions["userid"], msg, types.EventSuccess)
		if err != nil {
			a.logger.Error("error posting event dice roll solo mode", "error", err.Error(), "body", string(body))
			return err
		}
		next, err := a.db.GetAutoPlayOptionsByChannelID(a.ctx, cmd.Actions["channel"], cmd.Actions["userid"])
		if err != nil {
			a.logger.Error("error getting auto play next by auto play id", "error", err.Error())
			return err
		}
		if len(next.NextEncounters) == 0 {
			a.logger.Info("no next encounters found")
			return fmt.Errorf("no next encounters found")
		}

		nextEncounterID := 0
		for _, n := range next.NextEncounters {
			a.logger.Info("next encounter", "kind", n.Objective.Kind, "values", n.Objective.Values)
			if n.Objective.Kind == types.ObjectiveDiceRoll && slices.Contains(n.Objective.Values, rolled.Result) {
				nextEncounterID = n.NextEncounterID
			}
		}

		// register encounter used
		if nextEncounterID != 0 {
			// add to registry
			cmd.Actions["text"] = fmt.Sprintf("%s;%s;%s", "choice", cmd.Actions["channel"], cmd.Actions["userid"])
			cmd.Actions["auto_play_id"] = strconv.Itoa(cmd.UpstreamID)
			cmd.Actions["encounter_id"] = strconv.Itoa(nextEncounterID)
			cmd.Actions["command"] = parser.NextSolo
			err = a.db.RegisterActivitiesAutoPlay(a.ctx, cmd.UpstreamID, nextEncounterID, cmd.Actions)
			if err != nil {
				a.logger.Error("error registering activities auto play", "error", err.Error())
				return err
			}

			cmd.Actions["result_dice"] = msg
			processed = true
		}

	}
	if processed {
		err := a.db.UpdateProcessedAutoPlay(a.ctx, cmd.ID, true, cmd.Actions)
		if err != nil {
			a.logger.Error("error updating auto play activity processed", "error", err.Error())
			return err
		}
	}

	return nil
}

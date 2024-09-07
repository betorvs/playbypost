package worker

import (
	"strings"

	"github.com/betorvs/playbypost/core/parser"
	"github.com/betorvs/playbypost/core/sys/web/types"
)

func (a *WorkerAPI) parseAutoPlayCommand(cmd types.AutoPlayActivities) error {
	// Add your code here to handle the auto play command
	announce, last, err := a.db.GetAnnounceByEncounterID(a.ctx, cmd.EncounterID, cmd.AutoPlayID)
	if err != nil {
		a.logger.Error("error getting announce by encounter id", "error", err.Error())
		return err
	}
	a.logger.Info("announce found", "announce", announce)
	processed := false
	switch {
	case strings.HasPrefix(cmd.Actions["command"], parser.StartSolo):
		// start solo mode
		a.logger.Info("start solo mode")
		// send message to chat
		body, err := a.client.PostEvent(cmd.Actions["channel"], cmd.Actions["userid"], announce, types.EventInformation)
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
		processed = true

	case strings.HasPrefix(cmd.Actions["command"], parser.NextSolo):
		// next solo mode
		a.logger.Info("next solo mode")
		body, err := a.client.PostEvent(cmd.Actions["channel"], cmd.Actions["userid"], announce, types.EventInformation)
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
			body, err := a.client.PostEvent(cmd.Actions["channel"], cmd.Actions["userid"], "End of solo game. Congratulations", types.EventEnd)
			if err != nil {
				a.logger.Error("error posting event end solo mode", "error", err.Error(), "body", string(body))
				return err
			}
			// close channel
			err = a.db.CloseAutoPlayChannel(a.ctx, cmd.Actions["channel"], cmd.AutoPlayID)
			if err != nil {
				a.logger.Error("error closing auto play channel", "error", err.Error())
				return err
			}
			cmd.Actions["end"] = "true"
		}
		processed = true

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

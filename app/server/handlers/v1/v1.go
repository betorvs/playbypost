package v1

import (
	"context"
	"log/slog"

	"github.com/betorvs/playbypost/app/server/handlers"
	"github.com/betorvs/playbypost/app/server/worker"
	"github.com/betorvs/playbypost/core/rpg"
	"github.com/betorvs/playbypost/core/sys/db"
	"github.com/betorvs/playbypost/core/sys/web/cli"
	"github.com/betorvs/playbypost/core/sys/web/server"
)

type MainApi struct {
	Session  *handlers.Session
	Worker   *worker.WorkerAPI
	logger   *slog.Logger
	s        *server.SvrWeb
	db       db.DBClient
	ctx      context.Context
	dice     rpg.Roll
	client   *cli.Cli
	rpg      *rpg.RPGSystem
	autoPlay *rpg.RPGSystem
}

func NewMainApi(ctx context.Context, dice rpg.Roll, db db.DBClient, l *slog.Logger, s *server.SvrWeb, client *cli.Cli, rpgSystem *rpg.RPGSystem) *MainApi {
	session := handlers.NewSession(l, db, s, ctx)
	autoPlay := rpg.LoadRPGSystemsDefault(rpg.AutoPlay)
	worker := worker.NewWorkerAPI(ctx, dice, db, l, client, rpgSystem, autoPlay)
	return &MainApi{
		Session:  session,
		Worker:   worker,
		ctx:      ctx,
		dice:     dice,
		db:       db,
		logger:   l,
		s:        s,
		client:   client,
		rpg:      rpgSystem,
		autoPlay: autoPlay,
	}
}

// func (a *MainApi) Execute() {
// 	a.logger.Info("starting scheduler main api execution", "time", time.Now())
// 	activities, err := a.db.GetStageEncounterActivities(a.ctx)
// 	if err != nil {
// 		a.logger.Error("error getting stage encounter activities", "error", err.Error())
// 		return
// 	}
// 	for _, activity := range activities {
// 		if !activity.Processed {
// 			a.logger.Info("activity", "activity", activity)
// 			// execute activity
// 			err := a.parseCommand(activity)
// 			if err != nil {
// 				a.logger.Error("error parsing command", "error", err.Error())
// 				continue
// 			}

// 		}
// 	}
// }

// func (a *MainApi) parseCommand(cmd types.StageEncounterActivities) error {
// 	// call back to slack
// 	enc, err := a.db.GetStageEncounterByEncounterID(a.ctx, cmd.EncounterID)
// 	if err != nil {
// 		a.logger.Error("error getting encounter by id", "error", err.Error())
// 		return err
// 	}
// 	a.logger.Info("encounter found in parseCommand", "encounter", enc)
// 	processed := false
// 	// parse command
// 	switch {
// 	case cmd.Actions["command"] == parser.ChangeEncounterToStarted:
// 		// change encounter to started
// 		a.logger.Info("change encounter to started")
// 		err := a.db.UpdatePhase(a.ctx, cmd.EncounterID, int(types.Started))
// 		if err != nil {
// 			a.logger.Error("error updating encounter phase to started", "error", err.Error())
// 			return err
// 		}
// 		encounterAnnounce := fmt.Sprintf("Encounter %s has started: _%s_", enc.Title, enc.Announcement)
// 		body, err := a.client.PostEvent(cmd.Actions["channel"], "ALL", encounterAnnounce, types.EventAnnounce)
// 		if err != nil {
// 			a.logger.Error("error posting event start encounter", "error", err.Error(), "body", string(body))
// 			return err
// 		}
// 		processed = true

// 	case cmd.Actions["command"] == parser.ChangeEncounterToRunning:
// 		// change encounter to running
// 		a.logger.Info("change encounter to running")
// 		err := a.db.UpdatePhase(a.ctx, cmd.EncounterID, int(types.Running))
// 		if err != nil {
// 			a.logger.Error("error updating encounter phase to running", "error", err.Error())
// 			return err
// 		}
// 		// call back to slack
// 		encounterAnnounce := fmt.Sprintf("Encounter %s is running. Make your actions!", enc.Title)
// 		body, err := a.client.PostEvent(cmd.Actions["channel"], "ALL", encounterAnnounce, types.EventAnnounce)
// 		if err != nil {
// 			a.logger.Error("error posting event running encounter", "error", err.Error(), "body", string(body))
// 			return err
// 		}
// 		processed = true

// 	case cmd.Actions["command"] == parser.ChangeEncounterToFinished:
// 		// change encounter to finished
// 		a.logger.Info("change encounter to finished")
// 		err := a.db.UpdatePhase(a.ctx, cmd.EncounterID, int(types.Finished))
// 		if err != nil {
// 			a.logger.Error("error updating encounter phase to finished", "error", err.Error())
// 			return err
// 		}
// 		// call back to slack
// 		encounterAnnounce := fmt.Sprintf("Encounter %s finished.", enc.Title)
// 		body, err := a.client.PostEvent(cmd.Actions["channel"], "ALL", encounterAnnounce, types.EventAnnounce)
// 		if err != nil {
// 			a.logger.Error("error posting event finishing encounter", "error", err.Error(), "body", string(body))
// 			return err
// 		}
// 		// get next encounter and start it

// 		processed = true

// 	case cmd.Actions["command"] == parser.RollInitiative:
// 		a.logger.Info("roll initiative")

// 	case strings.HasPrefix(cmd.Actions["command"], parser.AttackNPC):
// 		a.logger.Info("attack npc")

// 	case strings.HasPrefix(cmd.Actions["command"], parser.AttackPlayer):
// 		a.logger.Info("attack player")

// 	case strings.HasPrefix(cmd.Actions["command"], parser.Task):
// 		a.logger.Info("task command")
// 		player, err := a.db.GetPlayerByUserID(a.ctx, cmd.Actions["userid"], cmd.Actions["channel"])
// 		if err != nil {
// 			a.logger.Error("error getting player by user id", "error", err.Error())
// 			return err
// 		}
// 		a.logger.Info("player found", "player", player)
// 		a.logger.Info("rpg system", "rpg", a.rpg.BaseSystem)
// 		creature := rules.RestoreCreature()
// 		creature.RPG = a.rpg
// 		types.PlayerToCreature(&player, creature)
// 		a.logger.Info("creature found", "creature", creature)
// 		// get task
// 		taskID, err := strconv.Atoi(cmd.Actions["task_id"])
// 		if err != nil {
// 			a.logger.Error("error getting task id", "error", err.Error())
// 			return err
// 		}
// 		task, err := a.db.GetSTaskFromRunningTaskID(a.ctx, taskID)
// 		if err != nil {
// 			a.logger.Error("error getting task id", "error", err.Error())
// 			return err
// 		}
// 		a.logger.Info("task found", "task", task)
// 		result, err := a.executeTask(task, creature)
// 		if err != nil {
// 			a.logger.Error("error executing task", "error", err.Error())
// 			return err
// 		}
// 		// call back to slack
// 		eventKind := types.EventFailure
// 		if result.Success {
// 			eventKind = types.EventSuccess
// 		}
// 		resultMessage := fmt.Sprintf("Result: %d with dices rolled: %s", result.Result, result.Rolled)
// 		taskResult := fmt.Sprintf("Task %s executed by @%s. %s", task.Description, cmd.Actions["userid"], resultMessage)
// 		body, err := a.client.PostEvent(cmd.Actions["channel"], "ALL", taskResult, eventKind)
// 		if err != nil {
// 			a.logger.Error("error posting event task result", "error", err.Error(), "body", string(body))
// 			return err
// 		}
// 		cmd.Actions["result"] = resultMessage
// 		processed = true

// 	default:
// 		a.logger.Info("command not defined")
// 	}
// 	// update activity processed
// 	if processed {
// 		err = a.db.UpdateProcessedActivities(a.ctx, cmd.ID, true, cmd.Actions)
// 		if err != nil {
// 			a.logger.Error("error updating stage encounter activities", "error", err.Error())
// 			return err
// 		}
// 	}

// 	return nil
// }

// func (a *MainApi) executeTask(task types.Task, creature *rules.Creature) (rules.Result, error) {
// 	result := rules.Result{}
// 	switch task.Kind {
// 	case types.SkillCheck:
// 		a.logger.Info("skill check", "ability", task.Ability, "skill", task.Skill, "target", task.Target)
// 		check := rules.Check{
// 			Ability: task.Ability,
// 			Skill:   task.Skill,
// 			Target:  task.Target,
// 		}
// 		if task.Ability != "" {
// 			check.Override = task.Ability
// 		}
// 		res, err := creature.SkillCheck(&a.dice, check)
// 		if err != nil {
// 			a.logger.Error("creature error skill check", "error", err.Error())
// 			return result, err
// 		}
// 		result = res

// 	}
// 	return result, nil
// }

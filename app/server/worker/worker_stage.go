package worker

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/betorvs/playbypost/core/initiative"
	"github.com/betorvs/playbypost/core/parser"
	"github.com/betorvs/playbypost/core/rpg/base"
	"github.com/betorvs/playbypost/core/rules"
	"github.com/betorvs/playbypost/core/sys/web/types"
	"github.com/betorvs/playbypost/core/utils"
)

func (a *WorkerAPI) parseCommand(cmd types.Activity) error {
	// call back to slack
	enc, err := a.db.GetStageEncounterByEncounterID(a.ctx, cmd.EncounterID)
	if err != nil {
		a.logger.Error("error getting encounter by id", "error", err.Error())
		return err
	}
	a.logger.Info("encounter found in parseCommand", "encounter", enc)
	processed := false
	// parse command
	switch {
	case cmd.Actions["command"] == parser.CloseStage:
		// close stage
		a.logger.Info("close stage")
		// send message to chat
		closingMessage := fmt.Sprintf("Stage %s closed. Congratulations!", cmd.Actions["display_text"])
		body, err := a.client.PostEvent(cmd.Actions["channel"], "ALL", closingMessage, types.EventEnd)
		if err != nil {
			a.logger.Error("error posting event close stage", "error", err.Error(), "body", string(body))
			return err
		}
		processed = true
	case cmd.Actions["command"] == parser.ChangeEncounterToStarted:
		// change encounter to started
		a.logger.Info("change encounter to started")
		err := a.db.UpdatePhase(a.ctx, cmd.EncounterID, int(types.Started))
		if err != nil {
			a.logger.Error("error updating encounter phase to started", "error", err.Error())
			return err
		}
		encounterAnnounce := fmt.Sprintf("Encounter %s has started: _%s_", enc.Title, enc.Announcement)
		body, err := a.client.PostEvent(cmd.Actions["channel"], "ALL", encounterAnnounce, types.EventAnnounce)
		if err != nil {
			a.logger.Error("error posting event start encounter", "error", err.Error(), "body", string(body))
			return err
		}
		processed = true

	case cmd.Actions["command"] == parser.ChangeEncounterToRunning:
		// change encounter to running
		a.logger.Info("change encounter to running")
		err := a.db.UpdatePhase(a.ctx, cmd.EncounterID, int(types.Running))
		if err != nil {
			a.logger.Error("error updating encounter phase to running", "error", err.Error())
			return err
		}
		// call back to slack
		encounterAnnounce := fmt.Sprintf("Encounter %s is running. Make your actions!", enc.Title)
		body, err := a.client.PostEvent(cmd.Actions["channel"], "ALL", encounterAnnounce, types.EventAnnounce)
		if err != nil {
			a.logger.Error("error posting event running encounter", "error", err.Error(), "body", string(body))
			return err
		}
		processed = true

	case cmd.Actions["command"] == parser.ChangeEncounterToFinished:
		// change encounter to finished
		a.logger.Info("change encounter to finished")
		err := a.db.UpdatePhase(a.ctx, cmd.EncounterID, int(types.Finished))
		if err != nil {
			a.logger.Error("error updating encounter phase to finished", "error", err.Error())
			return err
		}
		// call back to slack
		encounterAnnounce := fmt.Sprintf("Encounter %s finished.", enc.Title)
		body, err := a.client.PostEvent(cmd.Actions["channel"], "ALL", encounterAnnounce, types.EventAnnounce)
		if err != nil {
			a.logger.Error("error posting event finishing encounter", "error", err.Error(), "body", string(body))
			return err
		}
		// get next encounter and start it
		nextEncounter, err := a.db.GetNextEncounterByEncounterID(a.ctx, cmd.EncounterID)
		if err != nil {
			a.logger.Error("error getting next encounter by encounter id", "error", err.Error())
			return err
		}
		// TODO: check future args column
		if nextEncounter.NextEncounterID != 0 {
			// change next encounter to started
			err := a.db.UpdatePhase(a.ctx, nextEncounter.NextEncounterID, int(types.Started))
			if err != nil {
				a.logger.Error("error updating next encounter phase to started", "error", err.Error())
				return err
			}
			encounterToAnnounce, err := a.db.GetStageEncounterByEncounterID(a.ctx, nextEncounter.NextEncounterID)
			if err != nil {
				a.logger.Error("error getting encounter to announce", "error", err.Error())
				return err
			}
			encounterAnnounce := fmt.Sprintf("Encounter %s has started: _%s_", encounterToAnnounce.Title, encounterToAnnounce.Announcement)
			body, err := a.client.PostEvent(cmd.Actions["channel"], "ALL", encounterAnnounce, types.EventAnnounce)
			if err != nil {
				a.logger.Error("error posting event start next encounter", "error", err.Error(), "body", string(body))
				return err
			}
		}

		processed = true

	case cmd.Actions["command"] == parser.RollInitiative:
		a.logger.Info("roll initiative")
		if types.Phase(enc.Phase) != types.Running {
			return fmt.Errorf("encounter not in running phase")
		}
		msg, err := a.rollInitiative(enc)
		if err != nil {
			a.logger.Error("error rolling initiative", "error", err.Error())
			return err
		}
		a.logger.Info("initiative rolled", "msg", msg)
		body, err := a.client.PostEvent(cmd.Actions["channel"], "ALL", msg, types.EventAnnounce)
		if err != nil {
			a.logger.Error("error posting event roll initiative", "error", err.Error(), "body", string(body))
			return err
		}
		cmd.Actions["result"] = msg
		processed = true

	case strings.HasPrefix(cmd.Actions["command"], parser.HealthStatus):
		a.logger.Info("health status")
		players, npcs, err := a.db.GetCreatureFromParticipantsList(a.ctx, enc.PC, enc.NPC, a.rpg)
		if err != nil {
			a.logger.Error("error getting creature from participants list", "error", err.Error())
			return err
		}
		if strings.Contains(cmd.Actions["command"], "npc") {
			npcID, err := strconv.Atoi(cmd.Actions["npc_id"])
			if err != nil {
				a.logger.Error("error getting npc id", "error", err.Error())
				return err
			}
			npc, nOk := npcs[npcID]
			if !nOk {
				a.logger.Error("npc not found", "npc_id", npcID)
				return fmt.Errorf("npc not found")
			}
			// send message to slack
			msg := fmt.Sprintf("NPC %s health status: %d", npc.Name(), npc.HealthStatus())
			body, err := a.client.PostEvent(cmd.Actions["channel"], cmd.Actions["userid"], msg, types.EventAnnounce)
			if err != nil {
				a.logger.Error("error posting event health status", "error", err.Error(), "body", string(body))
				return err
			}
		} else {
			playerID, err := strconv.Atoi(cmd.Actions["player_id"])
			if err != nil {
				a.logger.Error("error getting player id", "error", err.Error())
				return err
			}
			player, pOk := players[playerID]
			if !pOk {
				a.logger.Error("player not found", "player_id", playerID)
				return fmt.Errorf("player not found")
			}
			// send message to slack
			msg := fmt.Sprintf("Player %s health status: %d", player.Name(), player.HealthStatus())
			body, err := a.client.PostEvent(cmd.Actions["channel"], cmd.Actions["userid"], msg, types.EventAnnounce)
			if err != nil {
				a.logger.Error("error posting event health status", "error", err.Error(), "body", string(body))
				return err
			}
		}
		processed = true

	case strings.HasPrefix(cmd.Actions["command"], parser.CurrentInitiative):
		a.logger.Info("current initiative")
		if enc.InitiativeID == 0 {
			a.logger.Error("initiative not rolled")
			return fmt.Errorf("initiative not rolled")
		}
		a.logger.Info("initiative id", "initiative", enc.InitiativeID, "actions_initiative_id", cmd.Actions["initiative_id"])
		initiative, err := a.db.GetInitiativeByID(a.ctx, enc.InitiativeID)
		if err != nil {
			a.logger.Error("error getting initiative by id", "error", err.Error())
			return err
		}
		a.logger.Info("initiative found", "current_name", initiative.Current())
		// send message to slack
		msg := fmt.Sprintf("Current initiative: %s", initiative.Current())
		body, err := a.client.PostEvent(cmd.Actions["channel"], "ALL", msg, types.EventAnnounce)
		if err != nil {
			a.logger.Error("error posting event current initiative", "error", err.Error(), "body", string(body))
			return err
		}
		processed = true

	case strings.HasPrefix(cmd.Actions["command"], parser.AttackNPC):
		a.logger.Info("attack npc")
		// check initiative
		if enc.InitiativeID == 0 {
			a.logger.Error("initiative not rolled")
			return fmt.Errorf("initiative not rolled")
		}
		a.logger.Info("initiative id", "initiative", enc.InitiativeID, "actions_initiative_id", cmd.Actions["initiative_id"])
		initiative, err := a.db.GetInitiativeByID(a.ctx, enc.InitiativeID)
		if err != nil {
			a.logger.Error("error getting initiative by id", "error", err.Error())
			return err
		}
		a.logger.Info("initiative found", "current", initiative.Current())
		// check if first to play is the player
		players, npcs, err := a.db.GetCreatureFromParticipantsList(a.ctx, enc.PC, enc.NPC, a.rpg)
		if err != nil {
			a.logger.Error("error getting creature from participants list", "error", err.Error())
			return err
		}
		playerID, err := strconv.Atoi(cmd.Actions["player_id"])
		if err != nil {
			a.logger.Error("error getting player id", "error", err.Error())
			return err
		}
		npcID, err := strconv.Atoi(cmd.Actions["npc_id"])
		if err != nil {
			a.logger.Error("error getting npc id", "error", err.Error())
			return err
		}
		attacker, aOk := players[playerID]
		if !aOk {
			a.logger.Error("player not found", "player_id", playerID)
			return fmt.Errorf("player not found")
		}
		defensor, dOK := npcs[npcID]
		if !dOK {
			a.logger.Error("npc not found", "npc_id", npcID)
			return fmt.Errorf("npc not found")
		}
		a.logger.Info("npc found", "npc_id", npcID, "npc", defensor)
		a.logger.Info("player found", "player_id", playerID, "player", attacker)
		if initiative.Current() == attacker.Name() {
			a.logger.Info("player ready to play", "player_name", attacker.Name(), "npc_name", defensor.Name())
			attack := rules.NewAttack(cmd.Actions["command"], "longsword", rules.Melee, attacker, defensor, &a.dice, a.logger)
			attack.Call()
			a.logger.Info("attack result", "attack", attack, "defensor_health", attack.Defensor.HealthStatus())
			err := attack.Defensor.Update(a.ctx, npcID, a.db.UpdateNPC)
			if err != nil {
				a.logger.Error("error updating npc", "error", err.Error())
				return err
			}
			emoji := types.EventFailure
			if attack.Response.Success {
				emoji = types.EventSuccess
			}
			attackResult := fmt.Sprintf("Player %s attacked NPC %s. Result: %s", attacker.Name(), defensor.Name(), attack.Response.Text)
			body, err := a.client.PostEvent(cmd.Actions["channel"], "ALL", attackResult, emoji)
			if err != nil {
				a.logger.Error("error posting event attack npc", "error", err.Error(), "body", string(body))
				return err
			}
			cmd.Actions["result"] = attackResult
			// update initiative
			err = a.db.UpdateNextPlayer(a.ctx, enc.InitiativeID, initiative.Next())
			if err != nil {
				a.logger.Error("error updating next player", "error", err.Error())
				return err
			}
			if attack.Defensor.IsDead() {
				_, err = a.db.DeactivateParticipant(a.ctx, enc.InitiativeID, attack.Defensor.Name())
				if err != nil {
					a.logger.Error("error deactivating participant", "error", err.Error())
					return err
				}
				// send message to slack
				msg := fmt.Sprintf("NPC %s is dead", defensor.Name())
				body, err = a.client.PostEvent(cmd.Actions["channel"], "ALL", msg, types.EventDead)
				if err != nil {
					a.logger.Error("error posting event npc dead", "error", err.Error(), "body", string(body))
				}
				initiative.RemoveByName(defensor.Name())
			}
			// send to slack
			msg := fmt.Sprintf("Next participant: %s", initiative.Current())
			body, err = a.client.PostEvent(cmd.Actions["channel"], "ALL", msg, types.EventAnnounce)
			if err != nil {
				a.logger.Error("error posting event next participant", "error", err.Error(), "body", string(body))
			}

			a.logger.Info("next participant", "next", initiative.Current())
			processed = true
		}

	case strings.HasPrefix(cmd.Actions["command"], parser.AttackPlayer):
		a.logger.Info("attack player")
		// check initiative
		if enc.InitiativeID == 0 {
			a.logger.Error("initiative not rolled")
			return fmt.Errorf("initiative not rolled")
		}
		a.logger.Info("initiative id", "initiative", enc.InitiativeID, "actions_initiative_id", cmd.Actions["initiative_id"])
		playerID, err := strconv.Atoi(cmd.Actions["player_id"])
		if err != nil {
			a.logger.Error("error getting player id", "error", err.Error())
			return err
		}
		npcID, err := strconv.Atoi(cmd.Actions["npc_id"])
		if err != nil {
			a.logger.Error("error getting npc id", "error", err.Error())
			return err
		}
		initiative, err := a.db.GetInitiativeByID(a.ctx, enc.InitiativeID)
		if err != nil {
			a.logger.Error("error getting initiative by id", "error", err.Error())
			return err
		}
		a.logger.Info("initiative found", "current", initiative.Current())
		players, npcs, err := a.db.GetCreatureFromParticipantsList(a.ctx, enc.PC, enc.NPC, a.rpg)
		if err != nil {
			a.logger.Error("error getting creature from participants list", "error", err.Error())
			return err
		}
		attacker, aOk := npcs[npcID]
		if !aOk {
			a.logger.Error("npc not found", "npc_id", npcID)
			return fmt.Errorf("npc not found")
		}
		defensor, dOK := players[playerID]
		if !dOK {
			a.logger.Error("player not found", "player_id", playerID)
			return fmt.Errorf("player not found")
		}
		a.logger.Info("npc found", "npc_id", npcID, "npc", attacker)
		a.logger.Info("player found", "player_id", playerID, "player", defensor)
		if initiative.Current() == attacker.Name() {
			a.logger.Info("npc ready to play", "npc_name", attacker.Name(), "player_name", defensor.Name())
			attack := rules.NewAttack(cmd.Actions["command"], "longsword", rules.Melee, attacker, defensor, &a.dice, a.logger)
			attack.Call()
			a.logger.Info("attack result", "attack", attack, "defensor_health", attack.Defensor.HealthStatus())
			err := attack.Defensor.Update(a.ctx, playerID, a.db.UpdatePlayer)
			if err != nil {
				a.logger.Error("error updating player", "error", err.Error())
				return err
			}
			emoji := types.EventFailure
			if attack.Response.Success {
				emoji = types.EventSuccess
			}
			attackResult := fmt.Sprintf("NPC %s attacked Player %s. Result: %s", attacker.Name(), defensor.Name(), attack.Response.Text)
			body, err := a.client.PostEvent(cmd.Actions["channel"], "ALL", attackResult, emoji)
			if err != nil {
				a.logger.Error("error posting event attack player", "error", err.Error(), "body", string(body))
				return err
			}
			cmd.Actions["result"] = attackResult
			// update initiative
			err = a.db.UpdateNextPlayer(a.ctx, enc.InitiativeID, initiative.Next())
			if err != nil {
				a.logger.Error("error updating next player", "error", err.Error())
				return err
			}
			if attack.Defensor.IsDead() {
				_, err = a.db.DeactivateParticipant(a.ctx, enc.InitiativeID, defensor.Name())
				if err != nil {
					a.logger.Error("error deactivating participant", "error", err.Error())
					return err
				}
				// send message to slack
				msg := fmt.Sprintf("Player %s is dead", defensor.Name())
				body, err = a.client.PostEvent(cmd.Actions["channel"], "ALL", msg, types.EventDead)
				if err != nil {
					a.logger.Error("error posting event player dead", "error", err.Error(), "body", string(body))
				}
				initiative.RemoveByName(defensor.Name())
			}

			// send next player to slack
			a.logger.Info("next participant", "next", initiative.Current())
			// send to slack
			msg := fmt.Sprintf("Next participant: %s", initiative.Current())
			body, err = a.client.PostEvent(cmd.Actions["channel"], "ALL", msg, types.EventAnnounce)
			if err != nil {
				a.logger.Error("error posting event next participant", "error", err.Error(), "body", string(body))
			}

			processed = true
		}

	case strings.HasPrefix(cmd.Actions["command"], parser.Task):
		a.logger.Info("task command")
		player, err := a.db.GetPlayerByUserIDChannel(a.ctx, cmd.Actions["userid"], cmd.Actions["channel"], a.rpg)
		if err != nil {
			a.logger.Error("error getting player by user id", "error", err.Error())
			return err
		}
		a.logger.Info("player found", "player", player)
		a.logger.Info("rpg system", "rpg", a.rpg.BaseSystem)
		creature := base.RestoreCreature()
		creature.RPG = a.rpg
		character := types.PlayerToCreature(&player, creature, a.lib)
		a.logger.Info("creature found", "creature", creature)
		// get task
		taskID, err := strconv.Atoi(cmd.Actions["task_id"])
		if err != nil {
			a.logger.Error("error getting task id", "error", err.Error())
			return err
		}
		task, err := a.db.GetStageTaskFromRunningTaskID(a.ctx, taskID)
		if err != nil {
			a.logger.Error("error getting task id", "error", err.Error())
			return err
		}
		a.logger.Info("task found", "task", task)

		result, err := a.executeTask(task, character)
		if err != nil {
			a.logger.Error("error executing task", "error", err.Error())
			return err
		}
		// call back to slack
		eventKind := types.EventFailure
		if result.Success {
			eventKind = types.EventSuccess
		}
		resultMessage := fmt.Sprintf("Result: %d with dices rolled: %s", result.Result, result.Rolled)
		taskResult := fmt.Sprintf("Task %s executed by @%s. %s", task.Description, cmd.Actions["userid"], resultMessage)
		body, err := a.client.PostEvent(cmd.Actions["channel"], "ALL", taskResult, eventKind)
		if err != nil {
			a.logger.Error("error posting event task result", "error", err.Error(), "body", string(body))
			return err
		}
		cmd.Actions["result"] = resultMessage
		processed = true

	default:
		a.logger.Info("command not defined")
	}
	// update activity processed
	if processed {
		err = a.db.UpdateProcessedActivities(a.ctx, cmd.ID, true, cmd.Actions)
		if err != nil {
			a.logger.Error("error updating stage encounter activities", "error", err.Error())
			return err
		}
	}

	return nil
}

func (a *WorkerAPI) executeTask(task types.Task, creature rules.RolePlaying) (rules.Result, error) {
	result := rules.Result{}
	switch task.Kind {
	case types.SkillCheck:
		a.logger.Info("skill check", "ability", task.Ability, "skill", task.Skill, "target", task.Target)
		check := rules.Check{
			Ability: task.Ability,
			Skill:   task.Skill,
			Target:  task.Target,
		}
		if task.Ability != "" {
			check.Override = task.Ability
		}
		res, err := creature.SkillCheck(&a.dice, check, a.logger, a.lib)
		if err != nil {
			a.logger.Error("creature error skill check", "error", err.Error())
			return result, err
		}
		result = res

	}
	return result, nil
}

func (a *WorkerAPI) rollInitiative(enc types.StageEncounter) (string, error) {

	party := make(map[string]int)
	players, npcs, err := a.db.GetCreatureFromParticipantsList(a.ctx, enc.PC, enc.NPC, a.rpg)
	if err != nil {
		return "", err
	}

	for _, p := range players {
		a.logger.Info("players found", "players", p)
		i, _ := p.InitiativeBonus()
		party[p.Name()] = i
		a.logger.Info("participant found roll", "name", p.Name, "init bonus", i)
	}
	for _, n := range npcs {
		a.logger.Info("npcs found roll", "npcs", n)
		i, _ := n.InitiativeBonus()
		party[n.Name()] = i
		a.logger.Info("npc found", "name", n.Name, "init bonus", i)
	}
	randomInit := utils.RandomString(6)
	name := fmt.Sprintf("init-%s-encID-%d", randomInit, enc.ID)
	init := initiative.NewInitiative(a.dice, party, name, a.rpg.InitiativeDice())
	a.logger.Info("initiative rolled", "initiative", fmt.Sprintf("%+v", init))
	a.logger.Info("participants", "participants", init.Participants)
	initID, err := a.db.SaveInitiativeTx(a.ctx, init, enc.ID)
	if err != nil {
		return "", err
	}
	msg := fmt.Sprintf("initiative id %d, and first to play %s", initID, init.Participants[0].Name)
	// msg := fmt.Sprintf("initiative id %s, and first to play %s", name, init.Participants[0].Name)
	return msg, nil
	// return name, nil
}

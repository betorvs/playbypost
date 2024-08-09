package pg

import (
	"context"
	"database/sql"
	"slices"

	"github.com/betorvs/playbypost/core/rpg"
	"github.com/betorvs/playbypost/core/rules"
	"github.com/betorvs/playbypost/core/sys/web/types"
	"github.com/betorvs/playbypost/core/utils"
)

func (db *DBX) GetStage(ctx context.Context) ([]types.Stage, error) {
	stages := []types.Stage{}
	query := "SELECT id, display_text, story_id, storyteller_id FROM stage"
	rows, err := db.Conn.QueryContext(ctx, query)
	if err != nil {
		db.Logger.Error("query on stage failed", "error", err.Error())
		return stages, err
	}
	defer rows.Close()
	for rows.Next() {
		var s types.Stage
		if err := rows.Scan(&s.ID, &s.Text, &s.StoryID, &s.StorytellerID); err != nil {
			db.Logger.Error("scan error on stage", "error", err.Error())
		}
		stages = append(stages, s)
	}
	// Check for errors from iterating over rows.
	if err := rows.Err(); err != nil {
		db.Logger.Error("rows err on stage", "error", err.Error())
	}
	return stages, nil
}

func (db *DBX) GetStageByStoryID(ctx context.Context, id int) ([]types.Stage, error) {
	stages := []types.Stage{}
	query := "SELECT id, display_text, story_id, storyteller_id FROM stage WHERE story_id = $1"
	rows, err := db.Conn.QueryContext(ctx, query, id)
	if err != nil {
		db.Logger.Error("query on stage by story_id failed", "error", err.Error())
		return stages, err
	}
	defer rows.Close()
	for rows.Next() {
		var s types.Stage
		if err := rows.Scan(&s.ID, &s.Text, &s.StoryID, &s.StorytellerID); err != nil {
			db.Logger.Error("scan error on stage by story_id", "error", err.Error())
		}
		stages = append(stages, s)
	}
	// Check for errors from iterating over rows.
	if err := rows.Err(); err != nil {
		db.Logger.Error("rows err on stage by story_id", "error", err.Error())
	}
	return stages, nil
}

func (db *DBX) GetStageByStageID(ctx context.Context, id int) (types.StageAggregated, error) {
	aggr := types.StageAggregated{}
	query := "SELECT sa.id, sa.display_text, sa.story_id, sa.storyteller_id, sa.encoding_key, sy.title, sy.announcement, sy.notes, sy.writer_id, u.userid, sc.channel, sc.active FROM stage AS sa JOIN story AS sy ON sa.story_id = sy.id JOIN users AS u ON sa.storyteller_id = u.id LEFT JOIN stage_channel AS sc ON sc.stage_id = sa.id WHERE sa.id = $1"
	rows, err := db.Conn.QueryContext(ctx, query, id)
	if err != nil {
		db.Logger.Error("query on stage failed", "error", err.Error())
		return aggr, err
	}
	defer rows.Close()
	for rows.Next() {
		var sa types.Stage
		var encodingKey, announce, notes string
		var story types.Story
		var channel sql.NullString
		var channelActive sql.NullBool
		if err := rows.Scan(&sa.ID, &sa.Text, &sa.StoryID, &sa.StorytellerID, &encodingKey, &story.Title, &announce, &notes, &story.WriterID, &sa.UserID, &channel, &channelActive); err != nil {
			db.Logger.Error("scan error on stage", "error", err.Error())
		}
		story.Announcement, _ = utils.DecryptText(announce, encodingKey)
		story.Notes, _ = utils.DecryptText(notes, encodingKey)
		aggr.Stage = sa
		aggr.Story = story
		if channel.Valid && channelActive.Valid {
			c := types.Channel{}
			c.Active = channelActive.Bool
			c.Channel = channel.String
			aggr.Channel = c
		}
	}
	// Check for errors from iterating over rows.
	if err := rows.Err(); err != nil {
		db.Logger.Error("rows err on stage", "error", err.Error())
	}
	return aggr, nil
}

func (db *DBX) GetStageEncounterByEncounterID(ctx context.Context, id int) (types.StageEncounter, error) {
	// enc := types.StageEncounter{}
	// query := "select se.ID, se.display_text, e.title, se.storyteller_id, se.phase, e.notes, e.announcement, s.encoding_key from stage_encounters AS se JOIN encounters AS e ON se.encounters_id = e.id JOIN stage AS s ON se.stage_id = s.id WHERE se.ID = $1"
	// rows, err := db.Conn.QueryContext(ctx, query, id)
	// if err != nil {
	// 	db.Logger.Error("query on stage_encounters by stage_encounters.ID failed", "error", err.Error())
	// 	return enc, err
	// }
	// defer rows.Close()
	// for rows.Next() {
	// 	var s types.StageEncounter
	// 	var notes, announce, encodingKey string
	// 	if err := rows.Scan(&s.ID, &s.Text, &s.Title, &s.StorytellerID, &s.Phase, &notes, &announce, &encodingKey); err != nil {
	// 		db.Logger.Error("scan error on stage_encounters by stage_encounters.ID ", "error", err.Error())
	// 	}
	// 	s.Notes, _ = utils.DecryptText(notes, encodingKey)
	// 	s.Announcement, _ = utils.DecryptText(announce, encodingKey)
	// 	enc = s
	// }
	// // Check for errors from iterating over rows.
	// if err := rows.Err(); err != nil {
	// 	db.Logger.Error("rows err on stage_encounters by stage_encounters.ID", "error", err.Error())
	// }
	// // participants
	// p, n, err := db.getParticipantsByStageEncounterID(ctx, enc.ID)
	// if err != nil {
	// 	db.Logger.Error("error on GetStageEncounterByEncounterID when calling getParticipantsByStageEncounterID", "error", err.Error())
	// }
	// enc.Players = p
	// enc.NPC = n
	return db.getStageEncounterByEncounterID(ctx, id, -1)
}

func (db *DBX) GetStageEncountersByStageID(ctx context.Context, id int) ([]types.StageEncounter, error) {
	list := []types.StageEncounter{}
	query := "select se.ID, se.display_text, e.title, se.storyteller_id, e.notes, e.announcement, e.writer_id, e.story_id, s.encoding_key, se.updated_at, se.phase from stage_encounters AS se JOIN encounters AS e ON se.encounters_id = e.id JOIN stage AS s ON se.stage_id = s.id WHERE s.id = $1"
	rows, err := db.Conn.QueryContext(ctx, query, id)
	if err != nil {
		db.Logger.Error("query on stage_encounters by stage_id failed", "error", err.Error())
		return list, err
	}
	defer rows.Close()
	for rows.Next() {
		var s types.StageEncounter
		var notes, announce, encodingKey string
		var updatedAt sql.NullTime
		if err := rows.Scan(&s.ID, &s.Text, &s.Title, &s.StorytellerID, &notes, &announce, &s.WriterID, &s.StoryID, &encodingKey, &updatedAt, &s.Phase); err != nil {
			db.Logger.Error("scan error on encounstage_encounters by stage_id ", "error", err.Error(), "updated_at", updatedAt)
		}
		s.Notes, _ = utils.DecryptText(notes, encodingKey)
		s.Announcement, _ = utils.DecryptText(announce, encodingKey)
		if s.Phase != int(types.Finished) {
			list = append(list, s)
		}

	}
	// Check for errors from iterating over rows.
	if err := rows.Err(); err != nil {
		db.Logger.Error("rows err on stage_encounters by stage_id", "error", err.Error())
	}
	return list, nil
}

func (db *DBX) GetRunningStageByChannelID(ctx context.Context, channelID, userID string) (types.RunningStage, error) {
	db.Logger.Info("GetRunningStageByChannelID", "channelID", channelID, "userID", userID)
	running := types.RunningStage{}
	aggr := types.StageAggregated{}
	query := "SELECT sa.id, sa.display_text, sa.story_id, sa.storyteller_id, sa.encoding_key, sy.title, sy.announcement, sy.notes, sy.writer_id, u.userid, sc.channel, sc.active FROM stage AS sa JOIN story AS sy ON sa.story_id = sy.id JOIN users AS u ON sa.storyteller_id = u.id LEFT JOIN stage_channel AS sc ON sc.stage_id = sa.id WHERE sc.channel = $1"
	rows, err := db.Conn.QueryContext(ctx, query, channelID)
	if err != nil {
		db.Logger.Error("query on stage failed", "error", err.Error())
		return running, err
	}
	defer rows.Close()
	for rows.Next() {
		var sa types.Stage
		var encodingKey, announce, notes string
		var story types.Story
		var channel sql.NullString
		var channelActive sql.NullBool
		if err := rows.Scan(&sa.ID, &sa.Text, &sa.StoryID, &sa.StorytellerID, &encodingKey, &story.Title, &announce, &notes, &story.WriterID, &sa.UserID, &channel, &channelActive); err != nil {
			db.Logger.Error("scan error on stage", "error", err.Error())
		}
		story.Announcement, _ = utils.DecryptText(announce, encodingKey)
		story.Notes, _ = utils.DecryptText(notes, encodingKey)
		aggr.Stage = sa
		aggr.Story = story
		if channel.Valid && channelActive.Valid {
			c := types.Channel{}
			c.Active = channelActive.Bool
			c.Channel = channel.String
			aggr.Channel = c
		}
		running.StageAggregated = aggr
	}
	// Check for errors from iterating over rows.
	if err := rows.Err(); err != nil {
		db.Logger.Error("rows err on stage", "error", err.Error())
	}
	storyteller := false
	// get running encounter
	enc, err := db.getStageEncounterByEncounterID(ctx, -1, int(types.Running))
	if err != nil {
		db.Logger.Error("rows err on getStageEncounterByEncounterID", "error", err.Error())
	}
	running.Encounter = enc
	if running.StageAggregated.Stage.UserID == userID {
		storyteller = true
		if enc.ID == 0 {
			db.Logger.Info("running encounter not found", "encounter", enc)
			encs, err := db.GetStageEncountersByStageID(ctx, running.StageAggregated.Stage.ID)
			if err != nil {
				db.Logger.Error("rows err on stage_encounters by stage_id", "error", err.Error())
			}
			db.Logger.Info("encounters list", "encounters", encs)
			running.Encounters = encs
		}

	}

	// 	encOptions := []types.GenericIDName{}
	// 	p := types.PhaseAtoi(enc.Phase)
	// 	db.Logger.Info("phase", "phase", p)
	// 	count := 1
	// 	encOptions = append(encOptions, types.GenericIDName{ID: count, Name: fmt.Sprintf("change-encounter-to-%s", p.NextPhase().String())})
	// 	if len(enc.NPC) > 0 {
	// 		encOptions = append(encOptions, types.GenericIDName{ID: count + 1, Name: "roll-initiative"})
	// 		for _, v := range enc.NPC {
	// 			count++
	// 			encOptions = append(encOptions, types.GenericIDName{ID: count, Name: fmt.Sprintf("act-as-npc-%s", v.Name)})
	// 		}
	// 	}
	// 	running.Options = encOptions
	// }
	// players
	if !storyteller {
		players, err := db.GetPlayerByUserID(ctx, userID, channelID)
		if err != nil {
			db.Logger.Error("rows err on players", "error", err.Error())
		}
		running.Players = players
		// options
		// user options should be tasks and combat options
		// storyteller should be encounter phases and npc actions and start next encounter

		options, err := db.getRunningTaskByEncounterID(ctx, enc.ID)
		if err != nil {
			db.Logger.Error("rows err on getRunningTaskByEncounterID", "error", err.Error())
		}
		running.Encounter.Options = options
		// count := len(options)
		// if count > 0 {
		// 	for _, v := range enc.NPC {
		// 		count++
		// 		running.Options = append(running.Options, types.GenericIDName{ID: count, Name: fmt.Sprintf("attack-npc-%s", v.Name)})
		// 	}
		// }
	}

	return running, nil

}

// stage_running_tasks
func (db *DBX) getRunningTaskByEncounterID(ctx context.Context, id int) ([]types.GenericIDName, error) {
	tasks := []types.GenericIDName{}
	query := "select sa.display_text, sa.id from stage_running_tasks AS sa WHERE sa.stage_encounters_id = $1"
	rows, err := db.Conn.QueryContext(ctx, query, id)
	if err != nil {
		db.Logger.Error("query on stage_running_tasks by encounter_id failed", "error", err.Error())
		return tasks, err
	}
	defer rows.Close()
	for rows.Next() {
		var s types.GenericIDName
		if err := rows.Scan(&s.Name, &s.ID); err != nil {
			db.Logger.Error("scan error on stage_running_tasks by encounter_id ", "error", err.Error())
		}
		tasks = append(tasks, s)
	}
	// Check for errors from iterating over rows.
	if err := rows.Err(); err != nil {
		db.Logger.Error("rows err on stage_running_tasks by encounter_id", "error", err.Error())
	}
	return tasks, nil
}

func (db *DBX) getStageEncounterByEncounterID(ctx context.Context, id int, phase int) (types.StageEncounter, error) {
	enc := types.StageEncounter{}
	var r *sql.Rows
	switch {
	case phase == -1 && id != 0:
		//
		query := "select se.ID, se.display_text, e.title, se.storyteller_id, se.phase, e.notes, e.announcement, e.writer_id, e.story_id, s.encoding_key, i.id from stage_encounters AS se JOIN encounters AS e ON se.encounters_id = e.id JOIN stage AS s ON se.stage_id = s.id LEFT JOIN initiative AS i ON i.stage_encounters_id = se.id WHERE se.ID = $1"
		rows, err := db.Conn.QueryContext(ctx, query, id)
		if err != nil {
			db.Logger.Error("query on stage_encounters by stage_encounters.ID failed", "error", err.Error())
			return enc, err
		}
		defer rows.Close()
		r = rows
	case phase > 0 && id == -1:
		// can return a initiative id
		query := "select se.ID, se.display_text, e.title, se.storyteller_id, se.phase, e.notes, e.announcement, e.writer_id, e.story_id, s.encoding_key, i.id from stage_encounters AS se JOIN encounters AS e ON se.encounters_id = e.id JOIN stage AS s ON se.stage_id = s.id LEFT JOIN initiative AS i ON i.stage_encounters_id = se.id WHERE se.phase = $1"
		rows, err := db.Conn.QueryContext(ctx, query, phase)
		if err != nil {
			db.Logger.Error("query on stage_encounters by phase equal 3 (running) failed", "error", err.Error())
			return enc, err
		}
		defer rows.Close()
		r = rows
	}

	defer r.Close()
	for r.Next() {
		var s types.StageEncounter
		var notes, announce, encodingKey string
		var initiativeID sql.NullInt64
		if err := r.Scan(&s.ID, &s.Text, &s.Title, &s.StorytellerID, &s.Phase, &notes, &announce, &s.WriterID, &s.StoryID, &encodingKey, &initiativeID); err != nil {
			db.Logger.Error("scan error on stage_encounters by stage_encounters.ID ", "error", err.Error())
		}
		s.Notes, _ = utils.DecryptText(notes, encodingKey)
		s.Announcement, _ = utils.DecryptText(announce, encodingKey)
		if initiativeID.Valid {
			s.InitiativeID = int(initiativeID.Int64)
		}
		enc = s
	}
	// Check for errors from iterating over rows.
	if err := r.Err(); err != nil {
		db.Logger.Error("rows err on stage_encounters by stage_encounters.ID", "error", err.Error())
	}
	// participants
	p, n, err := db.getParticipantsByStageEncounterID(ctx, enc.ID)
	if err != nil {
		db.Logger.Error("error on db.getStageEncounterByEncounterID when calling getParticipantsByStageEncounterID", "error", err.Error())
	}
	db.Logger.Info("players and npcs", "encounter_id", enc.ID, "player", p, "npc", n)
	enc.PC = p
	enc.NPC = n
	return enc, nil
}

func (db *DBX) getParticipantsByStageEncounterID(ctx context.Context, id int) ([]types.GenericIDName, []types.GenericIDName, error) {
	players := []types.GenericIDName{}
	npcs := []types.GenericIDName{}
	query := "select sp.players_id, pl.character_name, pl.destroyed, snp.non_players_id, npc.npc_name, npc.destroyed from stage_encounters AS se LEFT JOIN stage_encounters_participants_players AS sp ON sp.stage_encounters_id = se.id LEFT JOIN players AS pl ON pl.id = sp.players_id LEFT JOIN stage_encounters_participants_non_players AS snp ON snp.stage_encounters_id = se.id LEFT JOIN non_players AS npc ON npc.id = snp.non_players_id WHERE se.ID = $1"
	rows, err := db.Conn.QueryContext(ctx, query, id)
	if err != nil {
		db.Logger.Error("query on stage_encounters_participants_players and stage_encounters_participants_non_players failed", "error", err.Error())
		return players, npcs, err
	}
	defer rows.Close()
	var p types.GenericIDName
	var n types.GenericIDName
	for rows.Next() {
		var pcID, npcID sql.NullInt64
		var pcName, npcName sql.NullString
		var pcDestroyed, npcDestroyed sql.NullBool
		if err := rows.Scan(&pcID, &pcName, &pcDestroyed, &npcID, &npcName, &npcDestroyed); err != nil {
			db.Logger.Error("scan error on stage_encounters_participants_players and stage_encounters_participants_non_players", "error", err.Error())
		}
		if pcID.Valid && pcName.Valid && pcDestroyed.Valid && !pcDestroyed.Bool {
			p.ID = int(pcID.Int64)
			p.Name = pcName.String
			if !slices.Contains(players, p) {
				players = append(players, p)
			}
		}
		if npcID.Valid && npcName.Valid && npcDestroyed.Valid && !npcDestroyed.Bool {
			n.ID = int(npcID.Int64)
			n.Name = npcName.String
			if !slices.Contains(npcs, n) {
				npcs = append(npcs, n)
			}
		}
	}
	// Check for errors from iterating over rows.
	if err := rows.Err(); err != nil {
		db.Logger.Error("rows err on stage_encounters_participants_players and stage_encounters_participants_non_players", "error", err.Error())
	}
	return players, npcs, err
}

// stage_encounter_activities
func (db *DBX) GetStageEncounterActivitiesByEncounterID(ctx context.Context, id int) ([]types.StageEncounterActivities, error) {
	list := []types.StageEncounterActivities{}
	query := "select sa.id, sa.actions, sa.stage_id, sa.encounter_id, sa.processed from stage_encounter_activities AS sa WHERE sa.encounter_id = $1"
	rows, err := db.Conn.QueryContext(ctx, query, id)
	if err != nil {
		db.Logger.Error("query on stage_encounter_activities by encounter_id failed", "error", err.Error())
		return list, err
	}
	defer rows.Close()
	for rows.Next() {
		var s types.StageEncounterActivities
		if err := rows.Scan(&s.ID, &s.Actions, &s.StageID, &s.EncounterID, &s.Processed); err != nil {
			db.Logger.Error("scan error on stage_encounter_activities by encounter_id ", "error", err.Error())
		}
		list = append(list, s)
	}
	// Check for errors from iterating over rows.
	if err := rows.Err(); err != nil {
		db.Logger.Error("rows err on stage_encounter_activities by encounter_id", "error", err.Error())
	}
	return list, nil
}

// stage_encounter_activities
func (db *DBX) GetStageEncounterActivities(ctx context.Context) ([]types.StageEncounterActivities, error) {
	list := []types.StageEncounterActivities{}
	query := "select id, stage_id, encounter_id, actions, processed from stage_encounter_activities"
	rows, err := db.Conn.QueryContext(ctx, query)
	if err != nil {
		db.Logger.Error("query on stage_encounter_activities failed", "error", err.Error())
		return list, err
	}
	defer rows.Close()
	for rows.Next() {
		var s types.StageEncounterActivities
		if err := rows.Scan(&s.ID, &s.StageID, &s.EncounterID, &s.Actions, &s.Processed); err != nil {
			db.Logger.Error("scan error on stage_encounter_activities ", "error", err.Error())
		}
		list = append(list, s)
	}
	// Check for errors from iterating over rows.
	if err := rows.Err(); err != nil {
		db.Logger.Error("rows err on stage_encounter_activities", "error", err.Error())
	}
	return list, nil
}

// select t.description, t.kind, t.ability, t.skill FROM tasks AS t JOIN stage_running_tasks AS s ON s.task_id = t.id WHERE s.task_id = 1
func (db *DBX) GetStageTaskFromRunningTaskID(ctx context.Context, taskID int) (types.Task, error) {
	t := types.Task{}
	query := "select t.description, t.kind, t.ability, t.skill, t.target FROM tasks AS t JOIN stage_running_tasks AS s ON s.task_id = t.id WHERE s.task_id = $1"
	rows, err := db.Conn.QueryContext(ctx, query, taskID)
	if err != nil {
		db.Logger.Error("query on tasks by task_id failed", "error", err.Error())
		return t, err
	}
	defer rows.Close()
	for rows.Next() {
		if err := rows.Scan(&t.Description, &t.Kind, &t.Ability, &t.Skill, &t.Target); err != nil {
			db.Logger.Error("scan error on tasks by task_id ", "error", err.Error())
		}
	}
	// Check for errors from iterating over rows.
	if err := rows.Err(); err != nil {
		db.Logger.Error("rows err on tasks by task_id", "error", err.Error())
	}
	return t, nil
}

func (db *DBX) GetCreatureFromParticipantsList(ctx context.Context, players []types.GenericIDName, npcs []types.GenericIDName, rpgSystem *rpg.RPGSystem) (map[int]*rules.Creature, map[int]*rules.Creature, error) {
	db.Logger.Info("GetCreatureFromParticipantsList", "players", players, "npcs", npcs)
	// players
	playersMap := map[int]*rules.Creature{}
	creatureMap := map[int]*rules.Creature{}
	{
		query := "SELECT id, character_name, abilities, skills, extensions FROM players"
		rows, err := db.Conn.QueryContext(ctx, query)
		if err != nil {
			db.Logger.Error("query on players failed", "error", err.Error())
			return map[int]*rules.Creature{}, map[int]*rules.Creature{}, err
		}
		defer rows.Close()
		for rows.Next() {
			c := rules.RestoreCreature()
			c.RPG = rpgSystem
			var id int
			extended := rpg.NewExtended()
			if err := rows.Scan(&id, &c.Name, &c.Abilities, &c.Skills, &extended); err != nil {
				db.Logger.Error("scan error on players by id ", "error", err.Error())
			}
			c.Extension = rpg.NewExtendedSystem(rpgSystem, extended)
			creatureMap[id] = c
		}
	}
	for _, v := range players {
		if c, ok := creatureMap[v.ID]; ok {
			playersMap[v.ID] = c
		}
	}
	// npcs
	npcsMap := map[int]*rules.Creature{}
	creatureMap2 := map[int]*rules.Creature{}
	{
		query := "SELECT id, npc_name, abilities, skills, extensions FROM non_players"
		rows, err := db.Conn.QueryContext(ctx, query)
		if err != nil {
			db.Logger.Error("query on non_players by id failed", "error", err.Error())
			return map[int]*rules.Creature{}, map[int]*rules.Creature{}, err
		}
		defer rows.Close()
		for rows.Next() {
			c := rules.RestoreCreature()
			c.RPG = rpgSystem
			var id int
			extended := rpg.NewExtended()
			if err := rows.Scan(&id, &c.Name, &c.Abilities, &c.Skills, &extended); err != nil {
				db.Logger.Error("scan error on non_players by id ", "error", err.Error())
			}
			c.Extension = rpg.NewExtendedSystem(rpgSystem, extended)
			creatureMap2[id] = c
		}
	}
	for _, v := range npcs {
		if c, ok := creatureMap2[v.ID]; ok {
			npcsMap[v.ID] = c
		}
	}

	return playersMap, npcsMap, nil
}

// stage_next_encounter
func (db *DBX) GetNextEncounterByEncounterID(ctx context.Context, id int) (types.NextEncounter, error) {
	ne := types.NextEncounter{}
	query := "select display_text, stage_id, current_encounter_id, next_encounter_id FROM stage_next_encounter WHERE current_encounter_id = $1"
	rows, err := db.Conn.QueryContext(ctx, query, id)
	if err != nil {
		db.Logger.Error("query on stage_next_encounter by encounter_id failed", "error", err.Error())
		return ne, err
	}
	defer rows.Close()
	for rows.Next() {
		if err := rows.Scan(&ne.Text, &ne.StageID, &ne.EncounterID, &ne.NextEncounterID); err != nil {
			db.Logger.Error("scan error on stage_next_encounter by encounter_id ", "error", err.Error())
		}
	}
	// Check for errors from iterating over rows.
	if err := rows.Err(); err != nil {
		db.Logger.Error("rows err on stage_next_encounter by encounter_id", "error", err.Error())
	}
	return ne, nil
}

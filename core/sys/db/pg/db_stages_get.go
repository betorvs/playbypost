package pg

import (
	"context"
	"database/sql"
	"slices"

	"github.com/betorvs/playbypost/core/rpg"
	"github.com/betorvs/playbypost/core/rules"
	"github.com/betorvs/playbypost/core/sys/web/types"
	"github.com/betorvs/playbypost/core/utils"
	"github.com/lib/pq"
)

func (db *DBX) GetStage(ctx context.Context) ([]types.Stage, error) {
	stages := []types.Stage{}
	query := "SELECT id, display_text, story_id, creator_id, storyteller_id FROM stage WHERE finished = false" // dev:finder+query
	rows, err := db.Conn.QueryContext(ctx, query)
	if err != nil {
		db.Logger.Error("query on stage failed", "error", err.Error())
		return stages, err
	}
	defer rows.Close()
	for rows.Next() {
		var s types.Stage
		if err := rows.Scan(&s.ID, &s.Text, &s.StoryID, &s.CreatorID, &s.StorytellerID); err != nil {
			db.Logger.Error("scan error on stage", "error", err.Error())
		}
		stages = append(stages, s)
	}
	// Check for errors FROM iterating over rows.
	if err := rows.Err(); err != nil {
		db.Logger.Error("rows err on stage", "error", err.Error())
	}
	return stages, nil
}

func (db *DBX) GetStageByStoryID(ctx context.Context, id int) ([]types.Stage, error) {
	stages := []types.Stage{}
	query := "SELECT id, display_text, story_id, creator_id, storyteller_id FROM stage WHERE story_id = $1 AND finished = false" // dev:finder+query
	rows, err := db.Conn.QueryContext(ctx, query, id)
	if err != nil {
		db.Logger.Error("query on stage by story_id failed", "error", err.Error())
		return stages, err
	}
	defer rows.Close()
	for rows.Next() {
		var s types.Stage
		if err := rows.Scan(&s.ID, &s.Text, &s.StoryID, &s.CreatorID, &s.StorytellerID); err != nil {
			db.Logger.Error("scan error on stage by story_id", "error", err.Error())
		}
		stages = append(stages, s)
	}
	// Check for errors FROM iterating over rows.
	if err := rows.Err(); err != nil {
		db.Logger.Error("rows err on stage by story_id", "error", err.Error())
	}
	return stages, nil
}

func (db *DBX) GetStageByStageID(ctx context.Context, id int) (types.StageAggregated, error) {
	aggr := types.StageAggregated{}
	query := "SELECT sa.id, sa.display_text, sa.story_id, sa.creator_id, sa.storyteller_id, sa.encoding_key, sy.title, sy.announcement, sy.notes, sy.writer_id, u.userid, sc.channel, sc.active FROM stage AS sa JOIN story AS sy ON sa.story_id = sy.id JOIN users AS u ON sa.storyteller_id = u.id LEFT JOIN stage_channel AS sc ON sc.upstream_id = sa.id WHERE sa.id = $1 AND sa.finished = false" // dev:finder+query
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
		if err := rows.Scan(&sa.ID, &sa.Text, &sa.StoryID, &sa.CreatorID, &sa.StorytellerID, &encodingKey, &story.Title, &announce, &notes, &story.WriterID, &sa.UserID, &channel, &channelActive); err != nil {
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
	// Check for errors FROM iterating over rows.
	if err := rows.Err(); err != nil {
		db.Logger.Error("rows err on stage", "error", err.Error())
	}
	return aggr, nil
}

func (db *DBX) GetStageEncounterByEncounterID(ctx context.Context, id int) (types.StageEncounter, error) {
	return db.getStageEncounterByEncounterID(ctx, id, -1)
}

func (db *DBX) GetStageEncountersByStageID(ctx context.Context, id int) ([]types.StageEncounter, error) {
	list := []types.StageEncounter{}
	query := "SELECT se.ID, se.display_text, e.title, se.storyteller_id, e.notes, e.announcement, e.writer_id, e.story_id, s.encoding_key, se.updated_at, se.phase FROM stage_encounters AS se JOIN encounters AS e ON se.encounter_id = e.id JOIN stage AS s ON se.stage_id = s.id WHERE s.id = $1 AND s.finished = false" // dev:finder+query
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
			db.Logger.Error("scan error on stage_encounters by stage_id ", "error", err.Error(), "updated_at", updatedAt)
		}
		s.Notes, _ = utils.DecryptText(notes, encodingKey)
		s.Announcement, _ = utils.DecryptText(announce, encodingKey)
		if s.Phase != int(types.Finished) {
			list = append(list, s)
		}

	}
	// Check for errors FROM iterating over rows.
	if err := rows.Err(); err != nil {
		db.Logger.Error("rows err on stage_encounters by stage_id", "error", err.Error())
	}
	return list, nil
}

func (db *DBX) GetRunningStageByChannelID(ctx context.Context, channelID, userID string, rpgSystem *rpg.RPGSystem) (types.RunningStage, error) {
	db.Logger.Debug("GetRunningStageByChannelID", "channelID", channelID, "userID", userID)
	running := types.RunningStage{}
	aggr := types.StageAggregated{}
	query := "SELECT sa.id, sa.display_text, sa.story_id, sa.storyteller_id, sa.encoding_key, sy.title, sy.announcement, sy.notes, sy.writer_id, u.userid, sc.channel, sc.active FROM stage AS sa JOIN story AS sy ON sa.story_id = sy.id JOIN users AS u ON sa.storyteller_id = u.id LEFT JOIN stage_channel AS sc ON sc.upstream_id = sa.id WHERE sc.channel = $1 AND sa.finished = false" // dev:finder+query
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
	// Check for errors FROM iterating over rows.
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
			db.Logger.Debug("running encounter not found", "encounter", enc)
			encs, err := db.GetStageEncountersByStageID(ctx, running.StageAggregated.Stage.ID)
			if err != nil {
				db.Logger.Error("rows err on stage_encounters by stage_id", "error", err.Error())
			}
			db.Logger.Debug("encounters list", "encounters", encs)
			running.Encounters = encs
		}

	}

	// players
	if !storyteller {
		players, err := db.GetPlayerByUserID(ctx, userID, channelID, rpgSystem)
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
	}

	return running, nil

}

// stage_encounter_activities
func (db *DBX) GetStageEncounterActivitiesByEncounterID(ctx context.Context, id int) ([]types.Activity, error) {
	list := []types.Activity{}
	query := "SELECT sa.id, sa.actions, sa.upstream_id, sa.encounter_id, sa.processed FROM stage_encounter_activities AS sa WHERE sa.encounter_id = $1" // dev:finder+query
	rows, err := db.Conn.QueryContext(ctx, query, id)
	if err != nil {
		db.Logger.Error("query on stage_encounter_activities by encounter_id failed", "error", err.Error())
		return list, err
	}
	defer rows.Close()
	for rows.Next() {
		var s types.Activity
		if err := rows.Scan(&s.ID, &s.Actions, &s.UpstreamID, &s.EncounterID, &s.Processed); err != nil {
			db.Logger.Error("scan error on stage_encounter_activities by encounter_id ", "error", err.Error())
		}
		list = append(list, s)
	}
	// Check for errors FROM iterating over rows.
	if err := rows.Err(); err != nil {
		db.Logger.Error("rows err on stage_encounter_activities by encounter_id", "error", err.Error())
	}
	return list, nil
}

// stage_encounter_activities
func (db *DBX) GetStageEncounterActivities(ctx context.Context) ([]types.Activity, error) {
	list := []types.Activity{}
	query := "SELECT id, upstream_id, encounter_id, actions, processed FROM stage_encounter_activities" // dev:finder+query
	rows, err := db.Conn.QueryContext(ctx, query)
	if err != nil {
		db.Logger.Error("query on stage_encounter_activities failed", "error", err.Error())
		return list, err
	}
	defer rows.Close()
	for rows.Next() {
		var s types.Activity
		if err := rows.Scan(&s.ID, &s.UpstreamID, &s.EncounterID, &s.Actions, &s.Processed); err != nil {
			db.Logger.Error("scan error on stage_encounter_activities ", "error", err.Error())
		}
		list = append(list, s)
	}
	// Check for errors FROM iterating over rows.
	if err := rows.Err(); err != nil {
		db.Logger.Error("rows err on stage_encounter_activities", "error", err.Error())
	}
	return list, nil
}

// SELECT t.description, t.kind, t.ability, t.skill FROM tasks AS t JOIN stage_running_tasks AS s ON s.task_id = t.id WHERE s.task_id = 1
func (db *DBX) GetStageTaskFromRunningTaskID(ctx context.Context, taskID int) (types.Task, error) {
	t := types.Task{}
	query := "SELECT t.description, t.kind, t.ability, t.skill, t.target FROM tasks AS t JOIN stage_running_tasks AS s ON s.task_id = t.id WHERE s.task_id = $1" // dev:finder+query
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
	// Check for errors FROM iterating over rows.
	if err := rows.Err(); err != nil {
		db.Logger.Error("rows err on tasks by task_id", "error", err.Error())
	}
	return t, nil
}

func (db *DBX) GetCreatureFromParticipantsList(ctx context.Context, players []types.Options, npcs []types.Options, rpgSystem *rpg.RPGSystem) (map[int]*rules.Creature, map[int]*rules.Creature, error) {
	db.Logger.Debug("GetCreatureFromParticipantsList", "players", players, "npcs", npcs)
	// players
	playersMap := map[int]*rules.Creature{}
	creatureMap := map[int]*rules.Creature{}
	{
		query := "SELECT id, character_name, abilities, skills, extensions FROM players" // dev:finder+query
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
		query := "SELECT id, npc_name, abilities, skills, extensions FROM non_players" // dev:finder+query
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
func (db *DBX) GetNextEncounterByEncounterID(ctx context.Context, id int) (types.Next, error) {
	ne := types.Next{}
	query := "SELECT display_text, upstream_id, current_encounter_id, next_encounter_id FROM stage_next_encounter WHERE current_encounter_id = $1" // dev:finder+query
	rows, err := db.Conn.QueryContext(ctx, query, id)
	if err != nil {
		db.Logger.Error("query on stage_next_encounter by encounter_id failed", "error", err.Error())
		return ne, err
	}
	defer rows.Close()
	for rows.Next() {
		if err := rows.Scan(&ne.Text, &ne.UpstreamID, &ne.EncounterID, &ne.NextEncounterID); err != nil {
			db.Logger.Error("scan error on stage_next_encounter by encounter_id ", "error", err.Error())
		}
	}
	// Check for errors FROM iterating over rows.
	if err := rows.Err(); err != nil {
		db.Logger.Error("rows err on stage_next_encounter by encounter_id", "error", err.Error())
	}
	return ne, nil
}

// get next stage encounter by stage id
func (db *DBX) GetNextEncounterByStageID(ctx context.Context, id int) ([]types.Next, error) {
	next := []types.Next{}
	query := "SELECT s.id, s.upstream_id, s.current_encounter_id, s.next_encounter_id, s.display_text, o.kind, o.values FROM stage_next_encounter AS s JOIN stage_next_objectives AS o ON s.id = o.upstream_id WHERE s.upstream_id = $1" // dev:finder+query
	rows, err := db.Conn.QueryContext(ctx, query, id)
	if err != nil {
		db.Logger.Error("query on stage_next_encounter failed", "error", err.Error())
		return next, err
	}
	defer rows.Close()
	for rows.Next() {
		var n types.Next
		var o types.Objective
		var values []sql.NullInt64
		if err := rows.Scan(&n.ID, &n.UpstreamID, &n.EncounterID, &n.NextEncounterID, &n.Text, &o.Kind, pq.Array(&values)); err != nil {
			db.Logger.Error("scan error on stage_next_encounter ", "error", err.Error())
		}
		n.Objective = o
		if len(values) > 0 {
			for _, v := range values {
				if v.Valid {
					n.Objective.Values = append(n.Objective.Values, int(v.Int64))
				}
			}
		}

		next = append(next, n)
	}
	// Check for errors FROM iterating over rows.
	if err := rows.Err(); err != nil {
		db.Logger.Error("rows err on stage_next_encounter", "error", err.Error())
	}
	return next, nil

}

func (db *DBX) GetStageEncounterListByStoryID(ctx context.Context, storyID int) (types.EncounterList, error) {
	list := types.EncounterList{}
	query := "SELECT a.id, e.title AS encounter, e.id AS encounter_id, n.title AS next_encounter, n.id AS next_id FROM stage_next_encounter AS a JOIN encounters AS e ON e.id = a.current_encounter_id JOIN encounters AS n ON n.id = a.next_encounter_id WHERE e.story_id = $1" // dev:finder+query
	rows, err := db.Conn.QueryContext(ctx, query, storyID)
	if err != nil {
		db.Logger.Error("query on auto_play_next_encounter by story_id failed", "error", err.Error())
		return list, err
	}
	defer rows.Close()
	for rows.Next() {
		var next types.EncounterWithNext
		if err := rows.Scan(&next.ID, &next.Encounter, &next.EncounterID, &next.NextEncounter, &next.NextID); err != nil {
			db.Logger.Error("scan error on auto_play_next_encounter by story_id ", "error", err.Error())
		}
		list.Link = append(list.Link, next)
	}
	// Check for errors FROM iterating over rows.
	if err := rows.Err(); err != nil {
		db.Logger.Error("rows err on auto_play_next_encounter by story_id", "error", err.Error())
	}
	queryEncounter := "SELECT id, title AS name FROM encounters WHERE story_id = $1" // dev:finder+query
	rowsEncounter, err := db.Conn.QueryContext(ctx, queryEncounter, storyID)
	if err != nil {
		db.Logger.Error("query on encounters by story_id failed", "error", err.Error())
		return list, err
	}
	defer rowsEncounter.Close()
	for rowsEncounter.Next() {
		var generic types.Options
		if err := rowsEncounter.Scan(&generic.ID, &generic.Name); err != nil {
			db.Logger.Error("scan error on encounters by story_id ", "error", err.Error())
		}
		list.EncounterList = append(list.EncounterList, generic)
	}
	return list, nil
}

// stage_running_tasks
func (db *DBX) getRunningTaskByEncounterID(ctx context.Context, id int) ([]types.Options, error) {
	tasks := []types.Options{}
	query := "SELECT sa.display_text, sa.id FROM stage_running_tasks AS sa WHERE sa.stage_encounter_id = $1" // dev:finder+query
	rows, err := db.Conn.QueryContext(ctx, query, id)
	if err != nil {
		db.Logger.Error("query on stage_running_tasks by encounter_id failed", "error", err.Error())
		return tasks, err
	}
	defer rows.Close()
	for rows.Next() {
		var s types.Options
		if err := rows.Scan(&s.Name, &s.ID); err != nil {
			db.Logger.Error("scan error on stage_running_tasks by encounter_id ", "error", err.Error())
		}
		tasks = append(tasks, s)
	}
	// Check for errors FROM iterating over rows.
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
		query := "SELECT se.id, se.display_text, e.title, se.storyteller_id, se.phase, e.notes, e.announcement, e.writer_id, e.story_id, s.encoding_key, i.id FROM stage_encounters AS se JOIN encounters AS e ON se.encounter_id = e.id JOIN stage AS s ON se.stage_id = s.id LEFT JOIN initiative AS i ON i.stage_encounter_id = se.id WHERE se.ID = $1 AND s.finished = false" // dev:finder+query
		rows, err := db.Conn.QueryContext(ctx, query, id)
		if err != nil {
			db.Logger.Error("query on stage_encounters by stage_encounter.id failed", "error", err.Error())
			return enc, err
		}
		defer rows.Close()
		r = rows
	case phase > 0 && id == -1:
		// can return a initiative id
		query := "SELECT se.id, se.display_text, e.title, se.storyteller_id, se.phase, e.notes, e.announcement, e.writer_id, e.story_id, s.encoding_key, i.id FROM stage_encounters AS se JOIN encounters AS e ON se.encounter_id = e.id JOIN stage AS s ON se.stage_id = s.id LEFT JOIN initiative AS i ON i.stage_encounter_id = se.id WHERE se.phase = $1 AND s.finished = false" // dev:finder+query
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
			db.Logger.Error("scan error on stage_encounters by stage_encounter.ID ", "error", err.Error())
		}
		s.Notes, _ = utils.DecryptText(notes, encodingKey)
		s.Announcement, _ = utils.DecryptText(announce, encodingKey)
		if initiativeID.Valid {
			s.InitiativeID = int(initiativeID.Int64)
		}
		enc = s
	}
	// Check for errors FROM iterating over rows.
	if err := r.Err(); err != nil {
		db.Logger.Error("rows err on stage_encounters by stage_encounter.ID", "error", err.Error())
	}
	// participants
	p, n, err := db.getParticipantsByStageEncounterID(ctx, enc.ID)
	if err != nil {
		db.Logger.Error("error on db.getStageEncounterByEncounterID when calling getParticipantsByStageEncounterID", "error", err.Error())
	}
	db.Logger.Debug("players and npcs", "encounter_id", enc.ID, "player", p, "npc", n)
	enc.PC = p
	enc.NPC = n
	return enc, nil
}

func (db *DBX) getParticipantsByStageEncounterID(ctx context.Context, id int) ([]types.Options, []types.Options, error) {
	players := []types.Options{}
	npcs := []types.Options{}
	query := "SELECT sp.players_id, pl.character_name, pl.destroyed, snp.non_players_id, npc.npc_name, npc.destroyed FROM stage_encounters AS se LEFT JOIN stage_encounters_participants_players AS sp ON sp.stage_encounter_id = se.id LEFT JOIN players AS pl ON pl.id = sp.players_id LEFT JOIN stage_encounters_participants_non_players AS snp ON snp.stage_encounter_id = se.id LEFT JOIN non_players AS npc ON npc.id = snp.non_players_id WHERE se.ID = $1" // dev:finder+query
	rows, err := db.Conn.QueryContext(ctx, query, id)
	if err != nil {
		db.Logger.Error("query on stage_encounters_participants_players and stage_encounters_participants_non_players failed", "error", err.Error())
		return players, npcs, err
	}
	defer rows.Close()
	var p types.Options
	var n types.Options
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
	// Check for errors FROM iterating over rows.
	if err := rows.Err(); err != nil {
		db.Logger.Error("rows err on stage_encounters_participants_players and stage_encounters_participants_non_players", "error", err.Error())
	}
	return players, npcs, err
}

package db

// func (db *DBX) GetSceneByChannelID(ctx context.Context, id string) (types.Scene, error) {
// 	scene := types.Scene{}

// 	var story types.Story
// 	rows, err := db.Conn.QueryContext(ctx, "SELECT s.id, s.title, s.announcement, s.notes, s.master_id FROM story AS s JOIN story_channels AS c ON c.story_id = s.id WHERE c.channel = $1", id)
// 	if err != nil {
// 		db.logger.Error("query on story by id failed", "error", err.Error())
// 		return scene, err
// 	}
// 	defer rows.Close()
// 	for rows.Next() {
// 		if err := rows.Scan(&story.ID, &story.Title, &story.Announcement, &story.Notes, &story.MasterID); err != nil {
// 			db.logger.Error("scan error on story by id", "error", err.Error())
// 		}
// 	}
// 	// Check for errors from iterating over rows.
// 	if err := rows.Err(); err != nil {
// 		db.logger.Error("rows error on story by id", "error", err.Error())
// 	}
// 	scene.Story = story
// 	encounters, err := db.GetEncounterByStoryID(ctx, scene.Story.ID)
// 	if err != nil {
// 		db.logger.Error("cannot find any encounter for scene", "error", err.Error())
// 	}
// 	for _, enc := range encounters {
// 		if !enc.Finished && enc.Phase == types.Running {
// 			scene.Encounter = enc
// 		}
// 	}
// 	tasks, err := db.GetTasksByEncounterID(ctx, scene.Encounter.ID)
// 	if err != nil {
// 		db.logger.Error("cannot find any tasks for scene", "error", err.Error())
// 	}
// 	scene.Tasks = tasks
// 	combat := false
// 	for _, v := range scene.Tasks {
// 		if v.Kind == types.CombatTask {
// 			combat = true
// 			scene.Combat = true
// 		}
// 	}
// 	if combat {
// 		initiative, _, err := db.GetRunningInitiativeByEncounterID(ctx, scene.Encounter.ID)
// 		if err != nil {
// 			db.logger.Error("cannot find any initiative for scene", "error", err.Error())
// 		}
// 		scene.Initiative = &initiative
// 	}
// 	return scene, nil
// }

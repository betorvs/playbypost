package db

import (
	"context"
	"fmt"

	"github.com/betorvs/playbypost/core/rpg"
	"github.com/betorvs/playbypost/core/rpg/d10hm"
	"github.com/betorvs/playbypost/core/rpg/d10os"
	"github.com/betorvs/playbypost/core/rpg/d20e35"
	"github.com/betorvs/playbypost/core/rules"
	"github.com/betorvs/playbypost/core/sys/web/types"
)

func (db *DBX) GetPlayer(ctx context.Context, id int, npc bool, rpgSystem *rpg.RPGSystem) (*rules.Creature, error) {
	obj := rules.RestoreCreature()
	var query string
	switch rpgSystem.Name {
	case rpg.D2035:
		query = "SELECT p.character_name, p.destroyed, p.abilities, p.skills, p.rpg, x.level_total, x.hit_points, x.armor_class, x.class, x.multiclass, x.race, x.size, x.weapon FROM players as p JOIN extension_d20_e35_pc as x ON x.player_id = p.id WHERE p.id = $1"
		if npc {
			query = "SELECT p.npc_name, p.destroyed, p.abilities, p.skills, p.rpg, x.level_total, x.hit_points, x.armor_class, x.class, x.multiclass, x.race, x.size, x.weapon FROM non_players as p JOIN extension_d20_e35_npc as x ON x.player_id = p.id WHERE p.id = $1"
		}

	case rpg.D10HM:
		query = "SELECT p.character_name, p.destroyed, p.abilities, p.skills, p.rpg, x.health, x.defense, x.willpower, x.initiative, x.size, x.weapon FROM players as p JOIN extension_d10_homemade_pc as x ON x.player_id = p.id WHERE p.id = $1"
		if npc {
			query = "SELECT p.npc_name, p.destroyed, p.abilities, p.skills, p.rpg, x.health, x.defense, x.willpower, x.initiative, x.size, x.weapon FROM non_players as p JOIN extension_d10_homemade_npc as x ON x.player_id = p.id WHERE p.id = $1"
		}

	case rpg.D10OS:
		query = "SELECT p.character_name, p.destroyed, p.abilities, p.skills, p.rpg, x.health, x.willpower, x.initiative, x.size, x.armor, x.conscience_conviction, x.self_control_instinct, x.courage, x.weapon FROM players as p JOIN extension_d10_oldschool_pc as x ON x.player_id = p.id WHERE p.id = $1"
		if npc {
			query = "SELECT p.npc_name, p.destroyed, p.abilities, p.skills, p.rpg, x.health, x.willpower, x.initiative, x.size, x.armor, x.conscience_conviction, x.self_control_instinct, x.courage, x.weapon FROM non_players as p JOIN extension_d10_oldschool_npc as x ON x.player_id = p.id WHERE p.id = $1"
		}
	}

	rows, err := db.Conn.QueryContext(ctx, query, id)
	if err != nil {
		db.logger.Error("query on players/non_players by id failed", "non_player_query", npc, "error", err.Error())
		return obj, err
	}
	defer rows.Close()
	var destroyed bool
	var rpgName string
	for rows.Next() {
		switch rpgSystem.Name {
		case rpg.D2035:
			d2035Extension := d20e35.RestoreExtended()
			if err := rows.Scan(&obj.Name, &destroyed, &obj.Abilities, &obj.Skills, &rpgName, &d2035Extension.Level, &d2035Extension.HitPoints, &d2035Extension.ArmorClass, &d2035Extension.Class, &d2035Extension.Multiclass, &d2035Extension.Race, &d2035Extension.Size, &d2035Extension.Weapon); err != nil {
				db.logger.Error("scan on players/non_players by id  error", "error", err.Error())
			}
			obj.Extension = d2035Extension
			obj.RPG = rpgSystem

		case rpg.D10HM:
			d10HMExtension := d10hm.RestoreExtended()
			if err := rows.Scan(&obj.Name, &destroyed, &obj.Abilities, &obj.Skills, &rpgName, &d10HMExtension.Health, &d10HMExtension.Defense, &d10HMExtension.WillPower, &d10HMExtension.Initiative, &d10HMExtension.Size, &d10HMExtension.Weapon); err != nil {
				db.logger.Error("scan on players/non_players by id  error", "error", err.Error())
			}
			obj.Extension = d10HMExtension
			obj.RPG = rpgSystem

		case rpg.D10OS:
			d10OSExtension := d10os.RestoreExtended()
			if err := rows.Scan(&obj.Name, &destroyed, &obj.Abilities, &obj.Skills, &rpgName, &d10OSExtension.Health, &d10OSExtension.WillPower, &d10OSExtension.Initiative, &d10OSExtension.Size, &d10OSExtension.Armor, &d10OSExtension.Virtues.ConscienceConviction, &d10OSExtension.Virtues.SelfControlInstinct, &d10OSExtension.Virtues.Courage, &d10OSExtension.Weapon); err != nil {
				db.logger.Error("scan on players/non_players by id error", "error", err.Error())
			}
			obj.Extension = d10OSExtension
			obj.RPG = rpgSystem
		}

	}
	if destroyed {
		_ = obj.Destroy()
	}
	return obj, nil
}

func (db *DBX) GetPlayerByUserID(ctx context.Context, id int, npc bool, rpg *rpg.RPGSystem) (*rules.Creature, error) {
	res := rules.Creature{}

	var query string
	switch {
	case rpg.Name == "d20-3.5":
		query = "SELECT p.id, p.character_name, p.destroyed, p.abilities, p.skills, p.rpg, x.level_total, x.hit_points, x.armor_class, x.class, x.multiclass, x.race, x.size, x.weapon FROM players as p JOIN extension_d20_e35_pc as x ON x.player_id = p.id WHERE p.player_id = $1"
		if npc {
			query = "SELECT p.id, p.npc_name, p.destroyed, p.abilities, p.skills, p.rpg, x.level_total, x.hit_points, x.armor_class, x.class, x.multiclass, x.race, x.size, x.weapon FROM non_players as p JOIN extension_d20_e35_npc as x ON x.player_id = p.id WHERE p.player_id = $1"
		}

	case rpg.Name == "D10HomeMade":
		query = "SELECT p.id, p.character_name, p.destroyed, p.abilities, p.skills, p.rpg, x.health, x.defense, x.willpower, x.initiative, x.size, x.weapon FROM players as p JOIN extension_d10_homemade_pc as x ON x.player_id = p.id WHERE p.player_id = $1"
		if npc {
			query = "SELECT p.id, p.npc_name, p.destroyed, p.abilities, p.skills, p.rpg, x.health, x.defense, x.willpower, x.initiative, x.size, x.weapon FROM non_players as p JOIN extension_d10_homemade_npc as x ON x.player_id = p.id WHERE p.player_id = $1"
		}

	case rpg.Name == "D10OldSchool":
		query = "SELECT p.id, p.character_name, p.destroyed, p.abilities, p.skills, p.rpg, x.health, x.willpower, x.initiative, x.size, x.armor, x.conscience_conviction, x.self_control_instinct, x.courage, x.weapon FROM players as p JOIN extension_d10_oldschool_pc as x ON x.player_id = p.id WHERE p.player_id = $1"
		if npc {
			query = "SELECT p.id, p.npc_name, p.destroyed, p.abilities, p.skills, p.rpg, x.health, x.willpower, x.initiative, x.size, x.armor, x.conscience_conviction, x.self_control_instinct, x.courage, x.weapon FROM non_players as p JOIN extension_d10_oldschool_npc as x ON x.player_id = p.id WHERE p.player_id = $1"
		}
	}

	rows, err := db.Conn.QueryContext(ctx, query, id)
	if err != nil {
		db.logger.Error("query on players/non_players by player id failed", "non_player_query", npc, "error", err.Error())
		return &res, err
	}
	defer rows.Close()
	for rows.Next() {
		obj := rules.RestoreCreature()
		var id int
		var destroyed bool
		var rpgName string
		switch {
		case rpg.Name == "d20-3.5":
			d2035Extension := d20e35.RestoreExtended()
			if err := rows.Scan(&id, &obj.Name, &destroyed, &obj.Abilities, &obj.Skills, &rpgName, &d2035Extension.Level, &d2035Extension.HitPoints, &d2035Extension.ArmorClass, &d2035Extension.Class, &d2035Extension.Multiclass, &d2035Extension.Race, &d2035Extension.Size, &d2035Extension.Weapon); err != nil {
				db.logger.Error("scan on players/non_players by player id error", "error", err.Error())
			}
			obj.Extension = d2035Extension
			obj.RPG = rpg

		case rpg.Name == "D10HomeMade":
			d10HMExtension := d10hm.RestoreExtended()
			if err := rows.Scan(&id, &obj.Name, &destroyed, &obj.Abilities, &obj.Skills, &rpgName, &d10HMExtension.Health, &d10HMExtension.Defense, &d10HMExtension.WillPower, &d10HMExtension.Initiative, &d10HMExtension.Size, &d10HMExtension.Weapon); err != nil {
				db.logger.Error("scan on players/non_players by player id error", "error", err.Error())
			}
			obj.Extension = d10HMExtension
			obj.RPG = rpg

		case rpg.Name == "D10OldSchool":
			d10OSExtension := d10os.RestoreExtended()
			if err := rows.Scan(&id, &obj.Name, &destroyed, &obj.Abilities, &obj.Skills, &rpgName, &d10OSExtension.Health, &d10OSExtension.WillPower, &d10OSExtension.Initiative, &d10OSExtension.Size, &d10OSExtension.Armor, &d10OSExtension.Virtues.ConscienceConviction, &d10OSExtension.Virtues.SelfControlInstinct, &d10OSExtension.Virtues.Courage, &d10OSExtension.Weapon); err != nil {
				db.logger.Error("scan on players/non_players by player id error", "error", err.Error())
			}
			obj.Extension = d10OSExtension
			obj.RPG = rpg
		}
		if destroyed {
			_ = obj.Destroy()
		}
		res = *obj
	}
	return &res, nil
}

func (db *DBX) GetPlayersByStoryID(ctx context.Context, storyID int, npc bool, rpg *rpg.RPGSystem) (map[int]*rules.Creature, error) {
	res := map[int]*rules.Creature{}

	var query string
	switch {
	case rpg.Name == "d20-3.5":
		query = "SELECT p.id, p.character_name, p.destroyed, p.abilities, p.skills, p.rpg, x.level_total, x.hit_points, x.armor_class, x.class, x.multiclass, x.race, x.size, x.weapon FROM players as p JOIN extension_d20_e35_pc as x ON x.player_id = p.id WHERE p.story_id = $1"
		if npc {
			query = "SELECT p.id, p.npc_name, p.destroyed, p.abilities, p.skills, p.rpg, x.level_total, x.hit_points, x.armor_class, x.class, x.multiclass, x.race, x.size, x.weapon FROM non_players as p JOIN extension_d20_e35_npc as x ON x.player_id = p.id WHERE p.story_id = $1"
		}

	case rpg.Name == "D10HomeMade":
		query = "SELECT p.id, p.character_name, p.destroyed, p.abilities, p.skills, p.rpg, x.health, x.defense, x.willpower, x.initiative, x.size, x.weapon FROM players as p JOIN extension_d10_homemade_pc as x ON x.player_id = p.id WHERE p.story_id = $1"
		if npc {
			query = "SELECT p.id, p.npc_name, p.destroyed, p.abilities, p.skills, p.rpg, x.health, x.defense, x.willpower, x.initiative, x.size, x.weapon FROM non_players as p JOIN extension_d10_homemade_npc as x ON x.player_id = p.id WHERE p.story_id = $1"
		}

	case rpg.Name == "D10OldSchool":
		query = "SELECT p.id, p.character_name, p.destroyed, p.abilities, p.skills, p.rpg, x.health, x.willpower, x.initiative, x.size, x.armor, x.conscience_conviction, x.self_control_instinct, x.courage, x.weapon FROM players as p JOIN extension_d10_oldschool_pc as x ON x.player_id = p.id WHERE p.story_id = $1"
		if npc {
			query = "SELECT p.id, p.npc_name, p.destroyed, p.abilities, p.skills, p.rpg, x.health, x.willpower, x.initiative, x.size, x.armor, x.conscience_conviction, x.self_control_instinct, x.courage, x.weapon FROM non_players as p JOIN extension_d10_oldschool_npc as x ON x.player_id = p.id WHERE p.story_id = $1"
		}
	}

	rows, err := db.Conn.QueryContext(ctx, query, storyID)
	if err != nil {
		db.logger.Error("query on players/non_players by story id failed", "non_player_query", npc, "error", err.Error())
		return res, err
	}
	defer rows.Close()
	for rows.Next() {
		obj := rules.RestoreCreature()
		var id int
		var destroyed bool
		var rpgName string
		switch {
		case rpg.Name == "d20-3.5":
			d2035Extension := d20e35.RestoreExtended()
			if err := rows.Scan(&id, &obj.Name, &destroyed, &obj.Abilities, &obj.Skills, &rpgName, &d2035Extension.Level, &d2035Extension.HitPoints, &d2035Extension.ArmorClass, &d2035Extension.Class, &d2035Extension.Multiclass, &d2035Extension.Race, &d2035Extension.Size, &d2035Extension.Weapon); err != nil {
				db.logger.Error("scan on players/non_players by story id error", "error", err.Error())
			}
			obj.Extension = d2035Extension
			obj.RPG = rpg

		case rpg.Name == "D10HomeMade":
			d10HMExtension := d10hm.RestoreExtended()
			if err := rows.Scan(&id, &obj.Name, &destroyed, &obj.Abilities, &obj.Skills, &rpgName, &d10HMExtension.Health, &d10HMExtension.Defense, &d10HMExtension.WillPower, &d10HMExtension.Initiative, &d10HMExtension.Size, &d10HMExtension.Weapon); err != nil {
				db.logger.Error("scan on players/non_players by story id error", "error", err.Error())
			}
			obj.Extension = d10HMExtension
			obj.RPG = rpg

		case rpg.Name == "D10OldSchool":
			d10OSExtension := d10os.RestoreExtended()
			if err := rows.Scan(&id, &obj.Name, &destroyed, &obj.Abilities, &obj.Skills, &rpgName, &d10OSExtension.Health, &d10OSExtension.WillPower, &d10OSExtension.Initiative, &d10OSExtension.Size, &d10OSExtension.Armor, &d10OSExtension.Virtues.ConscienceConviction, &d10OSExtension.Virtues.SelfControlInstinct, &d10OSExtension.Virtues.Courage, &d10OSExtension.Weapon); err != nil {
				db.logger.Error("scan on players/non_players by story id error", "error", err.Error())
			}
			obj.Extension = d10OSExtension
			obj.RPG = rpg
		}
		if destroyed {
			_ = obj.Destroy()
		}
		res[id] = obj
	}
	return res, nil
}

func (db *DBX) GetPlayersByEncounterID(ctx context.Context, encounterID int, npc bool, rpg *rpg.RPGSystem) (map[int]*rules.Creature, error) {
	res := map[int]*rules.Creature{}

	var query string
	switch {
	case rpg.Name == "d20-3.5":
		query = "SELECT p.id, p.character_name, p.destroyed, p.abilities, p.skills, p.rpg, x.level_total, x.hit_points, x.armor_class, x.class, x.multiclass, x.race, x.size, x.weapon FROM players as p JOIN extension_d20_e35_pc as x ON x.player_id = p.id JOIN encounters_participants_players as e ON e.players_id = p.id WHERE e.encounters_id = $1"
		if npc {
			query = "SELECT p.id, p.npc_name, p.destroyed, p.abilities, p.skills, p.rpg, x.level_total, x.hit_points, x.armor_class, x.class, x.multiclass, x.race, x.size, x.weapon FROM non_players as p JOIN extension_d20_e35_npc as x ON x.player_id = p.id JOIN encounters_participants_non_players as e ON e.players_id = p.id WHERE e.encounters_id = $1"
		}

	case rpg.Name == "D10HomeMade":
		query = "SELECT p.id, p.character_name, p.destroyed, p.abilities, p.skills, p.rpg, x.health, x.defense, x.willpower, x.initiative, x.size, x.weapon FROM players as p JOIN extension_d10_homemade_pc as x ON x.player_id = p.id JOIN encounters_participants_players as e ON e.players_id = p.id WHERE e.encounters_id = $1"
		if npc {
			query = "SELECT p.id, p.npc_name, p.destroyed, p.abilities, p.skills, p.rpg, x.health, x.defense, x.willpower, x.initiative, x.size, x.weapon FROM non_players as p JOIN extension_d10_homemade_npc as x ON x.player_id = p.id JOIN encounters_participants_non_players as e ON e.players_id = p.id WHERE e.encounters_id = $1"
		}

	case rpg.Name == "D10OldSchool":
		query = "SELECT p.id, p.character_name, p.destroyed, p.abilities, p.skills, p.rpg, x.health, x.willpower, x.initiative, x.size, x.armor, x.conscience_conviction, x.self_control_instinct, x.courage, x.weapon FROM players as p JOIN extension_d10_oldschool_pc as x ON x.player_id = p.id JOIN encounters_participants_players as e ON e.players_id = p.id WHERE e.encounters_id = $1"
		if npc {
			query = "SELECT p.id, p.npc_name, p.destroyed, p.abilities, p.skills, p.rpg, x.health, x.willpower, x.initiative, x.size, x.armor, x.conscience_conviction, x.self_control_instinct, x.courage, x.weapon FROM non_players as p JOIN extension_d10_oldschool_npc as x ON x.player_id = p.id JOIN encounters_participants_non_players as e ON e.players_id = p.id WHERE e.encounters_id = $1"
		}
	}

	rows, err := db.Conn.QueryContext(ctx, query, encounterID)
	if err != nil {
		db.logger.Error("query on players/non_players by encounter id failed", "non_player_query", npc, "error", err.Error())
		return res, err
	}
	defer rows.Close()
	for rows.Next() {
		obj := rules.RestoreCreature()
		var id int
		var destroyed bool
		var rpgName string
		switch {
		case rpg.Name == "d20-3.5":
			d2035Extension := d20e35.RestoreExtended()
			if err := rows.Scan(&id, &obj.Name, &destroyed, &obj.Abilities, &obj.Skills, &rpgName, &d2035Extension.Level, &d2035Extension.HitPoints, &d2035Extension.ArmorClass, &d2035Extension.Class, &d2035Extension.Multiclass, &d2035Extension.Race, &d2035Extension.Size, &d2035Extension.Weapon); err != nil {
				db.logger.Error("scan on players/non_players by encounter id error", "error", err.Error())
			}
			obj.Extension = d2035Extension
			obj.RPG = rpg

		case rpg.Name == "D10HomeMade":
			d10HMExtension := d10hm.RestoreExtended()
			if err := rows.Scan(&id, &obj.Name, &destroyed, &obj.Abilities, &obj.Skills, &rpgName, &d10HMExtension.Health, &d10HMExtension.Defense, &d10HMExtension.WillPower, &d10HMExtension.Initiative, &d10HMExtension.Size, &d10HMExtension.Weapon); err != nil {
				db.logger.Error("scan on players/non_players by encounter id error", "error", err.Error())
			}
			obj.Extension = d10HMExtension
			obj.RPG = rpg

		case rpg.Name == "D10OldSchool":
			d10OSExtension := d10os.RestoreExtended()
			if err := rows.Scan(&id, &obj.Name, &destroyed, &obj.Abilities, &obj.Skills, &rpgName, &d10OSExtension.Health, &d10OSExtension.WillPower, &d10OSExtension.Initiative, &d10OSExtension.Size, &d10OSExtension.Armor, &d10OSExtension.Virtues.ConscienceConviction, &d10OSExtension.Virtues.SelfControlInstinct, &d10OSExtension.Virtues.Courage, &d10OSExtension.Weapon); err != nil {
				db.logger.Error("scan on players/non_players by encounter id error", "error", err.Error())
			}
			obj.Extension = d10OSExtension
			obj.RPG = rpg
		}
		if destroyed {
			_ = obj.Destroy()
		}
		res[id] = obj
	}
	return res, nil
}

func (db *DBX) GetSliceOfPlayersByStoryID(ctx context.Context, storyID int, npc bool, rpgSystem *rpg.RPGSystem) ([]types.Players, error) {
	res := []types.Players{}

	var query string
	switch rpgSystem.Name {
	case rpg.D2035:
		query = "SELECT p.id, p.character_name, p.destroyed, p.abilities, p.skills, p.rpg, x.level_total, x.hit_points, x.armor_class, x.class, x.multiclass, x.race, x.size, x.weapon FROM players as p JOIN extension_d20_e35_pc as x ON x.player_id = p.id WHERE p.story_id = $1"
		if npc {
			query = "SELECT p.id, p.npc_name, p.destroyed, p.abilities, p.skills, p.rpg, x.level_total, x.hit_points, x.armor_class, x.class, x.multiclass, x.race, x.size, x.weapon FROM non_players as p JOIN extension_d20_e35_npc as x ON x.player_id = p.id WHERE p.story_id = $1"
		}

	case rpg.D10HM:
		query = "SELECT p.id, p.character_name, p.destroyed, p.abilities, p.skills, p.rpg, x.health, x.defense, x.willpower, x.initiative, x.size, x.weapon FROM players as p JOIN extension_d10_homemade_pc as x ON x.player_id = p.id WHERE p.story_id = $1"
		if npc {
			query = "SELECT p.id, p.npc_name, p.destroyed, p.abilities, p.skills, p.rpg, x.health, x.defense, x.willpower, x.initiative, x.size, x.weapon FROM non_players as p JOIN extension_d10_homemade_npc as x ON x.player_id = p.id WHERE p.story_id = $1"
		}

	case rpg.D10OS:
		query = "SELECT p.id, p.character_name, p.destroyed, p.abilities, p.skills, p.rpg, x.health, x.willpower, x.initiative, x.size, x.armor, x.conscience_conviction, x.self_control_instinct, x.courage, x.weapon FROM players as p JOIN extension_d10_oldschool_pc as x ON x.player_id = p.id WHERE p.story_id = $1"
		if npc {
			query = "SELECT p.id, p.npc_name, p.destroyed, p.abilities, p.skills, p.rpg, x.health, x.willpower, x.initiative, x.size, x.armor, x.conscience_conviction, x.self_control_instinct, x.courage, x.weapon FROM non_players as p JOIN extension_d10_oldschool_npc as x ON x.player_id = p.id WHERE p.story_id = $1"
		}
	}

	rows, err := db.Conn.QueryContext(ctx, query, storyID)
	if err != nil {
		db.logger.Error("query on players/non_players by story id failed", "non_player_query", npc, "error", err.Error())
		return res, err
	}
	defer rows.Close()
	for rows.Next() {
		obj := types.Players{
			Abilities: make(map[string]int),
			Skills:    make(map[string]int),
			Extension: make(map[string]int),
			Details:   make(map[string]string),
		}
		abilities := rules.Abilities{}
		skills := rules.Skills{}

		switch rpgSystem.Name {
		case rpg.D2035:
			d2035Extension := d20e35.RestoreExtended()

			if err := rows.Scan(&obj.ID, &obj.Name, &obj.Destroyed, &abilities, &skills, &obj.RPG, &d2035Extension.Level, &d2035Extension.HitPoints, &d2035Extension.ArmorClass, &d2035Extension.Class, &d2035Extension.Multiclass, &d2035Extension.Race, &d2035Extension.Size, &d2035Extension.Weapon); err != nil {
				db.logger.Error("scan on players/non_players by story id error", "error", err.Error())
			}

			obj.Extension["hit_points"] = d2035Extension.HitPoints
			obj.Extension["level"] = d2035Extension.Level
			obj.Extension["armor_class"] = d2035Extension.ArmorClass
			obj.Details["class"] = fmt.Sprintf("%v", d2035Extension.Class)
			obj.Details["race"] = d2035Extension.Race
			obj.Details["size"] = d2035Extension.Size
			// obj.Extension = d2035Extension

		case rpg.D10HM:
			d10HMExtension := d10hm.RestoreExtended()
			if err := rows.Scan(&obj.ID, &obj.Name, &obj.Destroyed, &abilities, &skills, &obj.RPG, &d10HMExtension.Health, &d10HMExtension.Defense, &d10HMExtension.WillPower, &d10HMExtension.Initiative, &d10HMExtension.Size, &d10HMExtension.Weapon); err != nil {
				db.logger.Error("scan on players/non_players by story id error", "error", err.Error())
			}
			for k, v := range abilities {
				if v.DisplayName == "" {
					obj.Abilities[k] = v.Value
				} else {
					obj.Abilities[v.DisplayName] = v.Value
				}

			}
			for k, v := range skills {
				obj.Skills[k] = v.Value
			}
			obj.Extension["health"] = d10HMExtension.Health
			obj.Extension["will_power"] = d10HMExtension.WillPower
			obj.Extension["defense"] = d10HMExtension.Defense
			obj.Extension["initiative"] = d10HMExtension.Initiative

		case rpg.D10OS:
			d10OSExtension := d10os.RestoreExtended()
			if err := rows.Scan(&obj.ID, &obj.Name, &obj.Destroyed, &abilities, &skills, &obj.RPG, &d10OSExtension.Health, &d10OSExtension.WillPower, &d10OSExtension.Initiative, &d10OSExtension.Size, &d10OSExtension.Armor, &d10OSExtension.Virtues.ConscienceConviction, &d10OSExtension.Virtues.SelfControlInstinct, &d10OSExtension.Virtues.Courage, &d10OSExtension.Weapon); err != nil {
				db.logger.Error("scan on players/non_players by story id error", "error", err.Error())
			}
			for k, v := range abilities {
				if v.DisplayName == "" {
					obj.Abilities[k] = v.Value
				} else {
					obj.Abilities[v.DisplayName] = v.Value
				}

			}
			for k, v := range skills {
				obj.Skills[k] = v.Value
			}
			obj.Extension["health"] = d10OSExtension.Health
			obj.Extension["will_power"] = d10OSExtension.WillPower
			obj.Extension["initiative"] = d10OSExtension.Initiative
			obj.Extension["armor"] = d10OSExtension.Armor
		}
		res = append(res, obj)
	}
	return res, nil
}

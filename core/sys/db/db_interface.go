package db

import (
	"context"

	"github.com/betorvs/playbypost/core/initiative"
	"github.com/betorvs/playbypost/core/rpg"
	"github.com/betorvs/playbypost/core/rules"
	"github.com/betorvs/playbypost/core/sys/web/types"
)

type DBClient interface {
	Close() error
	// Storytellers
	// GetStorytellers active true to get only active users
	GetStorytellers(ctx context.Context, active bool) ([]types.Storyteller, error)
	CreateStorytellers(ctx context.Context, username, password string) (int, error)
	GetStorytellerByID(ctx context.Context, id int) (types.Storyteller, error)
	// GetUserByUserID(ctx context.Context, userid string) (types.User, error)
	GetStorytellerByUsername(ctx context.Context, username string) (types.Storyteller, error)
	GetUserCard(ctx context.Context) ([]types.Card, error)
	// Story
	// GetStory notfFinished true to get only not finished stories
	GetStory(ctx context.Context) ([]types.Story, error)
	GetStoriesByStorytellerID(ctx context.Context, id int) ([]types.Story, error)
	// CreateStory(ctx context.Context, title, announcement, notes string, masterID int) (int, error)
	CreateStoryTx(ctx context.Context, title, announcement, notes, encodingKey string, masterID int) (int, error)
	GetStoryIDByTitle(ctx context.Context, title string) (int, error)
	GetStoryByID(ctx context.Context, id int) (types.Story, error)
	GetStoryChannels(ctx context.Context) (map[string]int, error)
	// Encounters
	GetEncounters(ctx context.Context) ([]types.Encounters, error)
	GetEncounterByStoryID(ctx context.Context, storyID int) ([]types.Encounters, error)
	GetEncounterByID(ctx context.Context, id int) (types.Encounters, error)
	CreateEncounter(ctx context.Context, title, announcement, notes string, storyID, storytellerID int) (int, error)
	// UpdatePhase(ctx context.Context, id, phase int) error
	AddParticipants(ctx context.Context, encounterID int, npc bool, players []int) error
	CreateTask(ctx context.Context, title, subject, checks string, kind, target, encounterID int, options map[string]int) (int, error)
	GetTasksByEncounterID(ctx context.Context, id int) (map[string]types.Task, error)
	// scene
	// GetSceneByChannelID(ctx context.Context, id string) (types.Scene, error)
	// Initiative
	UpdateNextPlayer(ctx context.Context, id, nextPlayer int) error
	SaveInitiativeTx(ctx context.Context, i initiative.Initiative, encounterID int) (int, error)
	SaveInitiative(ctx context.Context, i initiative.Initiative, encounterID int) (int, error)
	GetRunningInitiativeByStoryID(ctx context.Context, storyID int) (initiative.Initiative, int, error)
	GetRunningInitiativeByEncounterID(ctx context.Context, encounterID int) (initiative.Initiative, int, error)
	DeactivateParticipant(ctx context.Context, id int, name string) (int, error)
	// Players
	SavePlayer(ctx context.Context, id, storyID int, npc bool, creature *rules.Creature, rpg *rpg.RPGSystem) (int, error)
	SavePlayerTx(ctx context.Context, id, storyID int, npc bool, creature *rules.Creature, rpgSystem *rpg.RPGSystem) (int, error)
	GetPlayer(ctx context.Context, id int, npc bool, rpg *rpg.RPGSystem) (*rules.Creature, error)
	GetPlayerByUserID(ctx context.Context, id int, npc bool, rpg *rpg.RPGSystem) (*rules.Creature, error)
	// GetPlayersByStoryID(ctx context.Context, storyID int, npc bool, rpg *rpg.RPGSystem) (map[int]*rules.Creature, error)
	GetSliceOfPlayersByStoryID(ctx context.Context, storyID int, npc bool, rpgSystem *rpg.RPGSystem) ([]types.Players, error)
	GetPlayersByEncounterID(ctx context.Context, encounterID int, npc bool, rpg *rpg.RPGSystem) (map[int]*rules.Creature, error)
	// Extension
	SaveExtension(ctx context.Context, playerId int, npc bool, rpg *rpg.RPGSystem, extension interface{}) (int, error)
	// Slack
	AddSlackInformation(ctx context.Context, username, userid, channel string) (int, error)
}

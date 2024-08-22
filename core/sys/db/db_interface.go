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
	// Writers
	// GetWriters active true to get only active users
	GetWriters(ctx context.Context, active bool) ([]types.Writer, error)
	CreateWriters(ctx context.Context, username, password string) (int, error)
	GetWriterByID(ctx context.Context, id int) (types.Writer, error)
	// GetUserByUserID(ctx context.Context, userid string) (types.User, error)
	GetWriterByUsername(ctx context.Context, username string) (types.Writer, error)
	// Users - Slack
	// GetUserCard(ctx context.Context) ([]types.Card, error)
	// Story
	// GetStory notfFinished true to get only not finished stories
	GetStory(ctx context.Context) ([]types.Story, error)
	GetStoriesByWriterID(ctx context.Context, id int) ([]types.Story, error)
	// CreateStory(ctx context.Context, title, announcement, notes string, masterID int) (int, error)
	CreateStoryTx(ctx context.Context, title, announcement, notes, encodingKey string, masterID int) (int, error)
	GetStoryIDByTitle(ctx context.Context, title string) (int, error)
	GetStoryByID(ctx context.Context, id int) (types.Story, error)
	GetStoryChannels(ctx context.Context) (map[string]int, error)
	// Encounters
	GetEncounters(ctx context.Context) ([]types.Encounter, error)
	GetEncounterByStoryID(ctx context.Context, storyID int) ([]types.Encounter, error)
	GetEncounterByID(ctx context.Context, id int) (types.Encounter, error)
	CreateEncounter(ctx context.Context, title, announcement, notes string, storyID, storytellerID int, first, last bool) (int, error)
	// UpdatePhase(ctx context.Context, id, phase int) error
	AddParticipants(ctx context.Context, encounterID int, npc bool, players []int) error
	// Tasks
	GetTask(ctx context.Context) ([]types.Task, error)
	CreateTask(ctx context.Context, description, ability, skill string, kind types.TaskKind, target int) (int, error)
	// GetTasksByEncounterID(ctx context.Context, id int) (map[string]types.Task, error)
	// stage
	GetStage(ctx context.Context) ([]types.Stage, error)
	GetStageByStageID(ctx context.Context, id int) (types.StageAggregated, error)
	GetStageByStoryID(ctx context.Context, id int) ([]types.Stage, error)
	GetStageEncounterByEncounterID(ctx context.Context, id int) (types.StageEncounter, error)
	GetStageEncountersByStageID(ctx context.Context, id int) ([]types.StageEncounter, error)
	GetRunningStageByChannelID(ctx context.Context, channelID, userID string) (types.RunningStage, error)
	CreateStageTx(ctx context.Context, text, userid string, storyID int) (int, error)
	AddChannelToStage(ctx context.Context, channel string, id int) (int, error)
	AddEncounterToStage(ctx context.Context, text string, stage_id, storyteller_id, encounter_id int) (int, error)
	GetNextEncounterByEncounterID(ctx context.Context, id int) (types.NextEncounter, error)
	AddNextEncounter(ctx context.Context, text string, stageID, encounterID, nextEncounterID int) error
	AddRunningTask(ctx context.Context, text string, stageID, taskID, StorytellerID, encounterID int) error
	GetStageEncounterActivities(ctx context.Context) ([]types.StageEncounterActivities, error)
	GetStageEncounterActivitiesByEncounterID(ctx context.Context, id int) ([]types.StageEncounterActivities, error)
	AddEncounterActivities(ctx context.Context, text string, stageID, encounterID int) error
	UpdatePhase(ctx context.Context, id, phase int) error
	RegisterActivities(ctx context.Context, stageID, encounterID int, actions types.Actions) error
	UpdateProcessedActivities(ctx context.Context, id int, processed bool, actions types.Actions) error
	GetStageTaskFromRunningTaskID(ctx context.Context, taskID int) (types.Task, error)
	GetCreatureFromParticipantsList(ctx context.Context, players []types.GenericIDName, npcs []types.GenericIDName, rpgSystem *rpg.RPGSystem) (map[int]*rules.Creature, map[int]*rules.Creature, error)
	// NPC
	GetNPCByStageID(ctx context.Context, id int) ([]types.Players, error)
	GenerateNPC(ctx context.Context, name string, stageID, storytellerID int, creature *rules.Creature) (int, error)
	UpdateNPC(ctx context.Context, id int, creature *rules.Creature, destroyed bool) error
	// users
	GetUser(ctx context.Context) ([]types.User, error)
	GetUserByUserID(ctx context.Context, id string) (types.User, error)
	CreateUserTx(ctx context.Context, userid string) (int, error)
	// Initiative
	UpdateNextPlayer(ctx context.Context, id, nextPlayer int) error
	SaveInitiativeTx(ctx context.Context, i initiative.Initiative, encounterID int) (int, error)
	SaveInitiative(ctx context.Context, i initiative.Initiative, encounterID int) (int, error)
	GetInitiativeByID(ctx context.Context, id int) (initiative.Initiative, error)
	// GetRunningInitiativeByStageID(ctx context.Context, stageID int) (initiative.Initiative, int, error)
	GetRunningInitiativeByEncounterID(ctx context.Context, encounterID int) (initiative.Initiative, int, error)
	DeactivateParticipant(ctx context.Context, id int, name string) (int, error)
	// Players
	// SavePlayer(ctx context.Context, id, storyID int, npc bool, creature *rules.Creature, rpg *rpg.RPGSystem) (int, error)
	SavePlayerTx(ctx context.Context, id, storyID int, creature *rules.Creature, rpg *rpg.RPGSystem) (int, error)
	UpdatePlayer(ctx context.Context, id int, creature *rules.Creature, destroyed bool) error
	GetPlayers(ctx context.Context) ([]types.Players, error)
	GetPlayerByID(ctx context.Context, id int) (types.Players, error)
	GetPlayerByPlayerID(ctx context.Context, id int) (types.Players, error)
	GetPlayerByStageID(ctx context.Context, id int) ([]types.Players, error)
	GetPlayerByUserID(ctx context.Context, id, channel string) (types.Players, error)
	// GetPlayersByEncounterID(ctx context.Context, encounterID int, npc bool, rpg *rpg.RPGSystem) (map[int]*rules.Creature, error)
	// GetPlayersByStoryID(ctx context.Context, storyID int, npc bool, rpg *rpg.RPGSystem) (map[int]*rules.Creature, error)
	// GetSliceOfPlayersByStageID(ctx context.Context, stageID int, npc bool, rpgSystem *rpg.RPGSystem) ([]types.Players, error)
	// GetPlayersByEncounterID(ctx context.Context, encounterID int, npc bool, rpg *rpg.RPGSystem) (map[int]*rules.Creature, error)
	// Extension
	SaveExtension(ctx context.Context, playerId int, npc bool, rpg *rpg.RPGSystem, extension interface{}) (int, error)
	// Chat
	AddChatInformation(ctx context.Context, username, userid, channel, chat string) (int, error)
	GetChatInformation(ctx context.Context) ([]types.ChatInfo, error)
	GetChatChannelInformation(ctx context.Context) ([]string, error)
	//  Auto Play
	GetAutoPlay(ctx context.Context) ([]types.AutoPlay, error)
	GetAutoPlayByID(ctx context.Context, autoPlayID int) (types.AutoPlay, error)
	GetNextEncounterByStoryID(ctx context.Context, storyID int) (types.AutoPlayEncounterList, error)
	GetAutoPlayOptByChannelID(ctx context.Context, channelID, userID string) (types.AutoPlayOptions, error)
	CreateAutoPlayTx(ctx context.Context, text string, storyID int, solo bool) (int, error)
	AddAutoPlayNext(ctx context.Context, text string, autoPlayID, encounterID, nextEncounterID int) error
	CreateAutoPlayChannelTx(ctx context.Context, channelID, userID string, autoPlayID int) (int, error)
	RegisterActivitiesAutoPlay(ctx context.Context, autoPlayID, encounterID int, actions types.Actions) error
	GetAutoPlayActivities(ctx context.Context) ([]types.AutoPlayActivities, error)
	GetAnnounceByEncounterID(ctx context.Context, encounterID, autoPlayID int) (string, bool, error)
	UpdateProcessedAutoPlay(ctx context.Context, id int, processed bool, actions types.Actions) error
	UpdateAutoPlayState(ctx context.Context, autoPlayID int, encounterID int) error
	CloseAutoPlayChannel(ctx context.Context, channelID string, autoPlayID int) error
}

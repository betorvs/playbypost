package db

import (
	"context"
	"time"

	"github.com/betorvs/playbypost/core/initiative"
	"github.com/betorvs/playbypost/core/rpg"
	"github.com/betorvs/playbypost/core/rules"
	"github.com/betorvs/playbypost/core/sys/web/types"
)

type DBClient interface {
	Close() error
	// Writers
	GetWriters(ctx context.Context, active bool) ([]types.Writer, error)
	CreateWriters(ctx context.Context, username, password string) (int, error)
	GetWriterByID(ctx context.Context, id int) (types.Writer, error)
	GetWriterByUsername(ctx context.Context, username string) (types.Writer, error)
	// story
	GetStory(ctx context.Context) ([]types.Story, error)
	CreateStoryTx(ctx context.Context, title, announcement, notes, encodingKey string, writerID int) (int, error)
	UpdateStoryTx(ctx context.Context, title, announcement, notes string, storyID int) (int, error)
	GetStoryIDByTitle(ctx context.Context, title string) (int, error)
	GetStoryByID(ctx context.Context, id int) (types.Story, error)
	GetStoriesByWriterID(ctx context.Context, id int) ([]types.Story, error)
	// Encounters
	GetEncounters(ctx context.Context) ([]types.Encounter, error)
	GetEncounterByStoryID(ctx context.Context, storyID int) ([]types.Encounter, error)
	GetEncounterByID(ctx context.Context, id int) (types.Encounter, error)
	CreateEncounterTx(ctx context.Context, title, announcement, notes string, storyID, storytellerID int, first, last bool) (int, error)
	UpdateEncounterTx(ctx context.Context, title, announcement, notes string, id, storyID int, first, last bool) (int, error)
	// Tasks
	GetTask(ctx context.Context) ([]types.Task, error)
	CreateTask(ctx context.Context, description, ability, skill string, kind types.TaskKind, target int) (int, error)
	// GetTasksByEncounterID(ctx context.Context, id int) (map[string]types.Task, error)
	// stage
	GetStage(ctx context.Context) ([]types.Stage, error)
	GetStageByStoryID(ctx context.Context, id int) ([]types.Stage, error)
	GetStageByStageID(ctx context.Context, id int) (types.StageAggregated, error)
	GetStageEncounterByEncounterID(ctx context.Context, id int) (types.StageEncounter, error)
	GetStageEncountersByStageID(ctx context.Context, id int) ([]types.StageEncounter, error)
	GetRunningStageByChannelID(ctx context.Context, channelID, userID string, rpgSystem *rpg.RPGSystem) (types.RunningStage, error)
	GetStageEncounterActivitiesByEncounterID(ctx context.Context, id int) ([]types.Activity, error)
	GetStageEncounterActivities(ctx context.Context) ([]types.Activity, error)
	GetStageTaskFromRunningTaskID(ctx context.Context, taskID int) (types.Task, error)
	GetCreatureFromParticipantsList(ctx context.Context, players []types.Options, npcs []types.Options, rpgSystem *rpg.RPGSystem) (map[int]*rules.Creature, map[int]*rules.Creature, error)
	GetNextEncounterByEncounterID(ctx context.Context, id int) (types.Next, error)
	CreateStageTx(ctx context.Context, text, userid string, storyID int) (int, error)
	AddChannelToStage(ctx context.Context, channel string, id int) (int, error)
	AddEncounterToStage(ctx context.Context, text string, stage_id, storyteller_id, encounter_id int) (int, error)
	UpdatePhase(ctx context.Context, id, phase int) error
	AddParticipants(ctx context.Context, encounterID int, npc bool, players []int) error
	AddNextEncounter(ctx context.Context, next []types.Next) error
	AddRunningTask(ctx context.Context, text string, stageID, taskID, StorytellerID, encounterID int) error
	AddEncounterActivities(ctx context.Context, text string, stageID, encounterID int) error
	RegisterActivities(ctx context.Context, stageID, encounterID int, actions types.Actions) error
	UpdateProcessedActivities(ctx context.Context, id int, processed bool, actions types.Actions) error
	CloseStage(ctx context.Context, id int) error
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
	GetInitiativeByID(ctx context.Context, id int) (initiative.Initiative, error)
	GetRunningInitiativeByEncounterID(ctx context.Context, encounterID int) (initiative.Initiative, int, error)
	DeactivateParticipant(ctx context.Context, id int, name string) (int, error)
	// Players
	GetPlayers(ctx context.Context, rpgSystem *rpg.RPGSystem) ([]types.Players, error)
	GetPlayerByID(ctx context.Context, id int, rpgSystem *rpg.RPGSystem) (types.Players, error)
	GetPlayerByPlayerID(ctx context.Context, id int, rpgSystem *rpg.RPGSystem) (types.Players, error)
	GetPlayerByStageID(ctx context.Context, id int, rpgSystem *rpg.RPGSystem) ([]types.Players, error)
	GetPlayerByUserID(ctx context.Context, id, channel string, rpgSystem *rpg.RPGSystem) (types.Players, error)
	SavePlayerTx(ctx context.Context, id, storyID int, creature *rules.Creature, rpgSystem *rpg.RPGSystem) (int, error)
	UpdatePlayer(ctx context.Context, id int, creature *rules.Creature, destroyed bool) error
	// Chat
	AddChatInformation(ctx context.Context, username, userid, channel, chat string) (int, error)
	GetChatInformation(ctx context.Context) ([]types.ChatInfo, error)
	GetChatChannelInformation(ctx context.Context) ([]string, error)
	//  Auto Play
	GetAutoPlay(ctx context.Context) ([]types.AutoPlay, error)
	GetAutoPlayByID(ctx context.Context, autoPlayID int) (types.AutoPlay, error)
	GetNextEncounterByStoryID(ctx context.Context, storyID int) (types.AutoPlayEncounterList, error)
	GetAutoPlayOptionsByChannelID(ctx context.Context, channelID, userID string) (types.AutoPlayOptions, error)
	GetAutoPlayActivities(ctx context.Context) ([]types.Activity, error)
	GetAnnounceByEncounterID(ctx context.Context, encounterID, autoPlayID int) (string, bool, error)
	GetNextEncounterByAutoPlayID(ctx context.Context, autoPlayID int) ([]types.Next, error)
	CreateAutoPlayTx(ctx context.Context, text string, storyID int, solo bool) (int, error)
	AddAutoPlayNext(ctx context.Context, next []types.Next) error
	CreateAutoPlayChannelTx(ctx context.Context, channelID, userID string, autoPlayID int) (int, error)
	RegisterActivitiesAutoPlay(ctx context.Context, autoPlayID, encounterID int, actions types.Actions) error
	UpdateProcessedAutoPlay(ctx context.Context, id int, processed bool, actions types.Actions) error
	UpdateAutoPlayGroup(ctx context.Context, id, count int, date time.Time) error
	UpdateAutoPlayState(ctx context.Context, autoPlayChannel string, encounterID int) error
	CloseAutoPlayChannel(ctx context.Context, channelID string, autoPlayID int) error
}

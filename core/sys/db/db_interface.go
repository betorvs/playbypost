package db

import (
	"context"
	"time"

	"github.com/betorvs/playbypost/core/initiative"
	"github.com/betorvs/playbypost/core/rpg"
	"github.com/betorvs/playbypost/core/rpg/base"
	"github.com/betorvs/playbypost/core/rules"
	"github.com/betorvs/playbypost/core/sys/web/types"
)

type DBClient interface {
	Close() error
	// Sessions
	CreateSession(ctx context.Context, session types.Session) error
	GetSessionByToken(ctx context.Context, token string) (types.Session, error)
	// GetSessionIDByToken(ctx context.Context, token string) (int64, error)
	UpdateSessionLastActivity(ctx context.Context, token string) error
	DeleteSessionByToken(ctx context.Context, token string) error
	DeleteSessionByID(ctx context.Context, sessionID int64) error
	DeleteExpiredSessions(ctx context.Context) error
	GetAllSessions(ctx context.Context) ([]types.Session, error)
	GetSessionByID(ctx context.Context, sessionID int64) (types.Session, error)
	// Session Events
	CreateSessionEvent(ctx context.Context, event types.SessionEvent) error
	GetSessionEvents(ctx context.Context) ([]types.SessionEvent, error)
	LogLoginAttempt(ctx context.Context, username, ipAddress, userAgent string, success bool) error
	LogSessionCreated(ctx context.Context, session types.Session, sessionID int64) error
	LogSessionDeleted(ctx context.Context, sessionID int64, reason string) error
	LogSessionExpired(ctx context.Context, sessionID int64) error
	LogLogout(ctx context.Context, sessionID int64, username string) error
	LogSessionValidated(ctx context.Context, sessionID int64, username string) error
	LogSessionInvalid(ctx context.Context, sessionID int64, reason string) error
	LogActivityUpdated(ctx context.Context, sessionID int64, username string) error
	LogCleanupExecuted(ctx context.Context, sessionsDeleted int) error
	// Writers
	GetWriters(ctx context.Context, active bool) ([]types.Writer, error)
	CreateWriters(ctx context.Context, username, password string) (int, error)
	CreateWriterUserAssociation(ctx context.Context, writerID, userID int) (int, error)
	GetUsersByWriterID(ctx context.Context, writerID int) ([]types.User, error)
	GetWritersByUserID(ctx context.Context, userID int) ([]types.Writer, error)
	CheckWriterExists(ctx context.Context, writerID int) (bool, error)
	CheckUserExists(ctx context.Context, userID int) (bool, error)
	CheckWriterUserAssociationExists(ctx context.Context, writerID, userID int) (bool, error)
	GetWriterUsersAssociation(ctx context.Context) ([]types.WriterUserAssociation, error)
	DeleteWriterUserAssociation(ctx context.Context, id int) error
	GetWriterByID(ctx context.Context, id int) (types.Writer, error)
	GetWriterByUsername(ctx context.Context, username string) (types.Writer, error)
	// story
	GetStory(ctx context.Context) ([]types.Story, error)
	CreateStoryTx(ctx context.Context, title, announcement, notes, encodingKey string, writerID int) (int, error)
	UpdateStoryTx(ctx context.Context, title, announcement, notes string, storyID int) (int, error)
	GetStoryIDByTitle(ctx context.Context, title string) (int, error)
	GetStoryByID(ctx context.Context, id int) (types.Story, error)
	GetStoriesByWriterID(ctx context.Context, id int) ([]types.Story, error)
	DeleteStoryByID(ctx context.Context, id int) error
	// Encounters
	GetEncounters(ctx context.Context) ([]types.Encounter, error)
	GetEncounterByStoryID(ctx context.Context, storyID int) ([]types.Encounter, error)
	GetEncounterByStoryIDWithPagination(ctx context.Context, storyID, limit, cursor int) ([]types.Encounter, int, int, error)
	GetEncounterByID(ctx context.Context, id int) (types.Encounter, error)
	CreateEncounterTx(ctx context.Context, title, announcement, notes string, storyID, storytellerID int, first, last bool) (int, error)
	UpdateEncounterTx(ctx context.Context, title, announcement, notes string, id, storyID int, first, last bool) (int, error)
	DeleteEncounterByID(ctx context.Context, id int) error
	// Tasks
	GetTask(ctx context.Context) ([]types.Task, error)
	GetTaskByID(ctx context.Context, id int) (types.Task, error)
	CreateTask(ctx context.Context, description, ability, skill string, kind types.TaskKind, target int) (int, error)
	UpdateTaskByID(ctx context.Context, description, ability, skill string, kind types.TaskKind, target, id int) error
	DeleteTaskByID(ctx context.Context, id int) error
	// GetTasksByEncounterID(ctx context.Context, id int) (map[string]types.Task, error)
	// stage
	GetStage(ctx context.Context) ([]types.Stage, error)
	GetStageByStoryID(ctx context.Context, id int) ([]types.Stage, error)
	GetStageByStageID(ctx context.Context, id int) (types.StageAggregated, error)
	GetStageEncounterByEncounterID(ctx context.Context, id int) (types.StageEncounter, error)
	GetStageEncountersByStageIDWithPagination(ctx context.Context, id, limit, cursor int) ([]types.StageEncounter, int, int, error)
	GetStageEncountersByStageID(ctx context.Context, id int) ([]types.StageEncounter, error)
	GetRunningStageByChannelID(ctx context.Context, channelID, userID string, rpgSystem *rpg.RPGSystem) (types.RunningStage, error)
	GetStageEncounterActivitiesByEncounterID(ctx context.Context, id int) ([]types.Activity, error)
	GetStageEncounterActivities(ctx context.Context) ([]types.Activity, error)
	GetStageTaskFromRunningTaskID(ctx context.Context, taskID int) (types.Task, error)
	GetCreatureFromParticipantsList(ctx context.Context, players []types.Options, npcs []types.Options, rpgSystem *rpg.RPGSystem) (map[int]rules.RolePlaying, map[int]rules.RolePlaying, error)
	GetNextEncounterByEncounterID(ctx context.Context, id int) (types.Next, error)
	GetNextEncounterByStageID(ctx context.Context, id int) ([]types.Next, error)
	GetStageEncounterListByStoryID(ctx context.Context, storyID int) (types.EncounterList, error)
	CreateStageTx(ctx context.Context, text, userid string, storyID, creatorID int) (int, error)
	AddChannelToStage(ctx context.Context, channel string, id int) (int, error)
	AddEncounterToStage(ctx context.Context, text string, stage_id, storyteller_id, encounter_id int) (int, error)
	UpdatePhase(ctx context.Context, id, phase int) error
	AddParticipants(ctx context.Context, encounterID int, npc bool, players []int) error
	AddNextEncounter(ctx context.Context, next []types.Next) error
	AddRunningTask(ctx context.Context, text string, stageID, taskID, StorytellerID, encounterID int) error
	RegisterActivities(ctx context.Context, stageID, encounterID int, actions types.Actions) error
	UpdateProcessedActivities(ctx context.Context, id int, processed bool, actions types.Actions) error
	CloseStage(ctx context.Context, id int) error
	DeleteStageNextEncounter(ctx context.Context, id int) error
	DeleteStageEncounterByID(ctx context.Context, id int) error
	// NPC
	GetNPCByStageID(ctx context.Context, id int) ([]types.Players, error)
	GenerateNPC(ctx context.Context, stageID, storytellerID int, creature *base.Creature, extension map[string]interface{}) (int, error)
	UpdateNPC(ctx context.Context, id int, creature *base.Creature, extension map[string]interface{}, destroyed bool) error
	// users
	GetUser(ctx context.Context) ([]types.User, error)
	GetUsersByID(ctx context.Context, id int) ([]types.User, error)
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
	GetPlayerByUserIDChannel(ctx context.Context, id, channel string, rpgSystem *rpg.RPGSystem) (types.Players, error)
	GetPlayerByUserID(ctx context.Context, id int, rpgSystem *rpg.RPGSystem) (types.Players, error)
	SavePlayerTx(ctx context.Context, id, storyID int, creature *base.Creature, extension map[string]interface{}) (int, error)
	UpdatePlayer(ctx context.Context, id int, creature *base.Creature, extension map[string]interface{}, destroyed bool) error
	UpdatePlayerDetails(ctx context.Context, id int, name, rpg string) error
	// Chat
	AddChatInformation(ctx context.Context, username, userid, channel, chat string) (int, error)
	GetChatInformation(ctx context.Context) ([]types.ChatInfo, error)
	GetChatChannelInformation(ctx context.Context) ([]string, error)
	GetChatRunningChannels(ctx context.Context, kind string) ([]types.RunningChannels, error)
	//  Auto Play
	GetAutoPlay(ctx context.Context) ([]types.AutoPlay, error)
	GetAutoPlayByID(ctx context.Context, autoPlayID int) (types.AutoPlay, error)
	GetAutoPlayEncounterListByStoryID(ctx context.Context, storyID int) (types.EncounterList, error)
	GetAutoPlayOptionsByChannelID(ctx context.Context, channelID, userID string) (types.AutoPlayOptions, error)
	GetAutoPlayActivities(ctx context.Context) ([]types.Activity, error)
	GetStoryAnnouncementByAutoPlayID(ctx context.Context, autoPlayID int) (string, string, error)
	GetAnnounceByEncounterID(ctx context.Context, encounterID, autoPlayID int) (string, bool, error)
	DescribeAutoPlayPublished(ctx context.Context, solo bool) ([]types.AutoPlayDescribed, error)
	GetNextEncounterByAutoPlayID(ctx context.Context, autoPlayID int) ([]types.Next, error)
	CreateAutoPlayTx(ctx context.Context, text string, storyID, creatorID int, solo bool) (int, error)
	AddAutoPlayNext(ctx context.Context, next []types.Next) error
	CreateAutoPlayChannelTx(ctx context.Context, channelID, userID string, autoPlayID int) (int, error)
	AddAutoPlayGroup(ctx context.Context, channelID, userID string) error
	RegisterActivitiesAutoPlay(ctx context.Context, autoPlayID, encounterID int, actions types.Actions) error
	UpdateProcessedAutoPlay(ctx context.Context, id int, processed bool, actions types.Actions) error
	UpdateAutoPlayGroup(ctx context.Context, id, count int, date time.Time) error
	UpdateAutoPlayState(ctx context.Context, autoPlayChannel string, encounterID int) error
	ChangePublishAutoPlay(ctx context.Context, autoPlayID int, publish bool) error
	CloseAutoPlayChannel(ctx context.Context, channelID string, autoPlayID int) error
	DeleteAutoPlayNextEncounter(ctx context.Context, id int) error
}

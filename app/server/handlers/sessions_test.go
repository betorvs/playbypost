package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"

	"github.com/betorvs/playbypost/core/sys/web/server"
	"github.com/betorvs/playbypost/core/sys/web/types"
	"github.com/betorvs/playbypost/core/tests/mock/dbclient"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"
)

// MockDBClient is a mock implementation of the DBClient interface

// type MockDBClient struct {
// 	mock.Mock
// }

// func (m *MockDBClient) Close() error {
// 	args := m.Called()
// 	return args.Error(0)
// }

// func (m *MockDBClient) CreateSession(ctx context.Context, session types.Session) error {
// 	args := m.Called(ctx, session)
// 	return args.Error(0)
// }

// func (m *MockDBClient) GetSessionByToken(ctx context.Context, token string) (types.Session, error) {
// 	args := m.Called(ctx, token)
// 	return args.Get(0).(types.Session), args.Error(1)
// }

// func (m *MockDBClient) DeleteSessionByToken(ctx context.Context, token string) error {
// 	args := m.Called(ctx, token)
// 	return args.Error(0)
// }

// func (m *MockDBClient) GetWriters(ctx context.Context, active bool) ([]types.Writer, error) {
// 	args := m.Called(ctx, active)
// 	return args.Get(0).([]types.Writer), args.Error(1)
// }

// func (m *MockDBClient) CreateWriters(ctx context.Context, username, password string) (int, error) {
// 	args := m.Called(ctx, username, password)
// 	return args.Int(0), args.Error(1)
// }

// func (m *MockDBClient) CreateWriterUserAssociation(ctx context.Context, writerID, userID int) (int, error) {
// 	args := m.Called(ctx, writerID, userID)
// 	return args.Int(0), args.Error(1)
// }

// func (m *MockDBClient) GetUsersByWriterID(ctx context.Context, writerID int) ([]types.User, error) {
// 	args := m.Called(ctx, writerID)
// 	return args.Get(0).([]types.User), args.Error(1)
// }

// func (m *MockDBClient) GetWritersByUserID(ctx context.Context, userID int) ([]types.Writer, error) {
// 	args := m.Called(ctx, userID)
// 	return args.Get(0).([]types.Writer), args.Error(1)
// }

// func (m *MockDBClient) CheckWriterExists(ctx context.Context, writerID int) (bool, error) {
// 	args := m.Called(ctx, writerID)
// 	return args.Bool(0), args.Error(1)
// }

// func (m *MockDBClient) CheckUserExists(ctx context.Context, userID int) (bool, error) {
// 	args := m.Called(ctx, userID)
// 	return args.Bool(0), args.Error(1)
// }

// func (m *MockDBClient) CheckWriterUserAssociationExists(ctx context.Context, writerID, userID int) (bool, error) {
// 	args := m.Called(ctx, writerID, userID)
// 	return args.Bool(0), args.Error(1)
// }

// func (m *MockDBClient) GetWriterUsersAssociation(ctx context.Context) ([]types.WriterUserAssociation, error) {
// 	args := m.Called(ctx)
// 	return args.Get(0).([]types.WriterUserAssociation), args.Error(1)
// }

// func (m *MockDBClient) DeleteWriterUserAssociation(ctx context.Context, id int) error {
// 	args := m.Called(ctx, id)
// 	return args.Error(0)
// }

// func (m *MockDBClient) GetWriterByID(ctx context.Context, id int) (types.Writer, error) {
// 	args := m.Called(ctx, id)
// 	return args.Get(0).(types.Writer), args.Error(1)
// }

// func (m *MockDBClient) GetWriterByUsername(ctx context.Context, username string) (types.Writer, error) {
// 	args := m.Called(ctx, username)
// 	return args.Get(0).(types.Writer), args.Error(1)
// }

// func (m *MockDBClient) GetStory(ctx context.Context) ([]types.Story, error) {
// 	args := m.Called(ctx)
// 	return args.Get(0).([]types.Story), args.Error(1)
// }

// func (m *MockDBClient) CreateStoryTx(ctx context.Context, title, announcement, notes, encodingKey string, writerID int) (int, error) {
// 	args := m.Called(ctx, title, announcement, notes, encodingKey, writerID)
// 	return args.Int(0), args.Error(1)
// }

// func (m *MockDBClient) UpdateStoryTx(ctx context.Context, title, announcement, notes string, storyID int) (int, error) {
// 	args := m.Called(ctx, title, announcement, notes, storyID)
// 	return args.Int(0), args.Error(1)
// }

// func (m *MockDBClient) GetStoryIDByTitle(ctx context.Context, title string) (int, error) {
// 	args := m.Called(ctx, title)
// 	return args.Int(0), args.Error(1)
// }

// func (m *MockDBClient) GetStoryByID(ctx context.Context, id int) (types.Story, error) {
// 	args := m.Called(ctx, id)
// 	return args.Get(0).(types.Story), args.Error(1)
// }

// func (m *MockDBClient) GetStoriesByWriterID(ctx context.Context, id int) ([]types.Story, error) {
// 	args := m.Called(ctx, id)
// 	return args.Get(0).([]types.Story), args.Error(1)
// }

// func (m *MockDBClient) DeleteStoryByID(ctx context.Context, id int) error {
// 	args := m.Called(ctx, id)
// 	return args.Error(0)
// }

// func (m *MockDBClient) GetEncounters(ctx context.Context) ([]types.Encounter, error) {
// 	args := m.Called(ctx)
// 	return args.Get(0).([]types.Encounter), args.Error(1)
// }

// func (m *MockDBClient) GetEncounterByStoryID(ctx context.Context, storyID int) ([]types.Encounter, error) {
// 	args := m.Called(ctx, storyID)
// 	return args.Get(0).([]types.Encounter), args.Error(1)
// }

// func (m *MockDBClient) GetEncounterByStoryIDWithPagination(ctx context.Context, storyID, limit, cursor int) ([]types.Encounter, int, int, error) {
// 	args := m.Called(ctx, storyID, limit, cursor)
// 	return args.Get(0).([]types.Encounter), args.Int(1), args.Int(2), args.Error(3)
// }

// func (m *MockDBClient) GetEncounterByID(ctx context.Context, id int) (types.Encounter, error) {
// 	args := m.Called(ctx, id)
// 	return args.Get(0).(types.Encounter), args.Error(1)
// }

// func (m *MockDBClient) CreateEncounterTx(ctx context.Context, title, announcement, notes string, storyID, storytellerID int, first, last bool) (int, error) {
// 	args := m.Called(ctx, title, announcement, notes, storyID, storytellerID, first, last)
// 	return args.Int(0), args.Error(1)
// }

// func (m *MockDBClient) UpdateEncounterTx(ctx context.Context, title, announcement, notes string, id, storyID int, first, last bool) (int, error) {
// 	args := m.Called(ctx, title, announcement, notes, id, storyID, first, last)
// 	return args.Int(0), args.Error(1)
// }

// func (m *MockDBClient) DeleteEncounterByID(ctx context.Context, id int) error {
// 	args := m.Called(ctx, id)
// 	return args.Error(0)
// }

// func (m *MockDBClient) GetTask(ctx context.Context) ([]types.Task, error) {
// 	args := m.Called(ctx)
// 	return args.Get(0).([]types.Task), args.Error(1)
// }

// func (m *MockDBClient) GetTaskByID(ctx context.Context, id int) (types.Task, error) {
// 	args := m.Called(ctx, id)
// 	return args.Get(0).(types.Task), args.Error(1)
// }

// func (m *MockDBClient) CreateTask(ctx context.Context, description, ability, skill string, kind types.TaskKind, target int) (int, error) {
// 	args := m.Called(ctx, description, ability, skill, kind, target)
// 	return args.Int(0), args.Error(1)
// }

// func (m *MockDBClient) UpdateTaskByID(ctx context.Context, description, ability, skill string, kind types.TaskKind, target, id int) error {
// 	args := m.Called(ctx, description, ability, skill, kind, target, id)
// 	return args.Error(0)
// }

// func (m *MockDBClient) DeleteTaskByID(ctx context.Context, id int) error {
// 	args := m.Called(ctx, id)
// 	return args.Error(0)
// }

// func (m *MockDBClient) GetStage(ctx context.Context) ([]types.Stage, error) {
// 	args := m.Called(ctx)
// 	return args.Get(0).([]types.Stage), args.Error(1)
// }

// func (m *MockDBClient) GetStageByStoryID(ctx context.Context, id int) ([]types.Stage, error) {
// 	args := m.Called(ctx, id)
// 	return args.Get(0).([]types.Stage), args.Error(1)
// }

// func (m *MockDBClient) GetStageByStageID(ctx context.Context, id int) (types.StageAggregated, error) {
// 	args := m.Called(ctx, id)
// 	return args.Get(0).(types.StageAggregated), args.Error(1)
// }

// func (m *MockDBClient) GetStageEncounterByEncounterID(ctx context.Context, id int) (types.StageEncounter, error) {
// 	args := m.Called(ctx, id)
// 	return args.Get(0).(types.StageEncounter), args.Error(1)
// }

// func (m *MockDBClient) GetStageEncountersByStageIDWithPagination(ctx context.Context, id, limit, cursor int) ([]types.StageEncounter, int, int, error) {
// 	args := m.Called(ctx, id, limit, cursor)
// 	return args.Get(0).([]types.StageEncounter), args.Int(1), args.Int(2), args.Error(3)
// }

// func (m *MockDBClient) GetStageEncountersByStageID(ctx context.Context, id int) ([]types.StageEncounter, error) {
// 	args := m.Called(ctx, id)
// 	return args.Get(0).([]types.StageEncounter), args.Error(1)
// }

// func (m *MockDBClient) GetRunningStageByChannelID(ctx context.Context, channelID, userID string, rpgSystem *rpg.RPGSystem) (types.RunningStage, error) {
// 	args := m.Called(ctx, channelID, userID, rpgSystem)
// 	return args.Get(0).(types.RunningStage), args.Error(1)
// }

// func (m *MockDBClient) GetStageEncounterActivitiesByEncounterID(ctx context.Context, id int) ([]types.Activity, error) {
// 	args := m.Called(ctx, id)
// 	return args.Get(0).([]types.Activity), args.Error(1)
// }

// func (m *MockDBClient) GetStageEncounterActivities(ctx context.Context) ([]types.Activity, error) {
// 	args := m.Called(ctx)
// 	return args.Get(0).([]types.Activity), args.Error(1)
// }

// func (m *MockDBClient) GetStageTaskFromRunningTaskID(ctx context.Context, taskID int) (types.Task, error) {
// 	args := m.Called(ctx, taskID)
// 	return args.Get(0).(types.Task), args.Error(1)
// }

// func (m *MockDBClient) GetCreatureFromParticipantsList(ctx context.Context, players []types.Options, npcs []types.Options, rpgSystem *rpg.RPGSystem) (map[int]rules.RolePlaying, map[int]rules.RolePlaying, error) {
// 	args := m.Called(ctx, players, npcs, rpgSystem)
// 	return args.Get(0).(map[int]rules.RolePlaying), args.Get(1).(map[int]rules.RolePlaying), args.Error(2)
// }

// func (m *MockDBClient) GetNextEncounterByEncounterID(ctx context.Context, id int) (types.Next, error) {
// 	args := m.Called(ctx, id)
// 	return args.Get(0).(types.Next), args.Error(1)
// }

// func (m *MockDBClient) GetNextEncounterByStageID(ctx context.Context, id int) ([]types.Next, error) {
// 	args := m.Called(ctx, id)
// 	return args.Get(0).([]types.Next), args.Error(1)
// }

// func (m *MockDBClient) GetStageEncounterListByStoryID(ctx context.Context, storyID int) (types.EncounterList, error) {
// 	args := m.Called(ctx, storyID)
// 	return args.Get(0).(types.EncounterList), args.Error(1)
// }

// func (m *MockDBClient) CreateStageTx(ctx context.Context, text, userid string, storyID, creatorID int) (int, error) {
// 	args := m.Called(ctx, text, userid, storyID, creatorID)
// 	return args.Int(0), args.Error(1)
// }

// func (m *MockDBClient) AddChannelToStage(ctx context.Context, channel string, id int) (int, error) {
// 	args := m.Called(ctx, channel, id)
// 	return args.Int(0), args.Error(1)
// }

// func (m *MockDBClient) AddEncounterToStage(ctx context.Context, text string, stage_id, storyteller_id, encounter_id int) (int, error) {
// 	args := m.Called(ctx, text, stage_id, storyteller_id, encounter_id)
// 	return args.Int(0), args.Error(1)
// }

// func (m *MockDBClient) UpdatePhase(ctx context.Context, id, phase int) error {
// 	args := m.Called(ctx, id, phase)
// 	return args.Error(0)
// }

// func (m *MockDBClient) AddParticipants(ctx context.Context, encounterID int, npc bool, players []int) error {
// 	args := m.Called(ctx, encounterID, npc, players)
// 	return args.Error(0)
// }

// func (m *MockDBClient) AddNextEncounter(ctx context.Context, next []types.Next) error {
// 	args := m.Called(ctx, next)
// 	return args.Error(0)
// }

// func (m *MockDBClient) AddRunningTask(ctx context.Context, text string, stageID, taskID, StorytellerID, encounterID int) error {
// 	args := m.Called(ctx, text, stageID, taskID, StorytellerID, encounterID)
// 	return args.Error(0)
// }

// func (m *MockDBClient) RegisterActivities(ctx context.Context, stageID, encounterID int, actions types.Actions) error {
// 	args := m.Called(ctx, stageID, encounterID, actions)
// 	return args.Error(0)
// }

// func (m *MockDBClient) UpdateProcessedActivities(ctx context.Context, id int, processed bool, actions types.Actions) error {
// 	args := m.Called(ctx, id, processed, actions)
// 	return args.Error(0)
// }

// func (m *MockDBClient) CloseStage(ctx context.Context, id int) error {
// 	args := m.Called(ctx, id)
// 	return args.Error(0)
// }

// func (m *MockDBClient) DeleteStageNextEncounter(ctx context.Context, id int) error {
// 	args := m.Called(ctx, id)
// 	return args.Error(0)
// }

// func (m *MockDBClient) DeleteStageEncounterByID(ctx context.Context, id int) error {
// 	args := m.Called(ctx, id)
// 	return args.Error(0)
// }

// func (m *MockDBClient) GetNPCByStageID(ctx context.Context, id int) ([]types.Players, error) {
// 	args := m.Called(ctx, id)
// 	return args.Get(0).([]types.Players), args.Error(1)
// }

// func (m *MockDBClient) GenerateNPC(ctx context.Context, stageID, storytellerID int, creature *base.Creature, extension map[string]interface{}) (int, error) {
// 	args := m.Called(ctx, stageID, storytellerID, creature, extension)
// 	return args.Int(0), args.Error(1)
// }

// func (m *MockDBClient) UpdateNPC(ctx context.Context, id int, creature *base.Creature, extension map[string]interface{}, destroyed bool) error {
// 	args := m.Called(ctx, id, creature, extension, destroyed)
// 	return args.Error(0)
// }

// func (m *MockDBClient) GetUser(ctx context.Context) ([]types.User, error) {
// 	args := m.Called(ctx)
// 	return args.Get(0).([]types.User), args.Error(1)
// }

// func (m *MockDBClient) GetUsersByID(ctx context.Context, id int) ([]types.User, error) {
// 	args := m.Called(ctx, id)
// 	return args.Get(0).([]types.User), args.Error(1)
// }

// func (m *MockDBClient) GetUserByUserID(ctx context.Context, id string) (types.User, error) {
// 	args := m.Called(ctx, id)
// 	return args.Get(0).(types.User), args.Error(1)
// }

// func (m *MockDBClient) CreateUserTx(ctx context.Context, userid string) (int, error) {
// 	args := m.Called(ctx, userid)
// 	return args.Int(0), args.Error(1)
// }

// func (m *MockDBClient) UpdateNextPlayer(ctx context.Context, id, nextPlayer int) error {
// 	args := m.Called(ctx, id, nextPlayer)
// 	return args.Error(0)
// }

// func (m *MockDBClient) SaveInitiativeTx(ctx context.Context, i initiative.Initiative, encounterID int) (int, error) {
// 	args := m.Called(ctx, i, encounterID)
// 	return args.Int(0), args.Error(1)
// }

// func (m *MockDBClient) GetInitiativeByID(ctx context.Context, id int) (initiative.Initiative, error) {
// 	args := m.Called(ctx, id)
// 	return args.Get(0).(initiative.Initiative), args.Error(1)
// }

// func (m *MockDBClient) GetRunningInitiativeByEncounterID(ctx context.Context, encounterID int) (initiative.Initiative, int, error) {
// 	args := m.Called(ctx, encounterID)
// 	return args.Get(0).(initiative.Initiative), args.Int(1), args.Error(2)
// }

// func (m *MockDBClient) DeactivateParticipant(ctx context.Context, id int, name string) (int, error) {
// 	args := m.Called(ctx, id, name)
// 	return args.Int(0), args.Error(1)
// }

// func (m *MockDBClient) GetPlayers(ctx context.Context, rpgSystem *rpg.RPGSystem) ([]types.Players, error) {
// 	args := m.Called(ctx, rpgSystem)
// 	return args.Get(0).([]types.Players), args.Error(1)
// }

// func (m *MockDBClient) GetPlayerByID(ctx context.Context, id int, rpgSystem *rpg.RPGSystem) (types.Players, error) {
// 	args := m.Called(ctx, id, rpgSystem)
// 	return args.Get(0).(types.Players), args.Error(1)
// }

// func (m *MockDBClient) GetPlayerByPlayerID(ctx context.Context, id int, rpgSystem *rpg.RPGSystem) (types.Players, error) {
// 	args := m.Called(ctx, id, rpgSystem)
// 	return args.Get(0).(types.Players), args.Error(1)
// }

// func (m *MockDBClient) GetPlayerByStageID(ctx context.Context, id int, rpgSystem *rpg.RPGSystem) ([]types.Players, error) {
// 	args := m.Called(ctx, id, rpgSystem)
// 	return args.Get(0).([]types.Players), args.Error(1)
// }

// func (m *MockDBClient) GetPlayerByUserIDChannel(ctx context.Context, id, channel string, rpgSystem *rpg.RPGSystem) (types.Players, error) {
// 	args := m.Called(ctx, id, channel, rpgSystem)
// 	return args.Get(0).(types.Players), args.Error(1)
// }

// func (m *MockDBClient) GetPlayerByUserID(ctx context.Context, id int, rpgSystem *rpg.RPGSystem) (types.Players, error) {
// 	args := m.Called(ctx, id, rpgSystem)
// 	return args.Get(0).(types.Players), args.Error(1)
// }

// func (m *MockDBClient) SavePlayerTx(ctx context.Context, id, storyID int, creature *base.Creature, extension map[string]interface{}) (int, error) {
// 	args := m.Called(ctx, id, storyID, creature, extension)
// 	return args.Int(0), args.Error(1)
// }

// func (m *MockDBClient) UpdatePlayer(ctx context.Context, id int, creature *base.Creature, extension map[string]interface{}, destroyed bool) error {
// 	args := m.Called(ctx, id, creature, extension, destroyed)
// 	return args.Error(0)
// }

// func (m *MockDBClient) UpdatePlayerDetails(ctx context.Context, id int, name, rpg string) error {
// 	args := m.Called(ctx, id, name, rpg)
// 	return args.Error(0)
// }

// func (m *MockDBClient) AddChatInformation(ctx context.Context, username, userid, channel, chat string) (int, error) {
// 	args := m.Called(ctx, username, userid, channel, chat)
// 	return args.Int(0), args.Error(1)
// }

// func (m *MockDBClient) GetChatInformation(ctx context.Context) ([]types.ChatInfo, error) {
// 	args := m.Called(ctx)
// 	return args.Get(0).([]types.ChatInfo), args.Error(1)
// }

// func (m *MockDBClient) GetChatChannelInformation(ctx context.Context) ([]string, error) {
// 	args := m.Called(ctx)
// 	return args.Get(0).([]string), args.Error(1)
// }

// func (m *MockDBClient) GetChatRunningChannels(ctx context.Context, kind string) ([]types.RunningChannels, error) {
// 	args := m.Called(ctx, kind)
// 	return args.Get(0).([]types.RunningChannels), args.Error(1)
// }

// func (m *MockDBClient) GetAutoPlay(ctx context.Context) ([]types.AutoPlay, error) {
// 	args := m.Called(ctx)
// 	return args.Get(0).([]types.AutoPlay), args.Error(1)
// }

// func (m *MockDBClient) GetAutoPlayByID(ctx context.Context, autoPlayID int) (types.AutoPlay, error) {
// 	args := m.Called(ctx, autoPlayID)
// 	return args.Get(0).(types.AutoPlay), args.Error(1)
// }

// func (m *MockDBClient) GetAutoPlayEncounterListByStoryID(ctx context.Context, storyID int) (types.EncounterList, error) {
// 	args := m.Called(ctx, storyID)
// 	return args.Get(0).(types.EncounterList), args.Error(1)
// }

// func (m *MockDBClient) GetAutoPlayOptionsByChannelID(ctx context.Context, channelID, userID string) (types.AutoPlayOptions, error) {
// 	args := m.Called(ctx, channelID, userID)
// 	return args.Get(0).(types.AutoPlayOptions), args.Error(1)
// }

// func (m *MockDBClient) GetAutoPlayActivities(ctx context.Context) ([]types.Activity, error) {
// 	args := m.Called(ctx)
// 	return args.Get(0).([]types.Activity), args.Error(1)
// }

// func (m *MockDBClient) GetStoryAnnouncementByAutoPlayID(ctx context.Context, autoPlayID int) (string, string, error) {
// 	args := m.Called(ctx, autoPlayID)
// 	return args.String(0), args.String(1), args.Error(2)
// }

// func (m *MockDBClient) GetAnnounceByEncounterID(ctx context.Context, encounterID, autoPlayID int) (string, bool, error) {
// 	args := m.Called(ctx, encounterID, autoPlayID)
// 	return args.String(0), args.Bool(1), args.Error(2)
// }

// func (m *MockDBClient) DescribeAutoPlayPublished(ctx context.Context, solo bool) ([]types.AutoPlayDescribed, error) {
// 	args := m.Called(ctx, solo)
// 	return args.Get(0).([]types.AutoPlayDescribed), args.Error(1)
// }

// func (m *MockDBClient) GetNextEncounterByAutoPlayID(ctx context.Context, autoPlayID int) ([]types.Next, error) {
// 	args := m.Called(ctx, autoPlayID)
// 	return args.Get(0).([]types.Next), args.Error(1)
// }

// func (m *MockDBClient) CreateAutoPlayTx(ctx context.Context, text string, storyID, creatorID int, solo bool) (int, error) {
// 	args := m.Called(ctx, text, storyID, creatorID, solo)
// 	return args.Int(0), args.Error(1)
// }

// func (m *MockDBClient) AddAutoPlayNext(ctx context.Context, next []types.Next) error {
// 	args := m.Called(ctx, next)
// 	return args.Error(0)
// }

// func (m *MockDBClient) CreateAutoPlayChannelTx(ctx context.Context, channelID, userID string, autoPlayID int) (int, error) {
// 	args := m.Called(ctx, channelID, userID, autoPlayID)
// 	return args.Int(0), args.Error(1)
// }

// func (m *MockDBClient) AddAutoPlayGroup(ctx context.Context, channelID, userID string) error {
// 	args := m.Called(ctx, channelID, userID)
// 	return args.Error(0)
// }

// func (m *MockDBClient) RegisterActivitiesAutoPlay(ctx context.Context, autoPlayID, encounterID int, actions types.Actions) error {
// 	args := m.Called(ctx, autoPlayID, encounterID, actions)
// 	return args.Error(0)
// }

// func (m *MockDBClient) UpdateProcessedAutoPlay(ctx context.Context, id int, processed bool, actions types.Actions) error {
// 	args := m.Called(ctx, id, processed, actions)
// 	return args.Error(0)
// }

// func (m *MockDBClient) UpdateAutoPlayGroup(ctx context.Context, id, count int, date time.Time) error {
// 	args := m.Called(ctx, id, count, date)
// 	return args.Error(0)
// }

// func (m *MockDBClient) UpdateAutoPlayState(ctx context.Context, autoPlayChannel string, encounterID int) error {
// 	args := m.Called(ctx, autoPlayChannel, encounterID)
// 	return args.Error(0)
// }

// func (m *MockDBClient) ChangePublishAutoPlay(ctx context.Context, autoPlayID int, publish bool) error {
// 	args := m.Called(ctx, autoPlayID, publish)
// 	return args.Error(0)
// }

// func (m *MockDBClient) CloseAutoPlayChannel(ctx context.Context, channelID string, autoPlayID int) error {
// 	args := m.Called(ctx, channelID, autoPlayID)
// 	return args.Error(0)
// }

// func (m *MockDBClient) DeleteAutoPlayNextEncounter(ctx context.Context, id int) error {
// 	args := m.Called(ctx, id)
// 	return args.Error(0)
// }

// func (m *MockDBClient) CreateSessionEvent(ctx context.Context, event types.SessionEvent) error {
// 	args := m.Called(ctx, event)
// 	return args.Error(0)
// }

// func (m *MockDBClient) DeleteExpiredSessions(ctx context.Context) error {
// 	args := m.Called(ctx)
// 	return args.Error(0)
// }

// func (m *MockDBClient) GetSessionEvents(ctx context.Context) ([]types.SessionEvent, error) {
// 	args := m.Called(ctx)
// 	return args.Get(0).([]types.SessionEvent), args.Error(1)
// }

func TestSignin(t *testing.T) {
	mockDB := new(dbclient.MockDBClient)
	s := &Session{
		db:     mockDB,
		s:      server.NewServer(0, slog.Default()),
		ctx:    context.Background(),
		logger: slog.Default(),
		Sessions: Sessions{
			Current: make(map[string]types.Session),
			mu:      &sync.Mutex{},
		},
	}

	creds := types.Credentials{Username: "testuser", Password: "password"}
	body, _ := json.Marshal(creds)
	req := httptest.NewRequest("POST", "/login", bytes.NewBuffer(body))
	rr := httptest.NewRecorder()

	// Generate a proper bcrypt hash for "password"
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password"), bcrypt.DefaultCost)
	writer := types.Writer{ID: 1, Username: "testuser", Password: string(hashedPassword)}
	mockDB.On("GetWriterByUsername", mock.Anything, "testuser").Return(writer, nil)
	mockDB.On("CreateSession", mock.Anything, mock.AnythingOfType("types.Session")).Return(nil)
	mockDB.On("LogLoginAttempt", mock.Anything, "testuser", "192.0.2.1:1234", "", true).Return(nil)

	s.Signin(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	mockDB.AssertExpectations(t)
}

func TestLogout(t *testing.T) {
	mockDB := new(dbclient.MockDBClient)
	s := &Session{
		db:     mockDB,
		s:      server.NewServer(0, slog.Default()),
		ctx:    context.Background(),
		logger: slog.Default(),
		Sessions: Sessions{
			Current: make(map[string]types.Session),
			mu:      &sync.Mutex{},
		},
	}

	req := httptest.NewRequest("POST", "/logout", nil)
	req.Header.Set("X-Access-Token", "testtoken")
	rr := httptest.NewRecorder()

	mockDB.On("DeleteSessionByToken", mock.Anything, "testtoken").Return(nil)

	s.Logout(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	mockDB.AssertExpectations(t)
}

func TestSessionExpiry(t *testing.T) {
	// Test expired session (expiry in the past)
	expiredSession := types.Session{
		Username: "testuser",
		Token:    "expired-token",
		Expiry:   time.Now().Add(-1 * time.Hour), // 1 hour ago
	}
	assert.True(t, expiredSession.IsExpired(), "Session should be expired")

	// Test valid session (expiry in the future)
	validSession := types.Session{
		Username: "testuser",
		Token:    "valid-token",
		Expiry:   time.Now().Add(1 * time.Hour), // 1 hour from now
	}
	assert.False(t, validSession.IsExpired(), "Session should not be expired")

	// Test session expiring right now
	nowSession := types.Session{
		Username: "testuser",
		Token:    "now-token",
		Expiry:   time.Now(),
	}
	assert.True(t, nowSession.IsExpired(), "Session should be expired when expiry is now")
}

func TestSessionTimezoneHandling(t *testing.T) {
	// Test that timezone handling works correctly
	now := time.Now()

	// Create a session with current time
	session := types.Session{
		Username:     "testuser",
		Token:        "test-token",
		Expiry:       now.Add(1 * time.Hour), // 1 hour from now
		UserID:       1,
		ClientType:   "web",
		ClientInfo:   "test",
		IPAddress:    "127.0.0.1",
		UserAgent:    "test",
		CreatedAt:    now,
		UpdatedAt:    now,
		LastActivity: now,
	}

	// The session should not be expired
	assert.False(t, session.IsExpired(), "Session should not be expired when expiry is 1 hour in the future")

	// Test with expired session
	expiredSession := types.Session{
		Username:     "testuser",
		Token:        "expired-token",
		Expiry:       now.Add(-1 * time.Hour), // 1 hour ago
		UserID:       1,
		ClientType:   "web",
		ClientInfo:   "test",
		IPAddress:    "127.0.0.1",
		UserAgent:    "test",
		CreatedAt:    now,
		UpdatedAt:    now,
		LastActivity: now,
	}

	// The session should be expired
	assert.True(t, expiredSession.IsExpired(), "Session should be expired when expiry is 1 hour in the past")
}

func TestGetClientContext(t *testing.T) {
	// Test with a web client
	req := httptest.NewRequest("POST", "/login", nil)
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/108.0.0.0 Safari/537.36")
	req.Header.Set("Referer", "http://localhost:3000/")
	clientType, clientInfo := getClientContext(req)
	assert.Equal(t, "web", clientType)
	var clientInfoMap map[string]string
	_ = json.Unmarshal([]byte(clientInfo), &clientInfoMap)
	assert.Equal(t, "web", clientInfoMap["kind"])

	// Test with a CLI client
	req = httptest.NewRequest("POST", "/login", nil)
	req.Header.Set("User-Agent", "Go-http-client/1.1")
	clientType, clientInfo = getClientContext(req)
	assert.Equal(t, "cli", clientType)
	_ = json.Unmarshal([]byte(clientInfo), &clientInfoMap)
	assert.Equal(t, "cli", clientInfoMap["kind"])
	assert.Equal(t, "local", clientInfoMap["source"])

	// Test with an unknown client
	req = httptest.NewRequest("POST", "/login", nil)
	req.Header.Set("User-Agent", "some-other-client")
	clientType, clientInfo = getClientContext(req)
	assert.Equal(t, "unknown", clientType)
	_ = json.Unmarshal([]byte(clientInfo), &clientInfoMap)
	assert.Equal(t, "unknown", clientInfoMap["kind"])
	assert.Equal(t, "some-other-client", clientInfoMap["userAgent"])
}

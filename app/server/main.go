/*
Copyright © 2024 Roberto Scudeller <beto.rvs@gmail.com>
*/
package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	v1 "github.com/betorvs/playbypost/app/server/handlers/v1"
	"github.com/betorvs/playbypost/app/server/worker"
	"github.com/betorvs/playbypost/core/rpg"
	"github.com/betorvs/playbypost/core/sys/db"
	"github.com/betorvs/playbypost/core/sys/library"
	"github.com/betorvs/playbypost/core/sys/scheduler"
	"github.com/betorvs/playbypost/core/sys/web/cli"
	"github.com/betorvs/playbypost/core/sys/web/server"
	"github.com/betorvs/playbypost/core/utils"
)

const (
	adminUser string = "admin"
)

var (
	Version string = "development"
	Port    int    = 3000 // 3000
	// ExternalPort int    = 8090 // 8090
)

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	slog.SetDefault(logger)
	logger.Info("starting admin web server", "version", Version, "port", Port)
	ctx := context.Background()

	conn := utils.LoadDBEnvVars()
	db := db.NewDB(conn, logger)
	defer func() {
		err := db.Close()
		if err != nil {
			logger.Error("error closing db", "error", err)
		}
	}()
	// rpg loading
	rpgFlag := flag.String("rpg", rpg.D10HM, fmt.Sprintf("-rpg [%s|%s]", rpg.D10HM, rpg.PFD20))
	stageWorker := flag.Bool("stage-worker", false, "-stage-worker")
	autoPlayWorker := flag.Bool("autoplay-worker", false, "-autoplay-worker")
	flag.Parse()
	rpgSystem := loadRPG(*rpgFlag, logger)
	d := rpg.Roll{
		RPGSystem: rpgSystem,
	}
	logger.Info("rpg system loaded", "system", rpgSystem)
	// creating library
	lib := library.New()
	lib.ImportRPGLibrary(rpgSystem.Name, logger, utils.LoadLibraryFiles())
	logger.Info("library loaded", "library", lib)
	// creating a random token for admin user
	adminToken := utils.RandomString(12)

	// web server
	srv := server.NewServer(Port, logger)
	// register health status
	srv.Register("GET /health", srv.GetHealth)

	client := cli.NewHeaders("http://localhost:8091", adminUser, adminToken)
	app := v1.NewMainApi(adminUser, ctx, d, db, logger, srv, client, rpgSystem, lib)

	srv.RegisterStatic()

	// sessions
	// srv.Register("POST /login", app.Session.Signin)
	// srv.Register("POST /logoff", app.Session.Logout)
	// srv.Register("POST /refresh", app.Session.Refresh)
	srv.Register("GET /api/v1/validate", app.Session.ValidateSession)
	srv.Register("POST /api/v1/login", app.Session.Signin)
	srv.Register("POST /api/v1/logout", app.Session.Logout)

	// session events
	srv.Register("GET /api/v1/session/events", app.GetSessionEvents)
	srv.Register("GET /api/v1/session/active", app.GetActiveSessions)
	srv.Register("GET /api/v1/session", app.Session.GetAllSessions)

	// writers
	srv.Register("GET /api/v1/writer", app.GetWriters)
	srv.Register("POST /api/v1/writer", app.CreateWriters)
	srv.Register("GET /api/v1/writer/association", app.GetWriterUsersAssociation)
	srv.Register("POST /api/v1/writer/user", app.CreateWriterUserAssociation)
	srv.Register("DELETE /api/v1/writer/association/{id}", app.DeleteWriterUserAssociation)

	// users
	srv.Register("GET /api/v1/user/{id}", app.GetUsersByUserID)

	// characters
	srv.Register("GET /api/v1/characters", app.GetCharacters)
	srv.Register("PUT /api/v1/characters/{id}", app.UpdateCharacter)

	// story
	srv.Register("GET /api/v1/story", app.GetStory)
	srv.Register("GET /api/v1/story/{id}", app.GetStoryById)
	srv.Register("GET /api/v1/story/writer/{id}", app.GetStoryByWriterId)
	srv.Register("POST /api/v1/story", app.CreateStory)
	srv.Register("PUT /api/v1/story/{id}", app.UpdateStoryByID)
	srv.Register("DELETE /api/v1/story/{id}", app.DeleteStoryByID)

	// stage
	srv.Register("GET /api/v1/stage", app.GetStage)
	srv.Register("GET /api/v1/stage/{id}", app.GetStageById)
	srv.Register("GET /api/v1/stage/story/{id}", app.GetStageByStoryId)
	srv.Register("POST /api/v1/stage", app.CreateStage)
	srv.Register("GET /api/v1/stage/player/{id}", app.GetPlayersByStageID)
	srv.Register("GET /api/v1/stage/npc/{id}", app.GetNPCByStageID)
	srv.Register("POST /api/v1/stage/npc", app.GenerateNPC)
	srv.Register("POST /api/v1/stage/channel", app.AddChannelToStage)
	srv.Register("GET /api/v1/stage/encounter/{id}", app.GetStageEncounterByEncounterID)
	srv.Register("GET /api/v1/stage/encounters/{id}", app.GetStageEncountersByStageID)
	srv.Register("PUT /api/v1/stage/encounter/{id}/{phase}", app.UpdateEncounterPhaseById)
	srv.Register("POST /api/v1/stage/encounter", app.AddEncounterToStage)
	srv.Register("DELETE /api/v1/stage/encounter/{id}", app.DeleteStageEncounterByID)
	srv.Register("POST /api/v1/stage/encounter/participants", app.AddParticipants)
	srv.Register("GET /api/v1/stage/encounter/story/{id}", app.GetStageEncounterListByStoryID)
	srv.Register("GET /api/v1/stage/encounter/next/{id}", app.GetNextEncounterByStageID)
	srv.Register("POST /api/v1/stage/encounter/next", app.AddNextEncounter)
	srv.Register("DELETE /api/v1/stage/encounter/next/{id}", app.DeleteStageNextEncounter)
	srv.Register("POST /api/v1/stage/encounter/task", app.AddRunningTask)
	srv.Register("GET /api/v1/stage/encounter/activities", app.GetStageEncounterActivities)
	srv.Register("GET /api/v1/stage/encounter/activities/{id}", app.GetStageEncounterActivitiesByEncounterID)
	srv.Register("PUT /api/v1/stage/{id}", app.CloseStage)

	// players
	srv.Register("POST /api/v1/player", app.GeneratePlayer)
	srv.Register("GET /api/v1/player", app.GetPlayers)
	srv.Register("GET /api/v1/player/{id}", app.GetPlayersByID)
	// srv.Register("GET /api/v1/player/stage/{id}", app.GetPlayersByStageID)
	// srv.Register("OPTIONS /api/v1/player/story/{id}", srv.Options)

	// encounters
	srv.Register("GET /api/v1/encounter", app.GetEncounters)
	srv.Register("GET /api/v1/encounter/{id}", app.GetEncounterById)
	srv.Register("GET /api/v1/encounter/story/{id}", app.GetEncounterByStoryId)
	srv.Register("POST /api/v1/encounter", app.CreateEncounter)
	srv.Register("PUT /api/v1/encounter/{id}", app.UpdateEncounterByID)
	srv.Register("DELETE /api/v1/encounter/{id}", app.DeleteEncounterByID)
	// tasks
	srv.Register("GET /api/v1/task", app.GetTask)
	srv.Register("GET /api/v1/task/{id}", app.GetTaskByID)
	srv.Register("POST /api/v1/task", app.CreateTasks)
	srv.Register("PUT /api/v1/task/{id}", app.UpdateTaskByID)
	srv.Register("DELETE /api/v1/task/{id}", app.DeleteTaskByID)

	// initiative
	srv.Register("POST /api/v1/initiative", app.GenerateInitiative)
	srv.Register("GET /api/v1/initiative/encounter/{id}", app.GetInitiativeByEncounterId)
	// command api
	srv.Register("POST /api/v1/command", app.ExecuteCommand)
	// chat api
	srv.Register("POST /api/v1/info", app.AddChatInfo)
	srv.Register("GET /api/v1/info/users", app.GetUsersInformation)
	srv.Register("GET /api/v1/info/channel", app.GetChannelsInformation)
	srv.Register("GET /api/v1/info/channel/{kind}", app.GetRunningChannels)
	srv.Register("GET /api/v1/info/phases", app.GetEncountersPhase)

	// auto play
	srv.Register("GET /api/v1/autoplay", app.GetAutoPlay)
	srv.Register("GET /api/v1/autoplay/{id}", app.GetAutoPlayByID)
	srv.Register("GET /api/v1/autoplay/encounter/story/{id}", app.GetAutoPlayEncounterListByStoryID)
	srv.Register("POST /api/v1/autoplay", app.CreateAutoPlay)
	srv.Register("PUT /api/v1/autoplay/publish/{id}", app.ChangePublishFlagAutoPlay)
	srv.Register("POST /api/v1/autoplay/next", app.AddAutoPlayNext)
	srv.Register("GET /api/v1/autoplay/next/{id}", app.GetAutoPlayNextEncounterByAutoPlayID)
	srv.Register("DELETE /api/v1/autoplay/next/{id}", app.DeleteAutoPlayNextEncounter)

	// validator auto play and stage and story
	srv.Register("GET /api/v1/validator/autoplay/{hashid}", app.GetValidateAutoPlay)
	srv.Register("GET /api/v1/validator/stage/{hashid}", app.GetValidateStage)
	srv.Register("GET /api/v1/validator/story/{hashid}", app.GetValidateStory)
	srv.Register("PUT /api/v1/validator/autoplay/{id}", app.RequestToValidateAutoPlay)
	srv.Register("PUT /api/v1/validator/stage/{id}", app.RequestToValidateStage)
	srv.Register("PUT /api/v1/validator/story/{id}", app.RequestToValidateStory)
	srv.Register("GET /api/v1/validator", app.GetAllValidations)
	// options
	srv.Register("OPTIONS /", srv.Options)

	app.Session.AddAdminSession(adminUser, adminToken)
	logger.Info("adding admin user one year token", "admin", adminUser, "token", adminToken)
	adminFile := utils.GetEnv("CREDS_FILE", "./creds")
	err := utils.Save(adminToken, adminFile)
	if err != nil {
		logger.Error("cannot write to file", "error", err)
		os.Exit(1)
	}

	if *stageWorker {
		app.Worker.StageActive = true
		logger.Info("starting stage worker job scheduler")
	}
	if *autoPlayWorker {
		app.Worker.AutoPlayActive = true
		logger.Info("starting auto play worker job scheduler")
	}

	jobScheduler := scheduler.NewJobScheduler(10 * time.Second)
	jobScheduler.JobQueue = app.Worker
	ctxJob, jobCancel := context.WithCancel(ctx)
	defer jobCancel()

	go func() {
		jobScheduler.Start(ctxJob)
	}()

	jobSchedulerHourly := scheduler.NewJobScheduler(1 * time.Hour)
	jobSchedulerHourly.JobQueue = app.Validator
	ctxJobHourly, jobCancelHourly := context.WithCancel(ctx)
	defer jobCancelHourly()

	go func() {
		jobSchedulerHourly.Start(ctxJobHourly)
	}()

	cleanupWorker := worker.NewCleanup(logger, db, ctx)
	cleanupWorker.Start()

	// starting a goroutine to server http requests
	// start web server in a goroutine
	go func() {
		if err := srv.Server.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			logger.Error("server listen and serve error", "error", err)
			os.Exit(1)
		}
		logger.Info("server stopped serving new connections.")
	}()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	ctxTimeout, ctxCancel := context.WithTimeout(ctx, 10*time.Second)
	defer ctxCancel()

	if err := srv.Server.Shutdown(ctxTimeout); err != nil {
		logger.Error("server shutdown error", "error", err)
	}
	logger.Info("graceful shutdown complete.")
}

func loadRPG(k string, logger *slog.Logger) *rpg.RPGSystem {
	logger.Info("loading rpg system", "system", k)
	var rpgSystem rpg.RPGSystem
	switch k {
	case rpg.D10HM:
		rpgSystem = *rpg.LoadRPGSystemsDefault(rpg.D10HM)

	case rpg.PFD20:
		rpgSystem = *rpg.LoadRPGSystemsDefault(rpg.PFD20)

	}
	return &rpgSystem
}

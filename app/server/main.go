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
	"github.com/betorvs/playbypost/core/rpg"
	"github.com/betorvs/playbypost/core/sys/db"
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
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	slog.SetDefault(logger)
	logger.Info("starting admin web server", "version", Version, "port", Port)
	ctx := context.Background()

	conn := utils.LoadDBEnvVars()
	db := db.NewDB(conn, logger)
	defer db.Close()
	// rpg loading
	rpgFlag := flag.String("rpg", rpg.D10HM, fmt.Sprintf("-rpg [%s|%s|%s|%s|%s]", rpg.D10HM, rpg.D10OS, rpg.D2035, rpg.Solo, rpg.Didactic))
	flag.Parse()
	rpgSystem := loadRPG(*rpgFlag, logger)
	d := rpg.Roll{
		RPGSystem: rpgSystem,
	}
	logger.Info("rpg system loaded", "system", rpgSystem, "ability", rpgSystem.Ability, "skill", rpgSystem.Skill)
	// creating a random token for admin user
	adminToken := utils.RandomString(12)

	// web server
	srv := server.NewServer(Port, logger)
	// register health status
	srv.Register("GET /health", srv.GetHealth)

	client := cli.NewHeaders("http://localhost:8091", adminUser, adminToken)
	app := v1.NewMainApi(ctx, d, db, logger, srv, client, rpgSystem)

	srv.RegisterStatic()

	srv.Register("POST /login", app.Session.Signin)
	// srv.Register("OPTIONS /login", srv.Options)
	srv.Register("POST /logoff", app.Session.Logout)
	srv.Register("POST /refresh", app.Session.Refresh)

	srv.Register("GET /api/v1/writer", app.GetWriters)
	// srv.Register("OPTIONS /api/v1/writer", srv.Options)
	srv.Register("POST /api/v1/writer", app.CreateWriters)

	srv.Register("GET /api/v1/story", app.GetStory)
	srv.Register("GET /api/v1/story/{id}", app.GetStoryById)
	// srv.Register("OPTIONS /api/v1/story/{id}", srv.Options)
	srv.Register("GET /api/v1/story/writer/{id}", app.GetStoryByWriterId)
	// srv.Register("OPTIONS /api/v1/story/writer/{id}", srv.Options)
	srv.Register("POST /api/v1/story", app.CreateStory)
	// srv.Register("OPTIONS /api/v1/story", srv.Options)

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
	srv.Register("GET /api/v1/stage/encounters/{id}", app.GetStageEncounterByStageID)
	srv.Register("PUT /api/v1/stage/encounter/{id}/{phase}", app.UpdateEncounterPhaseById)
	srv.Register("POST /api/v1/stage/encounter", app.AddEncounterToStage)
	srv.Register("POST /api/v1/stage/encounter/participants", app.AddParticipants)
	// stage_next_encounter
	srv.Register("POST /api/v1/stage/encounter/next", app.AddNextEncounter)
	// stage_running_tasks
	srv.Register("POST /api/v1/stage/encounter/task", app.AddRunningTask)
	// stage_encounter_activities
	srv.Register("GET /api/v1/stage/encounter/activities", app.GetStageEncounterActivities)
	srv.Register("GET /api/v1/stage/encounter/activities/{id}", app.GetStageEncounterActivitiesByEncounterID)

	srv.Register("POST /api/v1/player", app.GeneratePlayer)
	srv.Register("GET /api/v1/player", app.GetPlayers)
	srv.Register("GET /api/v1/player/{id}", app.GetPlayersByID)
	// srv.Register("GET /api/v1/player/stage/{id}", app.GetPlayersByStageID)
	// srv.Register("OPTIONS /api/v1/player/story/{id}", srv.Options)

	srv.Register("GET /api/v1/encounter", app.GetEncounters)
	srv.Register("GET /api/v1/encounter/{id}", app.GetEncounterById)
	srv.Register("GET /api/v1/encounter/story/{id}", app.GetEncounterByStoryId)
	// srv.Register("OPTIONS /api/v1/encounter/story/{id}", srv.Options)
	srv.Register("POST /api/v1/encounter", app.CreateEncounter)
	//
	srv.Register("GET /api/v1/task", app.GetTask)
	srv.Register("POST /api/v1/task", app.CreateTasks)

	// srv.Register("OPTIONS /api/v1/encounter", srv.Options)
	// srv.Register("PUT /api/v1/encounter/{id}/{phase}", app.UpdateEncounterPhaseById)

	srv.Register("POST /api/v1/initiative", app.GenerateInitiative)
	srv.Register("GET /api/v1/initiative/encounter/{id}", app.GetInitiativeByEncounterId)
	// command api
	srv.Register("POST /api/v1/command", app.ExecuteCommand)
	srv.Register("POST /api/v1/info", app.AddSlackInfo)
	srv.Register("GET /api/v1/info/users", app.GetUsersInformation)
	srv.Register("GET /api/v1/info/channel", app.GetChannelsInformation)
	srv.Register("GET /api/v1/info/phases", app.GetEncountersPhase)

	srv.Register("OPTIONS /*", srv.Options)

	app.Session.AddAdminSession(adminUser, adminToken)
	logger.Info("adding admin user one year token", "admin", adminUser, "token", adminToken)
	adminFile := utils.GetEnv("CREDS_FILE", "./creds")
	err := save(adminToken, adminFile)
	if err != nil {
		logger.Error("cannot write to file", "error", err)
		os.Exit(1)
	}

	jobScheduler := scheduler.NewJobScheduler(10 * time.Second)
	logger.Info("starting job scheduler")
	jobScheduler.JobQueue = app.Worker
	ctxJob, jobCancel := context.WithCancel(ctx)
	defer jobCancel()

	go func() {
		jobScheduler.Start(ctxJob)
	}()

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
		rpgSystem = *rpg.LoadRPGSystemsDefault(rpg.D10HM, logger)
		rpgSystem.InitDefinitions("./library/definitions-d10HM.json", logger)

	case rpg.D2035:
		rpgSystem = *rpg.LoadRPGSystemsDefault(rpg.D2035, logger)
		rpgSystem.InitDefinitions("./library/definitions-d20.json", logger)

	case rpg.D10OS:
		rpgSystem = *rpg.LoadRPGSystemsDefault(rpg.D10OS, logger)
		rpgSystem.InitDefinitions("./library/definitions-d10OS.json", logger)

	case rpg.Solo:
		rpgSystem = *rpg.LoadRPGSystemsDefault(rpg.Solo, logger)

	case rpg.Didactic:
		rpgSystem = *rpg.LoadRPGSystemsDefault(rpg.Didactic, logger)
	}
	return &rpgSystem
}

func save(value, file string) error {
	f, err := os.Create(file)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = f.WriteString(value)
	if err != nil {
		return err
	}
	return nil
}

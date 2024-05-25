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

	"github.com/betorvs/playbypost/app/server/handlers"
	"github.com/betorvs/playbypost/core/rpg"
	"github.com/betorvs/playbypost/core/sys/db"
	"github.com/betorvs/playbypost/core/sys/web/server"
	"github.com/betorvs/playbypost/core/utils"
)

var (
	Version      string = "development"
	Port         int    = 3000 // 3000
	ExternalPort int    = 8090 // 8090
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
	rpgFlag := flag.String("rpg", rpg.D10HM, fmt.Sprintf("-rpg [%s|%s|%s]", rpg.D10HM, rpg.D10OS, rpg.D2035))
	flag.Parse()
	rpgSystem := loadRPG(*rpgFlag, logger)
	d := rpg.Roll{
		RPGSystem: rpgSystem,
	}

	// web server
	srv := server.NewServer(Port, logger)
	// register health status
	srv.Register("GET /health", srv.GetHealth)

	app := handlers.NewMainApi(ctx, d, db, logger, srv, rpgSystem)

	srv.RegisterStatic()

	srv.Register("POST /login", app.Signin)
	// srv.Register("OPTIONS /login", srv.Options)
	srv.Register("POST /logoff", app.Logout)
	srv.Register("POST /refresh", app.Refresh)

	srv.Register("GET /api/v1/storyteller", app.GetStorytellers)
	// srv.Register("OPTIONS /api/v1/storyteller", srv.Options)
	srv.Register("GET /api/v1/storyteller/card", app.GetUsersCard)
	srv.Register("POST /api/v1/storyteller", app.CreateStorytellers)

	srv.Register("GET /api/v1/story", app.GetStory)
	srv.Register("GET /api/v1/story/{id}", app.GetStoryById)
	// srv.Register("OPTIONS /api/v1/story/{id}", srv.Options)
	srv.Register("GET /api/v1/story/master/{id}", app.GetStoryByMasterId)
	// srv.Register("OPTIONS /api/v1/story/master/{id}", srv.Options)
	srv.Register("POST /api/v1/story", app.CreateStory)
	// srv.Register("OPTIONS /api/v1/story", srv.Options)

	srv.Register("POST /api/v1/player", app.GeneratePlayer)
	srv.Register("GET /api/v1/player/{id}", app.GetPlayersByID)
	srv.Register("GET /api/v1/player/story/{id}", app.GetPlayersByStoryID)
	// srv.Register("OPTIONS /api/v1/player/story/{id}", srv.Options)

	srv.Register("GET /api/v1/encounter", app.GetEncounters)
	srv.Register("GET /api/v1/encounter/{id}", app.GetEncounterById)
	srv.Register("GET /api/v1/encounter/story/{id}", app.GetEncounterByStoryId)
	// srv.Register("OPTIONS /api/v1/encounter/story/{id}", srv.Options)
	srv.Register("POST /api/v1/encounter", app.CreateEncounter)
	srv.Register("POST /api/v1/encounter/task", app.CreateTasks)
	// srv.Register("OPTIONS /api/v1/encounter", srv.Options)
	// srv.Register("PUT /api/v1/encounter/{id}/{phase}", app.UpdateEncounterPhaseById)
	srv.Register("POST /api/v1/encounter/participants", app.AddParticipants)
	srv.Register("POST /api/v1/initiative", app.GenerateInitiative)
	srv.Register("GET /api/v1/initiative/encounter/{id}", app.GetInitiativeByEncounterId)
	// command api
	srv.Register("POST /api/v1/command", app.ExecuteCommand)
	srv.Register("POST /api/v1/info", app.AddSlackInfo)
	srv.Register("OPTIONS /*", srv.Options)

	logger.Info("starting chat web web server", "version", Version, "port", ExternalPort)
	// adding a token to admin user
	adminUser := "admin"
	adminToken := utils.RandomString(12)
	app.AddAdminSession(adminUser, adminToken)
	logger.Info("adding admin user one year token", "admin", adminUser, "token", adminToken)
	adminFile := utils.GetEnv("CREDS_FILE", "./creds")
	err := save(adminToken, adminFile)
	if err != nil {
		logger.Error("cannot write to file", "error", err)
		os.Exit(1)
	}

	// playerWeb := server.NewServer(ExternalPort, logger)
	// playerApi := handlers.NewPlayerApi(ctx, d, db, logger, playerWeb, rpgSystem)
	// playerWeb.Register("POST /api/v1/commands", playerApi.Commands)
	// playerWeb.Register("PUT /api/v1/reload", playerApi.ReloadCache)
	// starting a goroutine to server http requests
	// start web server in a goroutine
	go func() {
		if err := srv.Server.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			logger.Error("server listen and serve error", "error", err)
			os.Exit(1)
		}
		logger.Info("server stopped serving new connections.")
	}()
	// go func() {
	// 	if err := playerWeb.Server.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
	// 		logger.Error("chat web listen and serve error", "error", err)
	// 		os.Exit(1)
	// 	}
	// 	logger.Info("chat web stopped serving new connections.")
	// }()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	ctxTimeout, ctxCancel := context.WithTimeout(ctx, 10*time.Second)
	defer ctxCancel()

	if err := srv.Server.Shutdown(ctxTimeout); err != nil {
		logger.Error("server shutdown error", "error", err)
	}

	// if err := playerWeb.Server.Shutdown(ctxTimeout); err != nil {
	// 	logger.Error("chat web shutdown error", "error", err)
	// }
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

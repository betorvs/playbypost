package v1

import (
	"context"
	"log/slog"

	"github.com/betorvs/playbypost/app/server/handlers"
	"github.com/betorvs/playbypost/app/server/handlers/v1/validator"
	"github.com/betorvs/playbypost/app/server/worker"
	"github.com/betorvs/playbypost/core/rpg"
	"github.com/betorvs/playbypost/core/sys/db"
	"github.com/betorvs/playbypost/core/sys/library"
	"github.com/betorvs/playbypost/core/sys/web/cli"
	"github.com/betorvs/playbypost/core/sys/web/server"
)

type MainApi struct {
	Session   handlers.SessionChecker
	Worker    *worker.WorkerAPI
	Validator *validator.Validator
	logger    *slog.Logger
	lib       *library.Library
	s         *server.SvrWeb
	db        db.DBClient
	ctx       context.Context
	dice      rpg.Roll
	client    *cli.Cli
	rpg       *rpg.RPGSystem
	autoPlay  *rpg.RPGSystem
}

func NewMainApi(admin string, ctx context.Context, dice rpg.Roll, db db.DBClient, l *slog.Logger, s *server.SvrWeb, client *cli.Cli, rpgSystem *rpg.RPGSystem, lib *library.Library) *MainApi {
	session := handlers.NewSession(admin, l, db, s, ctx)
	autoPlay := rpg.LoadRPGSystemsDefault(rpg.AutoPlay)
	validator := validator.New(l, db, ctx)
	worker := worker.NewWorkerAPI(ctx, dice, db, l, client, rpgSystem, autoPlay, lib)
	return &MainApi{
		Session:   *session,
		Worker:    worker,
		Validator: validator,
		ctx:       ctx,
		dice:      dice,
		db:        db,
		logger:    l,
		s:         s,
		client:    client,
		lib:       lib,
		rpg:       rpgSystem,
		autoPlay:  autoPlay,
	}
}

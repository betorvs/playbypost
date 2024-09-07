package worker

import (
	"context"
	"log/slog"
	"sync"
	"time"

	"github.com/betorvs/playbypost/core/rpg"
	"github.com/betorvs/playbypost/core/sys/db"
	"github.com/betorvs/playbypost/core/sys/web/cli"
)

type Worker interface {
	Execute()
}

type WorkerAPI struct {
	StageActive    bool
	AutoPlayActive bool
	logger         *slog.Logger
	db             db.DBClient
	ctx            context.Context
	dice           rpg.Roll
	client         *cli.Cli
	rpg            *rpg.RPGSystem
	autoPlay       *rpg.RPGSystem
	AutoPlaySync   *sync.Mutex
	StageSync      *sync.Mutex
}

func NewWorkerAPI(ctx context.Context, dice rpg.Roll, db db.DBClient, l *slog.Logger, client *cli.Cli, rpgSystem *rpg.RPGSystem, auto *rpg.RPGSystem) *WorkerAPI {
	return &WorkerAPI{
		ctx:          ctx,
		dice:         dice,
		db:           db,
		logger:       l,
		client:       client,
		rpg:          rpgSystem,
		autoPlay:     auto,
		AutoPlaySync: &sync.Mutex{},
		StageSync:    &sync.Mutex{},
	}
}

func (a *WorkerAPI) Execute() {
	if a.StageActive {
		a.StageSync.Lock()
		a.logger.Info("starting scheduler worker api execution", "time", time.Now())
		activities, err := a.db.GetStageEncounterActivities(a.ctx)
		if err != nil {
			a.logger.Error("error getting stage encounter activities", "error", err.Error())
			return
		}
		for _, activity := range activities {
			if !activity.Processed {
				a.logger.Info("activity", "activity", activity)
				// execute activity
				err := a.parseCommand(activity)
				if err != nil {
					a.logger.Error("error parsing command", "error", err.Error())
					continue
				}

			}
		}
		a.StageSync.Unlock()
	}
	if a.AutoPlayActive {
		a.AutoPlaySync.Lock()
		a.logger.Info("starting scheduler auto play worker api execution", "time", time.Now())
		activities, err := a.db.GetAutoPlayActivities(a.ctx)
		if err != nil {
			a.logger.Error("error getting auto play activities", "error", err.Error())
			return
		}
		for _, activity := range activities {
			if !activity.Processed {
				a.logger.Info("activity", "activity", activity)
				// execute activity
				err := a.parseAutoPlayCommand(activity)
				if err != nil {
					a.logger.Error("error parsing auto play command", "error", err.Error())
					continue
				}

			}
		}
		a.AutoPlaySync.Unlock()
	}
}

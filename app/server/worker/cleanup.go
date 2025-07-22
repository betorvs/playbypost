package worker

import (
	"context"
	"log/slog"
	"time"

	"github.com/betorvs/playbypost/core/sys/db"
)

// Cleanup is a worker that cleans up expired sessions.

type Cleanup struct {
	logger *slog.Logger
	db     db.DBClient
	ctx    context.Context
}

// NewCleanup creates a new cleanup worker.
func NewCleanup(logger *slog.Logger, db db.DBClient, ctx context.Context) *Cleanup {
	return &Cleanup{
		logger: logger,
		db:     db,
		ctx:    ctx,
	}
}

// Start starts the cleanup worker.
func (c *Cleanup) Start() {
	ticker := time.NewTicker(1 * time.Minute)
	go func() {
		for range ticker.C {
			c.cleanup()
		}
	}()
}

func (c *Cleanup) cleanup() {
	c.logger.Info("cleaning up expired sessions")
	err := c.db.DeleteExpiredSessions(c.ctx)
	if err != nil {
		c.logger.Error("failed to delete expired sessions", "error", err)
	}
}

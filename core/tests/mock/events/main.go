package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/betorvs/playbypost/core/sys/web/types"
)

type app struct {
	logger *slog.Logger
}

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	slog.SetDefault(logger)
	mux := http.NewServeMux()
	server := &http.Server{
		Addr:    ":8091",
		Handler: mux,
	}
	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "{\"status\":\"OK\"}")
	})
	a := app{
		logger: logger,
	}
	mux.HandleFunc("POST /api/v1/event", a.events)
	mux.HandleFunc("GET /api/v1/validate", a.validate)
	if err := server.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
		logger.Error("listen and serve error", "error", err)
		os.Exit(1)
	}
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	ctxTimeout, ctxCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer ctxCancel()

	if err := server.Shutdown(ctxTimeout); err != nil {
		logger.Error("shutdown error", "error", err)
	}
	logger.Info("graceful shutdown complete.")
}

func (a *app) events(w http.ResponseWriter, r *http.Request) {
	// handle event
	obj := types.Event{}
	err := json.NewDecoder(r.Body).Decode(&obj)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "{\"msg\":\"json decode error\"}")
		return
	}
	a.logger.Info("event received", "event", obj)
	w.WriteHeader(http.StatusAccepted)
	fmt.Fprint(w, "{\"msg\":\"Accepted\"}")
}

func (a *app) validate(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "{\"msg\":\"authenticated\"}")
}

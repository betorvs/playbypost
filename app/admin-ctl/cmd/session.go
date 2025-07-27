/*
Copyright © 2024 Roberto Scudeller <beto.rvs@gmail.com>
*/
package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/spf13/cobra"
)

// sessionCmd represents the session command
var sessionCmd = &cobra.Command{
	Use:     "session [list|cleanup|events|stats]",
	Aliases: []string{"sessions"},
	Short:   "Manage sessions and view session analytics",
	Long:    `Commands for managing user sessions, viewing session events, and analyzing session patterns.`,
	Args:    cobra.ExactArgs(1),
	PreRun:  loadApp,
	Run: func(cmd *cobra.Command, args []string) {
		switch args[0] {
		case "list":
			listSessions()
		case "cleanup":
			cleanupSessions()
		case "events":
			listSessionEvents()
		case "stats":
			showSessionStats()
		case "delete":
			if sessionID == 0 {
				app.Logger.Error("session ID is required for delete command")
				os.Exit(1)
			}
			deleteSession(sessionID)
		default:
			app.Logger.Info("session command called")
		}
	},
}

func init() {
	rootCmd.AddCommand(sessionCmd)
	sessionCmd.Flags().Int64Var(&sessionID, "id", 0, "Session ID for specific operations")
}

func listSessions() {
	sessions, err := app.Web.GetAllSessions()
	if err != nil {
		app.Logger.Error("failed to get sessions", "error", err.Error())
		os.Exit(1)
	}

	for _, session := range sessions {
		switch outputFormat {
		case formatJSON:
			b, _ := json.Marshal(session)
			fmt.Println(string(b))
		case formatLog:
			app.Logger.Info("Session",
				"username", session.Username,
				"user_id", session.UserID,
				"client_type", session.ClientType,
				"ip_address", session.IPAddress,
				"created_at", session.CreatedAt,
				"last_activity", session.LastActivity,
				"expiry", session.Expiry,
				"expired", session.IsExpired(),
			)
		}
	}
}

func cleanupSessions() {
	app.Logger.Info("starting session cleanup")

	// Get count before cleanup
	sessions, err := app.Web.GetAllSessions()
	if err != nil {
		app.Logger.Error("failed to get sessions for cleanup", "error", err.Error())
		os.Exit(1)
	}

	expiredCount := 0
	for _, session := range sessions {
		if session.IsExpired() {
			expiredCount++
		}
	}

	app.Logger.Info("found expired sessions", "count", expiredCount)

	// Perform cleanup
	err = app.Web.DeleteExpiredSessions()
	if err != nil {
		app.Logger.Error("failed to cleanup sessions", "error", err.Error())
		os.Exit(1)
	}

	app.Logger.Info("session cleanup completed", "sessions_deleted", expiredCount)
}

func listSessionEvents() {
	events, err := app.Web.GetSessionEvents()
	if err != nil {
		app.Logger.Error("failed to get session events", "error", err.Error())
		os.Exit(1)
	}

	for _, event := range events {
		switch outputFormat {
		case formatJSON:
			b, _ := json.Marshal(event)
			fmt.Println(string(b))
		case formatLog:
			app.Logger.Info("Session Event",
				"id", event.ID,
				"session_id", event.SessionID,
				"event_type", event.EventType,
				"timestamp", event.Timestamp,
				"data", event.Data,
			)
		}
	}
}

func showSessionStats() {
	sessions, err := app.Web.GetAllSessions()
	if err != nil {
		app.Logger.Error("failed to get sessions for stats", "error", err.Error())
		os.Exit(1)
	}

	events, err := app.Web.GetSessionEvents()
	if err != nil {
		app.Logger.Error("failed to get session events for stats", "error", err.Error())
		os.Exit(1)
	}

	// Calculate statistics
	totalSessions := len(sessions)
	activeSessions := 0
	expiredSessions := 0
	clientTypes := make(map[string]int)

	for _, session := range sessions {
		if session.IsExpired() {
			expiredSessions++
		} else {
			activeSessions++
		}
		clientTypes[session.ClientType]++
	}

	// Event type statistics
	eventTypes := make(map[string]int)
	for _, event := range events {
		eventTypes[event.EventType]++
	}

	// Recent activity (last 24 hours)
	recentEvents := 0
	twentyFourHoursAgo := time.Now().Add(-24 * time.Hour)
	for _, event := range events {
		if event.Timestamp.After(twentyFourHoursAgo) {
			recentEvents++
		}
	}

	stats := map[string]interface{}{
		"total_sessions":    totalSessions,
		"active_sessions":   activeSessions,
		"expired_sessions":  expiredSessions,
		"client_types":      clientTypes,
		"total_events":      len(events),
		"recent_events_24h": recentEvents,
		"event_types":       eventTypes,
	}

	switch outputFormat {
	case formatJSON:
		b, _ := json.Marshal(stats)
		fmt.Println(string(b))
	case formatLog:
		app.Logger.Info("Session Statistics",
			"total_sessions", totalSessions,
			"active_sessions", activeSessions,
			"expired_sessions", expiredSessions,
			"total_events", len(events),
			"recent_events_24h", recentEvents,
		)
		app.Logger.Info("Client Type Distribution", "client_types", clientTypes)
		app.Logger.Info("Event Type Distribution", "event_types", eventTypes)
	}
}

func deleteSession(sessionID int64) {
	app.Logger.Info("deleting session", "session_id", sessionID)

	// Delete the session
	err := app.Web.DeleteSessionByID(sessionID)
	if err != nil {
		app.Logger.Error("failed to delete session", "error", err.Error())
		os.Exit(1)
	}

	app.Logger.Info("session deleted successfully", "session_id", sessionID)
}

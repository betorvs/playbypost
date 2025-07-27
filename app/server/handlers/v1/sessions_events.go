package v1

import (
	"net/http"
	"strconv"
)

func (a *MainApi) GetSessionEvents(w http.ResponseWriter, r *http.Request) {
	if a.Session.CheckAuth(r) {
		a.s.ErrJSON(w, http.StatusForbidden, "required authentication headers")
		return
	}
	events, err := a.db.GetSessionEvents(a.ctx)
	if err != nil {
		a.logger.Error("failed to get session events", "error", err)
		a.s.ErrJSON(w, http.StatusInternalServerError, "failed to get session events")
		return
	}
	a.s.JSON(w, events)
}

func (a *MainApi) GetActiveSessions(w http.ResponseWriter, r *http.Request) {
	if a.Session.CheckAuth(r) {
		a.s.ErrJSON(w, http.StatusForbidden, "required authentication headers")
		return
	}
	a.s.JSON(w, a.Session.GetActiveSessions())
}

// CleanupSessions manually triggers session cleanup
func (a *MainApi) CleanupSessions(w http.ResponseWriter, r *http.Request) {
	if a.Session.CheckAuth(r) {
		a.s.ErrJSON(w, http.StatusForbidden, "required authentication headers")
		return
	}

	err := a.db.DeleteExpiredSessions(a.ctx)
	if err != nil {
		a.logger.Error("failed to cleanup sessions", "error", err)
		a.s.ErrJSON(w, http.StatusInternalServerError, "failed to cleanup sessions")
		return
	}

	a.s.JSON(w, map[string]string{"message": "session cleanup completed"})
}

// DeleteSessionByID deletes a specific session by ID
func (a *MainApi) DeleteSessionByID(w http.ResponseWriter, r *http.Request) {
	if a.Session.CheckAuth(r) {
		a.s.ErrJSON(w, http.StatusForbidden, "required authentication headers")
		return
	}

	idString := r.PathValue("id")
	if idString == "" {
		a.s.ErrJSON(w, http.StatusBadRequest, "id cannot be empty")
		return
	}

	sessionID, err := strconv.ParseInt(idString, 10, 64)
	if err != nil {
		a.s.ErrJSON(w, http.StatusBadRequest, "invalid session id")
		return
	}

	err = a.Session.DeleteSessionByID(a.ctx, sessionID)
	if err != nil {
		a.logger.Error("failed to delete session by ID", "error", err)
		a.s.ErrJSON(w, http.StatusInternalServerError, "failed to delete session by ID")
		return
	}

	a.s.JSON(w, map[string]string{"message": "session deleted successfully"})
}

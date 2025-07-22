package v1

import "net/http"

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

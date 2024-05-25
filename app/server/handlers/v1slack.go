package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/betorvs/playbypost/core/sys/web/types"
	"github.com/lib/pq"
)

func (a MainApi) AddSlackInfo(w http.ResponseWriter, r *http.Request) {
	if a.checkAuth(r) {
		a.s.ErrJSON(w, http.StatusForbidden, "required authentication headers")
		return
	}
	obj := types.SlackInfo{}
	err := json.NewDecoder(r.Body).Decode(&obj)
	if err != nil {
		a.s.ErrJSON(w, http.StatusBadRequest, "json decode error")
		return
	}
	if obj.UserID == "" || obj.Username == "" || obj.Channel == "" {
		a.s.ErrJSON(w, http.StatusBadRequest, "empty body")
		return
	}
	_, err = a.db.AddSlackInformation(a.ctx, obj.UserID, obj.Username, obj.Channel)
	if err != nil {
		pgErr, ok := err.(*pq.Error)
		if ok {
			if pgErr.Code == "23505" {
				a.s.JSON(w, types.Msg{Msg: "already added"})
				return
			}
		}
		a.s.ErrJSON(w, http.StatusBadRequest, "error adding slack info to database")
		return
	}
	a.s.JSON(w, types.Msg{Msg: "added"})
}

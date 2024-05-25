package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/betorvs/playbypost/core/sys/web/types"
)

func (a MainApi) CreateTasks(w http.ResponseWriter, r *http.Request) {
	if a.checkAuth(r) {
		a.s.ErrJSON(w, http.StatusForbidden, "required authentication headers")
		return
	}
	obj := types.Task{}
	err := json.NewDecoder(r.Body).Decode(&obj)
	if err != nil {
		a.s.ErrJSON(w, http.StatusBadRequest, "json decode error")
		return
	}
	encounter, err := a.db.GetEncounterByID(a.ctx, obj.EncounterID)
	if err != nil {
		a.s.ErrJSON(w, http.StatusBadRequest, "encounter does not exist")
		return
	}

	res, err := a.db.CreateTask(a.ctx, obj.Title, obj.DisplayText, obj.Checks, int(obj.Kind), obj.Target, encounter.ID, obj.Options)
	if err != nil {
		a.s.ErrJSON(w, http.StatusBadRequest, "error creating task on database")
		return
	}
	msg := fmt.Sprintf("task id %v", res)
	a.s.JSON(w, types.Msg{Msg: msg})
}

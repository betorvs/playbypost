package v1

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/betorvs/playbypost/core/sys/web/types"
)

func (a MainApi) CreateTasks(w http.ResponseWriter, r *http.Request) {
	if a.Session.CheckAuth(r) {
		a.s.ErrJSON(w, http.StatusForbidden, "required authentication headers")
		return
	}
	obj := types.Task{}
	err := json.NewDecoder(r.Body).Decode(&obj)
	if err != nil {
		a.s.ErrJSON(w, http.StatusBadRequest, "json decode error")
		return
	}
	// encounter, err := a.db.GetEncounterByID(a.ctx, obj.EncounterID)
	// if err != nil {
	// 	a.s.ErrJSON(w, http.StatusBadRequest, "encounter does not exist")
	// 	return
	// }

	res, err := a.db.CreateTask(a.ctx, obj.Description, obj.Ability, obj.Skill, obj.Kind, obj.Target)
	if err != nil {
		a.logger.Error("create task issue", "error", err.Error())
		a.s.ErrJSON(w, http.StatusBadRequest, "error creating task on database")
		return
	}
	msg := fmt.Sprintf("task id %v", res)
	a.s.JSON(w, types.Msg{Msg: msg})
}

func (a MainApi) GetTask(w http.ResponseWriter, r *http.Request) {
	if a.Session.CheckAuth(r) {
		a.s.ErrJSON(w, http.StatusForbidden, "required authentication headers")
		return
	}
	obj, err := a.db.GetTask(a.ctx)
	if err != nil {
		a.s.ErrJSON(w, http.StatusBadRequest, "story database issue")
		return
	}
	a.s.JSON(w, obj)
}

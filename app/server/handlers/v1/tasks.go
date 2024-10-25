package v1

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

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

func (a MainApi) GetTaskByID(w http.ResponseWriter, r *http.Request) {
	if a.Session.CheckAuth(r) {
		a.s.ErrJSON(w, http.StatusForbidden, "required authentication headers")
		return
	}
	idString := r.PathValue("id")
	id, err := strconv.Atoi(idString)
	if err != nil {
		a.s.ErrJSON(w, http.StatusBadRequest, "id should be a integer")
		return
	}
	obj, err := a.db.GetTaskByID(a.ctx, id)
	if err != nil {
		a.s.ErrJSON(w, http.StatusBadRequest, "task does not exist")
		return
	}
	a.s.JSON(w, obj)
}

func (a MainApi) UpdateTaskByID(w http.ResponseWriter, r *http.Request) {
	if a.Session.CheckAuth(r) {
		a.s.ErrJSON(w, http.StatusForbidden, "required authentication headers")
		return
	}
	body := types.Task{}
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		a.s.ErrJSON(w, http.StatusBadRequest, "json decode error")
		return
	}
	idString := r.PathValue("id")
	id, err := strconv.Atoi(idString)
	if err != nil {
		a.s.ErrJSON(w, http.StatusBadRequest, "id should be a integer")
		return
	}
	if body.ID != 0 && body.ID != id {
		a.s.ErrJSON(w, http.StatusBadRequest, "id does not match with body")
		return
	}
	err = a.db.UpdateTaskByID(a.ctx, body.Description, body.Ability, body.Skill, body.Kind, body.Target, id)
	if err != nil {
		a.s.ErrJSON(w, http.StatusBadRequest, "error updating task on database")
		return
	}
	msg := fmt.Sprintf("task id %v updatred", id)
	a.s.JSON(w, types.Msg{Msg: msg})
}

func (a MainApi) DeleteTaskByID(w http.ResponseWriter, r *http.Request) {
	if a.Session.CheckAuth(r) {
		a.s.ErrJSON(w, http.StatusForbidden, "required authentication headers")
		return
	}
	idString := r.PathValue("id")
	id, err := strconv.Atoi(idString)
	if err != nil {
		a.s.ErrJSON(w, http.StatusBadRequest, "id should be a integer")
		return
	}
	err = a.db.DeleteTaskByID(a.ctx, id)
	if err != nil {
		a.s.ErrJSON(w, http.StatusBadRequest, "error deleting task on database")
		return
	}
	msg := fmt.Sprintf("task id %v deleted", id)
	a.s.JSON(w, types.Msg{Msg: msg})
}

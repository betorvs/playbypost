package v1

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/betorvs/playbypost/core/sys/web/types"
)

func (a MainApi) CreateAutoPlay(w http.ResponseWriter, r *http.Request) {
	if a.Session.CheckAuth(r) {
		a.s.ErrJSON(w, http.StatusForbidden, "required authentication headers")
		return
	}
	obj := types.AutoPlayStart{}
	err := json.NewDecoder(r.Body).Decode(&obj)
	if err != nil {
		a.logger.Error("json error ", "error", err.Error())
		a.s.ErrJSON(w, http.StatusBadRequest, "invalid json")
		return
	}
	if obj.StoryID == 0 || obj.Text == "" {
		a.s.ErrJSON(w, http.StatusBadRequest, "story_id and text cannot be empty")
		return
	}
	// create auto play
	res, err := a.db.CreateAutoPlayTx(a.ctx, obj.Text, obj.StoryID, obj.Solo)
	if err != nil {
		a.logger.Error("error creating auto play", "error", err.Error())
		a.s.ErrJSON(w, http.StatusInternalServerError, "error creating auto play")
		return
	}
	msg := fmt.Sprintf("auto_play_id %v", res)
	a.s.JSON(w, types.Msg{Msg: msg})
}

func (a MainApi) GetAutoPlay(w http.ResponseWriter, r *http.Request) {
	if a.Session.CheckAuth(r) {
		a.s.ErrJSON(w, http.StatusForbidden, "required authentication headers")
		return
	}
	obj, err := a.db.GetAutoPlay(a.ctx)
	if err != nil {
		a.s.ErrJSON(w, http.StatusBadRequest, "auto play database issue")
		return
	}
	a.s.JSON(w, obj)
}

func (a MainApi) GetAutoPlayByID(w http.ResponseWriter, r *http.Request) {
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
	obj, err := a.db.GetAutoPlayByID(a.ctx, id)
	if err != nil {
		a.s.ErrJSON(w, http.StatusBadRequest, "auto play database issue")
		return
	}
	a.s.JSON(w, obj)
}

func (a MainApi) GetNextEncounterByStoryId(w http.ResponseWriter, r *http.Request) {
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
	obj, err := a.db.GetNextEncounterByStoryID(a.ctx, id)
	if err != nil {
		a.s.ErrJSON(w, http.StatusBadRequest, "auto play database issue")
		return
	}
	a.s.JSON(w, obj)
}

func (a MainApi) AddAutoPlayNext(w http.ResponseWriter, r *http.Request) {
	if a.Session.CheckAuth(r) {
		a.s.ErrJSON(w, http.StatusForbidden, "required authentication headers")
		return
	}
	obj := types.Next{}
	err := json.NewDecoder(r.Body).Decode(&obj)
	if err != nil {
		a.s.ErrJSON(w, http.StatusBadRequest, "json decode error")
		return
	}
	if obj.NextEncounterID == 0 || obj.EncounterID == 0 || obj.UpstreamID == 0 {
		a.s.ErrJSON(w, http.StatusBadRequest, "next encounter id, encounter id and auto play id cannot be empty")
		return
	}
	if obj.Objective.Kind == "" {
		obj.Objective.Kind = types.ObjectiveDefault
		obj.Objective.Values = []int{0}
	}
	a.logger.Info("add auto play next", "obj", obj)
	err = a.db.AddAutoPlayNext(a.ctx, obj)
	if err != nil {
		a.s.ErrJSON(w, http.StatusBadRequest, "error adding next encounter to encounter on database")
		return
	}
	msg := fmt.Sprintf("encounter id %v next encounter updated", obj.EncounterID)
	a.s.JSON(w, types.Msg{Msg: msg})
}

func (a MainApi) GetAutoPlayNextEncounterByAutoPlayID(w http.ResponseWriter, r *http.Request) {
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
	obj, err := a.db.GetNextEncounterByAutoPlayID(a.ctx, id)
	if err != nil {
		a.s.ErrJSON(w, http.StatusBadRequest, "auto play database issue")
		return
	}
	a.s.JSON(w, obj)
}

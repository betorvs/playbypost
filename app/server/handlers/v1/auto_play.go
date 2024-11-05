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
	res, err := a.db.CreateAutoPlayTx(a.ctx, obj.Text, obj.StoryID, obj.CreatorID, obj.Solo)
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

func (a MainApi) GetAutoPlayEncounterListByStoryID(w http.ResponseWriter, r *http.Request) {
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
	obj, err := a.db.GetAutoPlayEncounterListByStoryID(a.ctx, id)
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
	obj := []types.Next{}
	err := json.NewDecoder(r.Body).Decode(&obj)
	if err != nil {
		a.s.ErrJSON(w, http.StatusBadRequest, "json decode error")
		return
	}
	valid, err := types.ValidateNextSlice(obj, types.UpstreamKindAutoPlay)
	if err != nil {
		a.logger.Error("error validating next encounter", "error", err.Error())
		a.s.ErrJSON(w, http.StatusBadRequest, err.Error())
		return
	}
	a.logger.Info("add auto play next", "obj", valid)
	err = a.db.AddAutoPlayNext(a.ctx, valid)
	if err != nil {
		a.logger.Error("error adding next encounter to encounter on database", "error", err.Error())
		a.s.ErrJSON(w, http.StatusBadRequest, "error adding next encounter to encounter on database")
		return
	}
	msg := fmt.Sprintf("encounter id %v next encounter updated", valid[0].EncounterID)
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

func (a MainApi) ChangePublishFlagAutoPlay(w http.ResponseWriter, r *http.Request) {
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
	// check auto play id
	autoPlay, err := a.db.GetAutoPlayByID(a.ctx, id)
	if err != nil {
		a.s.ErrJSON(w, http.StatusBadRequest, "auto play database issue")
		return
	}
	// auto play id not found
	if autoPlay.ID == 0 {
		a.s.ErrJSON(w, http.StatusBadRequest, "auto play id not found")
		return
	}
	// update auto play to publish
	// use the opposite value from DB
	// default value in DB: false
	value := !autoPlay.Publish
	err = a.db.ChangePublishAutoPlay(a.ctx, id, value)
	if err != nil {
		a.s.ErrJSON(w, http.StatusBadRequest, "auto play database issue")
		return
	}
	msg := fmt.Sprintf("auto play id %v publish %v", id, value)
	a.s.JSON(w, types.Msg{Msg: msg})
}

func (a MainApi) DeleteAutoPlayNextEncounter(w http.ResponseWriter, r *http.Request) {
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
	err = a.db.DeleteAutoPlayNextEncounter(a.ctx, id)
	if err != nil {
		a.s.ErrJSON(w, http.StatusBadRequest, "auto play database issue")
		return
	}
	msg := fmt.Sprintf("next encounter id %v deleted", id)
	a.s.JSON(w, types.Msg{Msg: msg})
}

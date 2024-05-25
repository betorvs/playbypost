package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/betorvs/playbypost/core/sys/web/types"
	"github.com/betorvs/playbypost/core/utils"
)

func (a MainApi) GetEncounters(w http.ResponseWriter, r *http.Request) {
	if a.checkAuth(r) {
		a.s.ErrJSON(w, http.StatusForbidden, "required authentication headers")
		return
	}
	obj, err := a.db.GetEncounters(a.ctx)
	if err != nil {
		a.s.ErrJSON(w, http.StatusBadRequest, "encounters issue")
		return
	}
	a.s.JSON(w, obj)
}

func (a MainApi) GetEncounterById(w http.ResponseWriter, r *http.Request) {
	if a.checkAuth(r) {
		a.s.ErrJSON(w, http.StatusForbidden, "required authentication headers")
		return
	}
	idString := r.PathValue("id")
	id, err := strconv.Atoi(idString)
	if err != nil {
		a.s.ErrJSON(w, http.StatusBadRequest, "id should be a integer")
		return
	}
	obj, err := a.db.GetEncounterByID(a.ctx, id)
	if err != nil {
		a.s.ErrJSON(w, http.StatusBadRequest, "encounters issue")
		return
	}
	a.s.JSON(w, obj)
}

func (a MainApi) GetEncounterByStoryId(w http.ResponseWriter, r *http.Request) {
	if a.checkAuth(r) {
		a.s.ErrJSON(w, http.StatusForbidden, "required authentication headers")
		return
	}
	headerUsername := r.Header.Get(types.HeaderUsername)
	idString := r.PathValue("id")
	id, err := strconv.Atoi(idString)
	if err != nil {
		a.s.ErrJSON(w, http.StatusBadRequest, "id should be a integer")
		return
	}
	obj, err := a.db.GetEncounterByStoryID(a.ctx, id)
	if err != nil {
		a.s.ErrJSON(w, http.StatusBadRequest, "encounters issue")
		return
	}
	masterID := -1
	if len(obj) > 0 {
		masterID = obj[0].StorytellerID
	}
	if a.Sessions.Current[headerUsername].UserID != masterID {
		a.s.JSON(w, obj)
		return
	}
	encounters := []types.Encounters{}
	for _, v := range obj {
		announce, _ := utils.DecryptText(v.Announcement, a.Sessions.Current[headerUsername].EncodingKey)
		note, _ := utils.DecryptText(v.Notes, a.Sessions.Current[headerUsername].EncodingKey)
		encounters = append(encounters, types.Encounters{
			ID:           v.ID,
			Title:        v.Title,
			Announcement: announce,
			Notes:        note,
			StoryID:      v.StoryID,
		})
	}

	a.s.JSON(w, encounters)
}

func (a MainApi) CreateEncounter(w http.ResponseWriter, r *http.Request) {
	obj := types.Encounters{}
	err := json.NewDecoder(r.Body).Decode(&obj)
	if err != nil {
		a.s.ErrJSON(w, http.StatusBadRequest, "json decode error")
		return
	}
	if obj.StoryID == 0 || obj.Title == "" {
		a.s.ErrJSON(w, http.StatusBadRequest, "title and story_id cannot be empty")
		return
	}
	story, err := a.db.GetStoryByID(a.ctx, obj.StoryID)
	if err != nil {
		a.s.ErrJSON(w, http.StatusBadRequest, "story does not exist")
		return
	}
	user, err := a.db.GetStorytellerByID(a.ctx, story.StorytellerID)
	if err != nil {
		a.s.ErrJSON(w, http.StatusBadRequest, "master id not found")
		return
	}
	announce, err := utils.EncryptText(obj.Announcement, user.EncodingKeys[story.ID])
	if err != nil {
		a.s.ErrJSON(w, http.StatusBadRequest, "announcement encoding fails")
		return
	}
	notes, err := utils.EncryptText(obj.Notes, user.EncodingKeys[story.ID])
	if err != nil {
		a.s.ErrJSON(w, http.StatusBadRequest, "notes encoding fails")
		return
	}
	res, err := a.db.CreateEncounter(a.ctx, obj.Title, announce, notes, obj.StoryID, obj.StorytellerID)
	if err != nil {
		a.s.ErrJSON(w, http.StatusBadRequest, "error creating encounter on database")
		return
	}
	msg := fmt.Sprintf("encounter id %v", res)
	a.s.JSON(w, types.Msg{Msg: msg})
}

// func (a MainApi) UpdateEncounterPhaseById(w http.ResponseWriter, r *http.Request) {
// 	idString := r.PathValue("id")
// 	id, err := strconv.Atoi(idString)
// 	if err != nil {
// 		a.s.ErrJSON(w, http.StatusBadRequest, "id should be a integer")
// 		return
// 	}
// 	phaseString := r.PathValue("phase")
// 	phase, err := strconv.Atoi(phaseString)
// 	if err != nil {
// 		a.s.ErrJSON(w, http.StatusBadRequest, "phase should be a integer")
// 		return
// 	}
// 	err = a.db.UpdatePhase(a.ctx, id, phase)
// 	if err != nil {
// 		a.s.ErrJSON(w, http.StatusBadRequest, "encounters issue")
// 		return
// 	}
// 	status := types.PhaseAtoi(phase)
// 	a.logger.Info("change phase worked", "phase", status)
// 	// starts a initiative
// 	a.s.JSON(w, types.Msg{Msg: fmt.Sprintf("change to phase: %s", status)})
// }

func (a MainApi) AddParticipants(w http.ResponseWriter, r *http.Request) {
	obj := types.Participants{}
	err := json.NewDecoder(r.Body).Decode(&obj)
	if err != nil {
		a.s.ErrJSON(w, http.StatusBadRequest, "json decode error")
		return
	}
	if len(obj.PlayersID) == 0 || obj.EncounterID == 0 {
		a.s.ErrJSON(w, http.StatusBadRequest, "players id list and encounter id cannot be empty")
		return
	}
	err = a.db.AddParticipants(a.ctx, obj.EncounterID, obj.NPC, obj.PlayersID)
	if err != nil {
		a.s.ErrJSON(w, http.StatusBadRequest, "error adding participants to encounter on database")
		return
	}
	msg := fmt.Sprintf("encounter id %v participants updated", obj.EncounterID)
	a.s.JSON(w, types.Msg{Msg: msg})
}

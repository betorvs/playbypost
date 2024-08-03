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
	writerID := -1
	if len(obj) > 0 {
		writerID = obj[0].WriterID
	}
	// if a.Sessions.Current[headerUsername].UserID != masterID {
	// 	a.s.JSON(w, obj)
	// 	return
	// }
	user, err := a.db.GetWriterByID(a.ctx, writerID)
	if err != nil {
		a.s.ErrJSON(w, http.StatusBadRequest, "writer id not found")
		return
	}
	if user.Username != headerUsername {
		a.logger.Info("username does not match with header", "username", user.Username, "header", headerUsername)
		a.s.JSON(w, obj)
		return
	}
	encounters := []types.Encounter{}
	for _, v := range obj {
		announce, _ := utils.DecryptText(v.Announcement, user.EncodingKeys[id])
		note, _ := utils.DecryptText(v.Notes, user.EncodingKeys[id])
		encounters = append(encounters, types.Encounter{
			ID:           v.ID,
			Title:        v.Title,
			Announcement: announce,
			Notes:        note,
			StoryID:      v.StoryID,
			WriterID:     v.WriterID,
		})
	}

	a.s.JSON(w, encounters)
}

func (a MainApi) CreateEncounter(w http.ResponseWriter, r *http.Request) {
	obj := types.Encounter{}
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
	user, err := a.db.GetWriterByID(a.ctx, story.WriterID)
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
	res, err := a.db.CreateEncounter(a.ctx, obj.Title, announce, notes, obj.StoryID, obj.WriterID)
	if err != nil {
		a.s.ErrJSON(w, http.StatusBadRequest, "error creating encounter on database")
		return
	}
	msg := fmt.Sprintf("encounter id %v", res)
	a.s.JSON(w, types.Msg{Msg: msg})
}

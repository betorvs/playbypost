package v1

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/betorvs/playbypost/core/sys/web/types"
	"github.com/betorvs/playbypost/core/utils"
)

func (a MainApi) GetEncounters(w http.ResponseWriter, r *http.Request) {
	if a.Session.CheckAuth(r) {
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
	obj, err := a.db.GetEncounterByID(a.ctx, id)
	if err != nil {
		a.s.ErrJSON(w, http.StatusBadRequest, "encounters issue")
		return
	}
	headerUsername := r.Header.Get(types.HeaderUsername)
	user, err := a.db.GetWriterByID(a.ctx, obj.WriterID)
	if err != nil {
		a.s.ErrJSON(w, http.StatusBadRequest, "writer id not found")
		return
	}
	if user.Username != headerUsername {
		a.logger.Debug("username does not match with header", "username", user.Username, "header", headerUsername)
		a.s.JSON(w, obj)
		return
	}
	announce, err := utils.DecryptText(obj.Announcement, user.EncodingKeys[obj.StoryID])
	if err != nil {
		a.s.ErrJSON(w, http.StatusBadRequest, "announcement decoding fails")
		return
	}
	note, err := utils.DecryptText(obj.Notes, user.EncodingKeys[obj.StoryID])
	if err != nil {
		a.s.ErrJSON(w, http.StatusBadRequest, "notes decoding fails")
		return
	}
	obj.Announcement = announce
	obj.Notes = note
	a.s.JSON(w, obj)
}

func (a MainApi) GetEncounterByStoryId(w http.ResponseWriter, r *http.Request) {
	if a.Session.CheckAuth(r) {
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
	writerID := -1
	var encounters []types.Encounter
	// get params
	limit := r.URL.Query().Get("limit")
	cursor := r.URL.Query().Get("cursor")

	if limit != "" {
		limitInt, err := strconv.Atoi(limit)
		if err != nil {
			a.s.ErrJSON(w, http.StatusBadRequest, "limit should be a integer")
			return
		}
		if limitInt > 10 {
			limitInt = 10
		}
		if limitInt < 3 {
			limitInt = 3
		}
		lastID := 0
		if cursor != "" {
			lastIDIntTmp, err := strconv.Atoi(cursor)
			if err != nil {
				a.s.ErrJSON(w, http.StatusBadRequest, "cursor should be a integer")
				return
			}
			a.logger.Debug("cursor", "cursor", lastIDIntTmp)
			lastID = lastIDIntTmp
		}
		obj, cursor, total, err := a.db.GetEncounterByStoryIDWithPagination(a.ctx, id, limitInt, lastID)
		if err != nil {
			a.s.ErrJSON(w, http.StatusBadRequest, "get encounters with pagination issue")
			return
		}
		if len(obj) > 0 {
			writerID = obj[0].WriterID
		}
		if cursor > 0 {
			// get URI from request
			uri := r.RequestURI
			uri = uri + "&cursor=" + strconv.Itoa(cursor)
			w.Header().Set("X-Cursor-URI", uri)
			w.Header().Set("X-Last-Id", strconv.Itoa(cursor))
			w.Header().Set("X-Total-Count", strconv.Itoa(total))
		}
		encounters = obj
	} else {
		obj, err := a.db.GetEncounterByStoryID(a.ctx, id)
		if err != nil {
			a.s.ErrJSON(w, http.StatusBadRequest, "encounters issue")
			return
		}

		if len(obj) > 0 {
			writerID = obj[0].WriterID
		}
		encounters = obj
	}

	user, err := a.db.GetWriterByID(a.ctx, writerID)
	if err != nil {
		a.s.ErrJSON(w, http.StatusBadRequest, "writer id not found")
		return
	}
	if user.Username != headerUsername {
		a.logger.Debug("username does not match with header", "username", user.Username, "header", headerUsername)
		a.s.JSON(w, encounters)
		return
	}
	encountersDecrypted := []types.Encounter{}
	for _, v := range encounters {
		announce, _ := utils.DecryptText(v.Announcement, user.EncodingKeys[id])
		note, _ := utils.DecryptText(v.Notes, user.EncodingKeys[id])
		encountersDecrypted = append(encountersDecrypted, types.Encounter{
			ID:             v.ID,
			Title:          v.Title,
			Announcement:   announce,
			Notes:          note,
			StoryID:        v.StoryID,
			WriterID:       v.WriterID,
			FirstEncounter: v.FirstEncounter,
			LastEncounter:  v.LastEncounter,
		})
	}

	a.s.JSON(w, encountersDecrypted)
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
	res, err := a.db.CreateEncounterTx(a.ctx, obj.Title, announce, notes, obj.StoryID, obj.WriterID, obj.FirstEncounter, obj.LastEncounter)
	if err != nil {
		m := fmt.Sprintf("error creating encounter on database %v", err.Error())
		a.s.ErrJSON(w, http.StatusBadRequest, m)
		return
	}
	msg := fmt.Sprintf("encounter id %v", res)
	a.s.JSON(w, types.Msg{Msg: msg})
}

func (a MainApi) UpdateEncounterByID(w http.ResponseWriter, r *http.Request) {
	if a.Session.CheckAuth(r) {
		a.s.ErrJSON(w, http.StatusForbidden, "required authentication headers")
		return
	}
	body := types.Encounter{}
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		a.s.ErrJSON(w, http.StatusBadRequest, "json decode error")
		return
	}
	if body.StoryID == 0 || body.Title == "" || body.Announcement == "" || body.Notes == "" {
		a.s.ErrJSON(w, http.StatusBadRequest, "title, story_id, announcement and notes cannot be empty")
		return
	}
	headerUsername := r.Header.Get(types.HeaderUsername)
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
	a.logger.Debug("encounter", "id", obj.ID, "text", obj.Title, "story_id", obj.StoryID, "writer_id", obj.WriterID)
	if body.ID != 0 && body.ID != id {
		a.s.ErrJSON(w, http.StatusBadRequest, "id does not match with body")
		return
	}
	if body.WriterID != obj.WriterID {
		a.s.ErrJSON(w, http.StatusBadRequest, "writer id does not match with body")
		return
	}
	user, err := a.db.GetWriterByID(a.ctx, obj.WriterID)
	if err != nil {
		a.s.ErrJSON(w, http.StatusBadRequest, "master id not found")
		return
	}
	if headerUsername != user.Username {
		a.logger.Debug("username does not match with header", "username", user.Username, "header", headerUsername)
		if headerUsername != a.Session.Admin() {
			a.s.ErrJSON(w, http.StatusForbidden, "username does not match with header")
			return
		}
	}
	announce, err := utils.EncryptText(body.Announcement, user.EncodingKeys[body.StoryID])
	if err != nil {
		a.logger.Error("announcement encoding fails", "error", err.Error())
		a.s.ErrJSON(w, http.StatusBadRequest, "announcement encoding fails")
		return
	}
	notes, err := utils.EncryptText(body.Notes, user.EncodingKeys[body.StoryID])
	if err != nil {
		a.logger.Error("notes encoding fails", "error", err.Error())
		a.s.ErrJSON(w, http.StatusBadRequest, "notes encoding fails")
		return
	}
	res, err := a.db.UpdateEncounterTx(a.ctx, body.Title, announce, notes, obj.ID, body.StoryID, body.FirstEncounter, body.LastEncounter)
	if err != nil {
		a.s.ErrJSON(w, http.StatusBadRequest, "error updating encounter on database")
		return
	}

	msg := fmt.Sprintf("encounter id %v updated", res)
	a.s.JSON(w, types.Msg{Msg: msg})
}

func (a MainApi) DeleteEncounterByID(w http.ResponseWriter, r *http.Request) {
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
	err = a.db.DeleteEncounterByID(a.ctx, id)
	if err != nil {
		a.s.ErrJSON(w, http.StatusBadRequest, "error deleting encounter on database")
		return
	}
	msg := fmt.Sprintf("encounter id %v deleted", id)
	a.s.JSON(w, types.Msg{Msg: msg})
}

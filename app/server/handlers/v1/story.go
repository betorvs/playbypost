package v1

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/betorvs/playbypost/core/sys/web/types"
	"github.com/betorvs/playbypost/core/utils"
)

func (a MainApi) GetStory(w http.ResponseWriter, r *http.Request) {
	if a.Session.CheckAuth(r) {
		a.s.ErrJSON(w, http.StatusForbidden, "required authentication headers")
		return
	}
	obj, err := a.db.GetStory(a.ctx)
	if err != nil {
		a.s.ErrJSON(w, http.StatusBadRequest, "story database issue")
		return
	}
	a.s.JSON(w, obj)
}

func (a MainApi) CreateStory(w http.ResponseWriter, r *http.Request) {
	if a.Session.CheckAuth(r) {
		a.s.ErrJSON(w, http.StatusForbidden, "required authentication headers")
		return
	}
	obj := types.Story{}
	err := json.NewDecoder(r.Body).Decode(&obj)
	if err != nil {
		a.logger.Error("json error ", "error", err.Error())
		a.s.ErrJSON(w, http.StatusBadRequest, "invalid json")
		return
	}
	if obj.WriterID == 0 || obj.Title == "" {
		a.s.ErrJSON(w, http.StatusBadRequest, "title and master_id cannot be empty")
		return
	}
	// user, err := a.db.GetWriterByID(a.ctx, obj.WriterID)
	// if err != nil {
	// 	a.s.ErrJSON(w, http.StatusBadRequest, "Writer id not found")
	// 	return
	// }
	newEncodingKey := utils.RandomString(16)
	announce, err := utils.EncryptText(obj.Announcement, newEncodingKey)
	if err != nil {
		a.s.ErrJSON(w, http.StatusInternalServerError, "announcement encoding fails")
		return
	}
	notes, err := utils.EncryptText(obj.Notes, newEncodingKey)
	if err != nil {
		a.s.ErrJSON(w, http.StatusInternalServerError, "notes encoding fails")
		return
	}
	res, err := a.db.CreateStoryTx(a.ctx, obj.Title, announce, notes, newEncodingKey, obj.WriterID)
	if err != nil {
		m := fmt.Sprintf("error creating story on database\n%v", err)
		a.s.ErrJSON(w, http.StatusBadGateway, m)
		return
	}
	msg := fmt.Sprintf("story id %v", res)
	a.s.JSON(w, types.Msg{Msg: msg})
}

func (a MainApi) UpdateStoryByID(w http.ResponseWriter, r *http.Request) {
	if a.Session.CheckAuth(r) {
		a.s.ErrJSON(w, http.StatusForbidden, "required authentication headers")
		return
	}
	body := types.Story{}
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		a.logger.Error("json error ", "error", err.Error())
		a.s.ErrJSON(w, http.StatusBadRequest, "invalid json")
		return
	}
	if body.Title == "" && body.Announcement == "" && body.Notes == "" {
		a.s.ErrJSON(w, http.StatusBadRequest, "title, announcement and notes cannot be empty")
		return
	}
	headerUsername := r.Header.Get(types.HeaderUsername)
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
	obj, err := a.db.GetStoryByID(a.ctx, id)
	if err != nil {
		a.s.ErrJSON(w, http.StatusBadRequest, "story issue")
		return
	}
	if body.WriterID != obj.WriterID {
		a.s.ErrJSON(w, http.StatusBadRequest, "writer id does not match with body")
		return
	}
	user, err := a.db.GetWriterByID(a.ctx, obj.WriterID)
	if err != nil {
		a.s.ErrJSON(w, http.StatusBadRequest, "writer id not found")
		return
	}
	if headerUsername != user.Username {
		a.logger.Debug("username does not match with header", "username", user.Username, "header", headerUsername)
		if headerUsername != a.Session.Admin() {
			a.s.ErrJSON(w, http.StatusForbidden, "username does not match with header")
			return
		}

	}
	announce, err := utils.EncryptText(body.Announcement, user.EncodingKeys[obj.ID])
	if err != nil {
		a.logger.Error("error announce", "writer_id", obj.WriterID)
		a.s.ErrJSON(w, http.StatusBadRequest, "announcement encoding fails")
		return
	}
	note, err := utils.EncryptText(body.Notes, user.EncodingKeys[obj.ID])
	if err != nil {
		a.logger.Error("error note ", "writer_id", obj.WriterID)
		a.s.ErrJSON(w, http.StatusBadRequest, "notes encoding fails")
		return
	}
	res, err := a.db.UpdateStoryTx(a.ctx, body.Title, announce, note, obj.ID)
	if err != nil {
		a.s.ErrJSON(w, http.StatusBadRequest, "story issue")
		return
	}

	msg := fmt.Sprintf("story id %v updatred", res)
	a.s.JSON(w, types.Msg{Msg: msg})
}

func (a MainApi) GetStoryById(w http.ResponseWriter, r *http.Request) {
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
	obj, err := a.db.GetStoryByID(a.ctx, id)
	if err != nil {
		a.s.ErrJSON(w, http.StatusBadRequest, "story issue")
		return
	}
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
	a.logger.Debug("obj from db", "obj", obj)
	a.logger.Debug("used from db", "user", user)
	announce, _ := utils.DecryptText(obj.Announcement, user.EncodingKeys[obj.ID])
	note, _ := utils.DecryptText(obj.Notes, user.EncodingKeys[obj.ID])
	story := types.Story{
		ID:           obj.ID,
		Title:        obj.Title,
		Announcement: announce,
		Notes:        note,
		WriterID:     obj.WriterID,
	}
	a.s.JSON(w, story)
}

func (a MainApi) GetStoryByWriterId(w http.ResponseWriter, r *http.Request) {
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
	obj, err := a.db.GetStoriesByWriterID(a.ctx, id)
	if err != nil {
		a.s.ErrJSON(w, http.StatusBadRequest, "story issue")
		return
	}
	user, err := a.db.GetWriterByID(a.ctx, id)
	if err != nil {
		a.s.ErrJSON(w, http.StatusBadRequest, "writer id not found")
		return
	}
	if user.Username != headerUsername {
		a.logger.Debug("username does not match with header")
		a.s.JSON(w, obj)
		return
	}
	stories := []types.Story{}
	for _, v := range obj {
		announce, _ := utils.DecryptText(v.Announcement, user.EncodingKeys[v.ID])
		note, _ := utils.DecryptText(v.Notes, user.EncodingKeys[v.ID])
		stories = append(stories, types.Story{
			ID:           v.ID,
			Title:        v.Title,
			Announcement: announce,
			Notes:        note,
			WriterID:     v.WriterID,
		})
	}

	a.s.JSON(w, stories)
}

func (a MainApi) DeleteStoryByID(w http.ResponseWriter, r *http.Request) {
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
	err = a.db.DeleteStoryByID(a.ctx, id)
	if err != nil {
		a.s.ErrJSON(w, http.StatusBadRequest, "story issue")
		return
	}
	msg := fmt.Sprintf("story id %v deleted", id)
	a.s.JSON(w, types.Msg{Msg: msg})
}

package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/betorvs/playbypost/core/sys/web/types"
	"github.com/betorvs/playbypost/core/utils"
)

func (a MainApi) GetStory(w http.ResponseWriter, r *http.Request) {
	if a.checkAuth(r) {
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
	if a.checkAuth(r) {
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
	if obj.StorytellerID == 0 || obj.Title == "" {
		a.s.ErrJSON(w, http.StatusBadRequest, "title and master_id cannot be empty")
		return
	}
	// user, err := a.db.GetStorytellerByID(a.ctx, obj.StorytellerID)
	// if err != nil {
	// 	a.s.ErrJSON(w, http.StatusBadRequest, "Storyteller id not found")
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
	res, err := a.db.CreateStoryTx(a.ctx, obj.Title, announce, notes, newEncodingKey, obj.StorytellerID)
	if err != nil {
		a.logger.Error("error ", "storyteller_id", obj.StorytellerID)
		a.s.ErrJSON(w, http.StatusBadGateway, "error creating story on database")
		return
	}
	msg := fmt.Sprintf("story id %v", res)
	a.s.JSON(w, types.Msg{Msg: msg})
}

func (a MainApi) GetStoryById(w http.ResponseWriter, r *http.Request) {
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
	obj, err := a.db.GetStoryByID(a.ctx, id)
	if err != nil {
		a.s.ErrJSON(w, http.StatusBadRequest, "story issue")
		return
	}
	if a.Sessions.Current[headerUsername].UserID != obj.StorytellerID {
		a.s.JSON(w, obj)
		return
	}
	announce, _ := utils.DecryptText(obj.Announcement, a.Sessions.Current[headerUsername].EncodingKey)
	note, _ := utils.DecryptText(obj.Notes, a.Sessions.Current[headerUsername].EncodingKey)
	story := types.Story{
		ID:            obj.ID,
		Title:         obj.Title,
		Announcement:  announce,
		Notes:         note,
		StorytellerID: obj.StorytellerID,
	}
	a.s.JSON(w, story)
}

func (a MainApi) GetStoryByMasterId(w http.ResponseWriter, r *http.Request) {
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
	obj, err := a.db.GetStoriesByStorytellerID(a.ctx, id)
	if err != nil {
		a.s.ErrJSON(w, http.StatusBadRequest, "story issue")
		return
	}
	user, err := a.db.GetStorytellerByID(a.ctx, id)
	if err != nil {
		a.s.ErrJSON(w, http.StatusBadRequest, "Storyteller id not found")
		return
	}
	if user.Username != headerUsername {
		a.logger.Info("username does not match with header")
		a.s.JSON(w, obj)
		return
	}
	stories := []types.Story{}
	for _, v := range obj {
		announce, _ := utils.DecryptText(v.Announcement, user.EncodingKeys[v.ID])
		note, _ := utils.DecryptText(v.Notes, user.EncodingKeys[v.ID])
		stories = append(stories, types.Story{
			ID:            v.ID,
			Title:         v.Title,
			Announcement:  announce,
			Notes:         note,
			StorytellerID: v.StorytellerID,
		})
	}

	a.s.JSON(w, stories)
}

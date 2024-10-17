package v1

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/betorvs/playbypost/app/server/handlers/v1/validator"
	"github.com/betorvs/playbypost/core/sys/web/types"
	"github.com/betorvs/playbypost/core/utils"
)

func (a MainApi) GetAllValidations(w http.ResponseWriter, r *http.Request) {
	if a.Session.CheckAuth(r) {
		a.s.ErrJSON(w, http.StatusForbidden, "required authentication headers")
		return
	}
	// a.s.JSON(w, validatorRequestToSlice(a.Validator.Request))
	a.s.JSON(w, a.Validator.Slice())
}

func (a MainApi) GetValidateAutoPlay(w http.ResponseWriter, r *http.Request) {
	if a.Session.CheckAuth(r) {
		a.s.ErrJSON(w, http.StatusForbidden, "required authentication headers")
		return
	}
	output := r.URL.Query().Get("output")
	hashID := r.PathValue("hashid")
	if output == "json" {
		id, err := strconv.Atoi(hashID)
		if err != nil {
			a.s.ErrJSON(w, http.StatusBadRequest, "id should be a integer")
			return
		}
		a.logger.Info("output json set and looking for id", "id", id)
		slice := a.Validator.Slice()
		for _, v := range slice {
			if v.ID == id {
				a.s.JSON(w, v)
				return
			}
		}
	}
	msg := types.Composed{}
	v, ok := a.Validator.Request[hashID]
	if ok {
		msg.Msg = fmt.Sprintf("%s with result: %v", v.Output, v.Valid)
		if len(v.Analise.Results) > 0 {
			details := make(map[string]string)
			for i, r := range v.Analise.Results {
				details[fmt.Sprintf("result_%d", i)] = r
			}
			msg.Details = details
		}

	} else {
		msg.Msg = "not found"
	}
	a.s.JSON(w, msg)
}

func (a MainApi) GetValidateStage(w http.ResponseWriter, r *http.Request) {
	if a.Session.CheckAuth(r) {
		a.s.ErrJSON(w, http.StatusForbidden, "required authentication headers")
		return
	}
	output := r.URL.Query().Get("output")
	hashID := r.PathValue("hashid")
	if output == "json" {
		id, err := strconv.Atoi(hashID)
		if err != nil {
			a.s.ErrJSON(w, http.StatusBadRequest, "id should be a integer")
			return
		}
		a.logger.Info("output json set and looking for id", "id", id)
		slice := a.Validator.Slice()
		for _, v := range slice {
			if v.ID == id {
				a.s.JSON(w, v)
				return
			}
		}
	}
	msg := types.Composed{}
	v, ok := a.Validator.Request[hashID]
	if ok {
		msg.Msg = fmt.Sprintf("%s with result: %v", v.Output, v.Valid)
		if len(v.Analise.Results) > 0 {
			details := make(map[string]string)
			for i, r := range v.Analise.Results {
				details[fmt.Sprintf("result_%d", i)] = r
			}
			msg.Details = details
		}

	} else {
		msg.Msg = "not found"
	}
	a.s.JSON(w, msg)
}

func (a MainApi) GetValidateStory(w http.ResponseWriter, r *http.Request) {
	if a.Session.CheckAuth(r) {
		a.s.ErrJSON(w, http.StatusForbidden, "required authentication headers")
		return
	}
	output := r.URL.Query().Get("output")
	hashID := r.PathValue("hashid")
	if output == "json" {
		id, err := strconv.Atoi(hashID)
		if err != nil {
			a.s.ErrJSON(w, http.StatusBadRequest, "id should be a integer")
			return
		}
		a.logger.Info("output json set and looking for id", "id", id)
		slice := a.Validator.Slice()
		for _, v := range slice {
			if v.ID == id {
				a.s.JSON(w, v)
				return
			}
		}
	}
	msg := types.Composed{}
	v, ok := a.Validator.Request[hashID]
	if ok {
		msg.Msg = fmt.Sprintf("%s with result: %v", v.Output, v.Valid)
		if len(v.Analise.Results) > 0 {
			details := make(map[string]string)
			for i, r := range v.Analise.Results {
				details[fmt.Sprintf("result_%d", i)] = r
			}
			msg.Details = details
		}

	} else {
		msg.Msg = "not found"
	}
	a.s.JSON(w, msg)
}

func (a MainApi) RequestToValidateAutoPlay(w http.ResponseWriter, r *http.Request) {
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
	autoPlay, err := a.db.GetAutoPlayByID(a.ctx, id)
	if err != nil {
		a.s.ErrJSON(w, http.StatusBadRequest, "error getting auto play")
		return
	}
	var hashID string
	msg := "auto play not found"
	if autoPlay.ID != 0 {
		hashID = utils.RandomString(8)
		go func() {
			a.Validator.AddRequest(autoPlay.ID, hashID, validator.KindAutoPlay)
			a.Validator.ValidateAutoPlay(&autoPlay, hashID)
		}()
		msg = fmt.Sprintf("hash_id %v", hashID)
	}
	a.s.JSON(w, types.Msg{Msg: msg})
}

func (a MainApi) RequestToValidateStage(w http.ResponseWriter, r *http.Request) {
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
	stage, err := a.db.GetStageByStageID(a.ctx, id)
	if err != nil {
		a.s.ErrJSON(w, http.StatusBadRequest, "error getting stage")
		return
	}
	var hashID string
	msg := "stage not found"
	if stage.Stage.ID != 0 {
		hashID = utils.RandomString(8)
		go func() {
			a.Validator.AddRequest(stage.Stage.ID, hashID, validator.KindStage)
			a.Validator.ValidateStage(&stage, hashID)
		}()
		msg = fmt.Sprintf("hash_id %v", id)
	}
	a.s.JSON(w, types.Msg{Msg: msg})
}

func (a MainApi) RequestToValidateStory(w http.ResponseWriter, r *http.Request) {
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

	story, err := a.db.GetStoryByID(a.ctx, id)
	if err != nil {
		a.s.ErrJSON(w, http.StatusBadRequest, "error getting story")
		return
	}
	var hashID string
	msg := "story not found"
	if story.ID != 0 {
		hashID = utils.RandomString(8)
		go func() {
			a.Validator.AddRequest(story.ID, hashID, validator.KindStory)
			a.Validator.ValidateStory(&story, hashID)
		}()
		msg = fmt.Sprintf("hash_id %v", hashID)
	}
	a.s.JSON(w, types.Msg{Msg: msg})
}

package v1

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/betorvs/playbypost/core/rpg"
	"github.com/betorvs/playbypost/core/sys/web/types"
	"github.com/betorvs/playbypost/core/utils"
)

func (a MainApi) GeneratePlayer(w http.ResponseWriter, r *http.Request) {
	if a.Session.CheckAuth(r) {
		a.s.ErrJSON(w, http.StatusForbidden, "required authentication headers")
		return
	}
	obj := types.GeneratePlayer{}
	err := json.NewDecoder(r.Body).Decode(&obj)
	if err != nil {
		a.s.ErrJSON(w, http.StatusBadRequest, "json decode error")
		return
	}
	stage, err := a.db.GetStageByStageID(a.ctx, obj.StageID)
	if err != nil {
		a.s.ErrJSON(w, http.StatusBadRequest, "stage database issue")
		return
	}
	// check user_id storyteller cannot be player
	if obj.UserID == stage.Stage.UserID {
		a.s.ErrJSON(w, http.StatusBadRequest, "storyteller cannot be a player")
		return
	}
	if stage.Stage.ID == 0 {
		a.s.ErrJSON(w, http.StatusBadRequest, "stage cannot be empty")
		return
	}
	if obj.UserID == "" {
		a.s.ErrJSON(w, http.StatusBadRequest, "user_id cannot be empty")
		return
	}
	user, err := a.db.GetUserByUserID(a.ctx, obj.UserID)
	if err != nil {
		a.s.ErrJSON(w, http.StatusBadRequest, "users database issue")
		return
	}
	userID := user.ID
	if userID == 0 {
		id, err := a.db.CreateUserTx(a.ctx, obj.UserID)
		if err != nil {
			a.s.ErrJSON(w, http.StatusBadRequest, "error adding user to database")
			return
		}
		userID = id
	}

	switch a.rpg.Name {
	case rpg.D10HM:
		creature, err := utils.GenD10Random(obj.Name, a.rpg)
		if err != nil {
			a.logger.Error("generate createure", "error", err.Error())
			a.s.ErrJSON(w, http.StatusBadRequest, "cannot generate a random player")
			return
		}
		// creature.
		res, err := a.db.SavePlayerTx(a.ctx, userID, obj.StageID, false, creature, a.rpg)
		if err != nil {
			a.logger.Error("generate creature", "userid", userID, "error", err.Error())
			a.s.ErrJSON(w, http.StatusBadGateway, "error saving new player")
			return
		}
		msg := fmt.Sprintf("player id %v", res)
		a.s.JSON(w, types.Msg{Msg: msg})
		return

	default:
		a.s.ErrJSON(w, http.StatusNotImplemented, "not implemented")
		return
	}
}

func (a MainApi) GetPlayersByStageID(w http.ResponseWriter, r *http.Request) {
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
	npc := false
	if r.URL.Query().Get("npc") == "true" {
		npc = true
	}
	a.logger.Info("get players by story id", "story-id", id, "query_npc", npc)
	obj, err := a.db.GetPlayerByStageID(a.ctx, id)
	if err != nil {
		a.s.ErrJSON(w, http.StatusBadRequest, "players database issue")
		return
	}
	a.logger.Info("players list", "obj", obj)
	a.s.JSON(w, obj)
}

func (a MainApi) GetPlayersByID(w http.ResponseWriter, r *http.Request) {
	idString := r.PathValue("id")
	id, err := strconv.Atoi(idString)
	if err != nil {
		a.s.ErrJSON(w, http.StatusBadRequest, "id should be a integer")
		return
	}
	npc := false
	if r.URL.Query().Get("npc") == "true" {
		npc = true
	}
	a.logger.Info("get players by id", "player-id", id, "query_npc", npc)
	obj, err := a.db.GetPlayerByID(a.ctx, id)
	if err != nil {
		a.s.ErrJSON(w, http.StatusBadRequest, "players database issue")
		return
	}
	a.s.JSON(w, obj)
}

func (a MainApi) GetPlayers(w http.ResponseWriter, r *http.Request) {
	if a.Session.CheckAuth(r) {
		a.s.ErrJSON(w, http.StatusForbidden, "required authentication headers")
		return
	}
	obj, err := a.db.GetPlayers(a.ctx)
	if err != nil {
		a.s.ErrJSON(w, http.StatusBadRequest, "players database issue")
		return
	}
	a.s.JSON(w, obj)
}

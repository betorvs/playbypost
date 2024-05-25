package handlers

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
	if a.checkAuth(r) {
		a.s.ErrJSON(w, http.StatusForbidden, "required authentication headers")
		return
	}
	obj := types.GeneratePlayer{}
	err := json.NewDecoder(r.Body).Decode(&obj)
	if err != nil {
		a.s.ErrJSON(w, http.StatusBadRequest, "json decode error")
		return
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
		res, err := a.db.SavePlayerTx(a.ctx, obj.PlayerID, obj.StoryID, false, creature, a.rpg)
		if err != nil {
			a.logger.Error("generate createure", "error", err.Error())
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

func (a MainApi) GetPlayersByStoryID(w http.ResponseWriter, r *http.Request) {
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
	npc := false
	if r.URL.Query().Get("npc") == "true" {
		npc = true
	}
	a.logger.Info("get players by story id", "story-id", id, "query_npc", npc)
	obj, err := a.db.GetSliceOfPlayersByStoryID(a.ctx, id, npc, a.rpg)
	if err != nil {
		a.s.ErrJSON(w, http.StatusBadRequest, "players database issue")
		return
	}
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
	obj, err := a.db.GetPlayer(a.ctx, id, npc, a.rpg)
	if err != nil {
		a.s.ErrJSON(w, http.StatusBadRequest, "players database issue")
		return
	}
	a.s.JSON(w, obj)
}

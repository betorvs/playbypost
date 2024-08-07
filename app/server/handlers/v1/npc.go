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

func (a MainApi) GetNPCByStageID(w http.ResponseWriter, r *http.Request) {
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
	a.logger.Info("get npc by story id", "story-id", id)
	obj, err := a.db.GetNPCByStageID(a.ctx, id)
	if err != nil {
		a.s.ErrJSON(w, http.StatusBadRequest, "npc database issue")
		return
	}
	a.logger.Info("npc list", "obj", obj)
	a.s.JSON(w, obj)
}

func (a MainApi) GenerateNPC(w http.ResponseWriter, r *http.Request) {
	if a.Session.CheckAuth(r) {
		a.s.ErrJSON(w, http.StatusForbidden, "required authentication headers")
		return
	}
	obj := types.GenerateNPC{}
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
	// check storytellerID with stage storyteller ID
	if stage.Stage.StorytellerID != obj.StorytellerID {
		a.s.ErrJSON(w, http.StatusBadRequest, "only storyteller can generate a npc")
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
		npc, err := a.db.GenerateNPC(a.ctx, obj.Name, stage.Stage.ID, stage.Stage.StorytellerID, creature)
		if err != nil {
			a.s.ErrJSON(w, http.StatusBadRequest, "npc database issue")
			return
		}
		if obj.EncounterID != 0 {
			err = a.db.AddParticipants(a.ctx, obj.EncounterID, true, []int{npc})
			if err != nil {
				a.s.ErrJSON(w, http.StatusBadRequest, "error adding participants to encounter on database")
				return
			}
		}

		msg := fmt.Sprintf("npc id %v", npc)

		a.s.JSON(w, types.Msg{Msg: msg})
		return

	default:
		a.s.ErrJSON(w, http.StatusNotImplemented, "not implemented")
		return
	}
}

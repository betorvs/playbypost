package v1

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/betorvs/playbypost/core/parser"
	"github.com/betorvs/playbypost/core/sys/web/types"
)

func (a MainApi) GenerateInitiative(w http.ResponseWriter, r *http.Request) {
	if a.Session.CheckAuth(r) {
		a.s.ErrJSON(w, http.StatusForbidden, "required authentication headers")
		return
	}
	obj := types.Initiative{}
	err := json.NewDecoder(r.Body).Decode(&obj)
	if err != nil {
		a.s.ErrJSON(w, http.StatusBadRequest, "json decode error")
		return
	}
	fmt.Printf("%+v \n", obj)
	enc, err := a.db.GetStageEncounterByEncounterID(a.ctx, obj.EncounterID)
	if err != nil {
		a.s.ErrJSON(w, http.StatusBadRequest, "stage encounter not found")
		return
	}

	if types.Phase(enc.Phase) == types.Running {
		actions := types.NewActions()
		actions["command"] = parser.RollInitiative
		actions["text"] = "created by web handler"
		actions["channel"] = obj.Channel
		actions["userid"] = obj.UserID
	}
	a.s.ErrJSON(w, http.StatusBadRequest, "encounter not in running phase")
}

func (a MainApi) GetInitiativeByEncounterId(w http.ResponseWriter, r *http.Request) {
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
	obj, initID, err := a.db.GetRunningInitiativeByEncounterID(a.ctx, id)
	if err != nil {
		a.s.ErrJSON(w, http.StatusBadRequest, "initiative issue")
		return
	}
	if initID == -1 {
		a.s.ErrJSON(w, http.StatusNotFound, "initiative not found")
		return
	}
	info := obj.NextInfo()
	party := []string{}
	for _, v := range obj.Participants {
		result := fmt.Sprintf("%s, initiative score of %d", v.Name, v.Result)
		party = append(party, result)
	}
	initiative := types.InitiativeShort{
		ID:           initID,
		Name:         obj.Name,
		NextPlayer:   obj.Participants[info].Name,
		Participants: party,
	}

	a.s.JSON(w, initiative)
}

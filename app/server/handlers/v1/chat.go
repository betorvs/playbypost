package v1

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/betorvs/playbypost/core/sys/web/types"
	"github.com/lib/pq"
)

func (a MainApi) AddChatInfo(w http.ResponseWriter, r *http.Request) {
	if a.Session.CheckAuth(r) {
		a.s.ErrJSON(w, http.StatusForbidden, "required authentication headers")
		return
	}
	obj := types.ChatInfo{}
	err := json.NewDecoder(r.Body).Decode(&obj)
	if err != nil {
		a.s.ErrJSON(w, http.StatusBadRequest, "json decode error")
		return
	}
	if obj.UserID == "" || obj.Username == "" || obj.Channel == "" || obj.Chat == "" {
		a.s.ErrJSON(w, http.StatusBadRequest, "empty body")
		return
	}
	a.logger.Info("chat info received", "username", obj.Username, "userid", obj.UserID, "channel", obj.Channel, "chat", obj.Chat)
	_, err = a.db.AddChatInformation(a.ctx, obj.Username, obj.UserID, obj.Channel, obj.Chat)
	if err != nil {
		pgErr, ok := err.(*pq.Error)
		if ok {
			if pgErr.Code == "23505" {
				a.s.JSON(w, types.Msg{Msg: "already added"})
				return
			}
		}
		a.s.ErrJSON(w, http.StatusBadRequest, "error adding slack info to database")
		return
	}
	a.s.JSON(w, types.Msg{Msg: "added"})
}

func (a MainApi) GetUsersInformation(w http.ResponseWriter, r *http.Request) {
	// if a.Session.CheckAuth(r) {
	// 	a.s.ErrJSON(w, http.StatusForbidden, "required authentication headers")
	// 	return
	// }
	obj, err := a.db.GetChatInformation(a.ctx)
	if err != nil {
		a.s.ErrJSON(w, http.StatusBadRequest, "slack information for users database issue")
		return
	}
	a.s.JSON(w, obj)
}

func (a MainApi) GetChannelsInformation(w http.ResponseWriter, r *http.Request) {
	// if a.Session.CheckAuth(r) {
	// 	a.s.ErrJSON(w, http.StatusForbidden, "required authentication headers")
	// 	return
	// }
	obj, err := a.db.GetChatChannelInformation(a.ctx)
	if err != nil {
		a.s.ErrJSON(w, http.StatusBadRequest, "slack information for channel database issue")
		return
	}
	a.s.JSON(w, obj)
}

func (a MainApi) GetEncountersPhase(w http.ResponseWriter, r *http.Request) {
	// if a.Session.CheckAuth(r) {
	// 	a.s.ErrJSON(w, http.StatusForbidden, "required authentication headers")
	// 	return
	// }
	det := make(map[string]string)
	for i := 0; i <= int(types.Finished); i++ {
		det[fmt.Sprintf("%d", i)] = types.PhaseAtoi(i).String()
	}
	obj := types.Composed{Msg: "Encounter Phases", Details: det}
	a.s.JSON(w, obj)
}

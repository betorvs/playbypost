package handlers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/betorvs/playbypost/core/sys/web/types"
)

func (a MainApi) ExecuteCommand(w http.ResponseWriter, r *http.Request) {
	headerUserID := r.Header.Get(types.HeaderUserID)
	headerStoryChannel := r.Header.Get(types.HeaderStoryChannel)
	if headerUserID == "" || headerStoryChannel == "" {
		a.s.ErrJSON(w, http.StatusForbidden, "required authentication headers")
		return
	}
	// scene, err := a.db.GetSceneByChannelID(a.ctx, headerStoryChannel)
	// if err != nil {
	// 	a.s.ErrJSON(w, http.StatusBadRequest, "channel id")
	// 	return
	// }
	id, err := strconv.Atoi(headerUserID)
	if err != nil {
		a.logger.Error("command strconv error", "error", err.Error())
		a.s.ErrJSON(w, http.StatusBadRequest, "bad headers")
		return
	}
	player, err := a.db.GetPlayerByUserID(a.ctx, id, false, a.rpg)
	if err != nil {
		a.s.ErrJSON(w, http.StatusBadRequest, "invalid user id")
		return
	}
	// return message
	// msg := fmt.Sprintf("player found '%s' and story id found '%d' ", player.Name, scene.Story.ID)
	msg := fmt.Sprintf("player found '%s' ", player.Name)
	//
	a.s.JSON(w, types.Msg{Msg: msg})
}

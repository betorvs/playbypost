package v1

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/betorvs/playbypost/core/sys/web/types"
)

type UpdateCharacterRequest struct {
	Name string `json:"name"`
	Rpg  string `json:"rpg"`
}

func (a MainApi) GetCharacters(w http.ResponseWriter, r *http.Request) {
	if a.Session.CheckAuth(r) {
		a.s.ErrJSON(w, http.StatusForbidden, "required authentication headers")
		return
	}
	username := r.Header.Get("X-Username")
	writer, err := a.db.GetWriterByUsername(a.ctx, username)
	if err != nil {
		a.s.ErrJSON(w, http.StatusBadRequest, "writer not found")
		return
	}
	obj, err := a.db.GetWriterUsersAssociation(a.ctx)
	if err != nil {
		a.s.ErrJSON(w, http.StatusBadRequest, "users database issue")
		return
	}
	// filter by writer id
	filtered := []types.WriterUserAssociation{}
	for _, v := range obj {
		if v.WriterID == writer.ID {
			filtered = append(filtered, v)
		}
	}
	// filtered length 0, return empty array: if theres no association, return empty array
	if len(filtered) == 0 {
		a.s.JSON(w, []types.Players{})
		return
	}

	// GetUserByUserID
	// characters := []types.User{}
	characters := []types.Players{}
	for _, v := range filtered {
		a.logger.Info("getting player by user id", "user id", v.UserID)
		players, err := a.db.GetPlayerByUserID(a.ctx, v.UserID, a.rpg)
		if err != nil {
			a.s.ErrJSON(w, http.StatusBadRequest, "player not found")
			return
		}
		characters = append(characters, players)
	}
	a.s.JSON(w, characters)
}

func (a MainApi) UpdateCharacter(w http.ResponseWriter, r *http.Request) {
	if a.Session.CheckAuth(r) {
		a.s.ErrJSON(w, http.StatusForbidden, "required authentication headers")
		return
	}

	idString := r.PathValue("id")
	if idString == "" {
		a.s.ErrJSON(w, http.StatusBadRequest, "id cannot be empty")
		return
	}
	characterID, err := strconv.Atoi(idString)
	if err != nil {
		a.s.ErrJSON(w, http.StatusBadRequest, "invalid id")
		return
	}

	var req UpdateCharacterRequest
	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		a.s.ErrJSON(w, http.StatusBadRequest, "invalid request body")
		return
	}

	err = a.db.UpdatePlayerDetails(a.ctx, characterID, req.Name, req.Rpg)
	if err != nil {
		a.s.ErrJSON(w, http.StatusInternalServerError, "failed to update character")
		return
	}

	a.s.JSON(w, map[string]string{"message": "character updated successfully"})
}

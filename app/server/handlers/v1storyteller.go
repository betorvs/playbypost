package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/betorvs/playbypost/core/sys/web/types"
	"golang.org/x/crypto/bcrypt"
)

func (a MainApi) GetStorytellers(w http.ResponseWriter, r *http.Request) {
	if a.checkAuth(r) {
		a.s.ErrJSON(w, http.StatusForbidden, "required authentication headers")
		return
	}
	obj, err := a.db.GetStorytellers(a.ctx, false)
	if err != nil {
		a.s.ErrJSON(w, http.StatusBadRequest, "users database issue")
		return
	}
	a.s.JSON(w, obj)
}

func (a MainApi) CreateStorytellers(w http.ResponseWriter, r *http.Request) {
	if a.checkAuth(r) {
		a.s.ErrJSON(w, http.StatusForbidden, "required authentication headers")
		return
	}
	obj := types.Storyteller{}
	err := json.NewDecoder(r.Body).Decode(&obj)
	if err != nil {
		a.s.ErrJSON(w, http.StatusBadRequest, "json decode error")
		return
	}
	if obj.Username == "" {
		a.s.ErrJSON(w, http.StatusBadRequest, "username and userid cannot be empty")
		return
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(obj.Password), 8)
	if err != nil {
		a.s.ErrJSON(w, http.StatusBadRequest, "error creating password hash")
		return
	}
	//utils.RandomString(16)
	res, err := a.db.CreateStorytellers(a.ctx, obj.Username, string(hashedPassword))
	if err != nil {
		a.s.ErrJSON(w, http.StatusBadRequest, "error creating user on database")
		return
	}
	msg := fmt.Sprintf("user id %v", res)
	a.s.JSON(w, types.Msg{Msg: msg})
}

func (a MainApi) GetUsersCard(w http.ResponseWriter, r *http.Request) {
	// if a.checkAuth(r) {
	// 	a.s.ErrJSON(w, http.StatusForbidden, "required authentication headers")
	// 	return
	// }
	obj, err := a.db.GetUserCard(a.ctx)
	if err != nil {
		a.s.ErrJSON(w, http.StatusBadRequest, "users database issue")
		return
	}
	a.s.JSON(w, obj)
}

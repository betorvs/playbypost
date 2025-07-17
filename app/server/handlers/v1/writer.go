package v1

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/betorvs/playbypost/core/sys/web/types"
	"golang.org/x/crypto/bcrypt"
)

func (a MainApi) GetWriters(w http.ResponseWriter, r *http.Request) {
	if a.Session.CheckAuth(r) {
		a.s.ErrJSON(w, http.StatusForbidden, "required authentication headers")
		return
	}
	obj, err := a.db.GetWriters(a.ctx, false)
	if err != nil {
		a.s.ErrJSON(w, http.StatusBadRequest, "users database issue")
		return
	}
	a.s.JSON(w, obj)
}

func (a MainApi) CreateWriters(w http.ResponseWriter, r *http.Request) {
	if a.Session.CheckAuth(r) {
		a.s.ErrJSON(w, http.StatusForbidden, "required authentication headers")
		return
	}
	obj := types.Writer{}
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
	res, err := a.db.CreateWriters(a.ctx, obj.Username, string(hashedPassword))
	if err != nil {
		a.s.ErrJSON(w, http.StatusBadRequest, "error creating user on database")
		return
	}
	msg := fmt.Sprintf("user id %v", res)
	a.s.JSON(w, types.Msg{Msg: msg})
}

func (a MainApi) CreateWriterUserAssociation(w http.ResponseWriter, r *http.Request) {
	if a.Session.CheckAuth(r) {
		a.s.ErrJSON(w, http.StatusForbidden, "required authentication headers")
		return
	}
	obj := types.WriterUserAssociation{}
	err := json.NewDecoder(r.Body).Decode(&obj)
	if err != nil {
		a.s.ErrJSON(w, http.StatusBadRequest, "json decode error")
		return
	}
	if obj.WriterID == 0 || obj.UserID == 0 {
		a.s.ErrJSON(w, http.StatusBadRequest, "writer_id and user_id cannot be empty")
		return
	}
	res, err := a.db.CreateWriterUserAssociation(a.ctx, obj.WriterID, obj.UserID)
	if err != nil {
		a.s.ErrJSON(w, http.StatusBadRequest, "error creating writer user association on database")
		return
	}
	msg := fmt.Sprintf("writer user association id %v", res)
	a.s.JSON(w, types.Msg{Msg: msg})
}

func (a MainApi) DeleteWriterUserAssociation(w http.ResponseWriter, r *http.Request) {
	if a.Session.CheckAuth(r) {
		a.s.ErrJSON(w, http.StatusForbidden, "required authentication headers")
		return
	}
	idString := r.PathValue("id")
	if idString == "" {
		a.s.ErrJSON(w, http.StatusBadRequest, "id cannot be empty")
		return
	}
	id, err := strconv.Atoi(idString)
	if err != nil {
		a.s.ErrJSON(w, http.StatusBadRequest, "invalid id")
		return
	}
	err = a.db.DeleteWriterUserAssociation(a.ctx, id)
	if err != nil {
		a.s.ErrJSON(w, http.StatusBadRequest, "error deleting writer user association")
		return
	}
	a.s.JSON(w, types.Msg{Msg: "writer user association deleted successfully"})
}

func (a MainApi) GetWriterUsersAssociation(w http.ResponseWriter, r *http.Request) {
	if a.Session.CheckAuth(r) {
		a.s.ErrJSON(w, http.StatusForbidden, "required authentication headers")
		return
	}
	associations, err := a.db.GetWriterUsersAssociation(a.ctx)
	if err != nil {
		a.s.ErrJSON(w, http.StatusBadRequest, "error getting writer user association")
		return
	}
	a.s.JSON(w, associations)
}

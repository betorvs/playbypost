package v1

import "net/http"

func (a MainApi) GetUsersByUserID(w http.ResponseWriter, r *http.Request) {
	if a.Session.CheckAuth(r) {
		a.s.ErrJSON(w, http.StatusForbidden, "required authentication headers")
		return
	}
	idString := r.PathValue("id")
	if idString == "" {
		a.s.ErrJSON(w, http.StatusBadRequest, "id cannot be empty")
		return
	}
	user, err := a.db.GetUserByUserID(a.ctx, idString)
	if err != nil {
		a.logger.Error("get users by user id", "error", err.Error())
		a.s.ErrJSON(w, http.StatusBadRequest, "users database issue")
		return
	}
	if user.ID == 0 {
		a.s.ErrJSON(w, http.StatusBadRequest, "user not found")
		return
	}
	a.s.JSON(w, user)
}

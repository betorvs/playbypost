package handlers

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"time"

	"github.com/betorvs/playbypost/core/sys/db"
	"github.com/betorvs/playbypost/core/sys/web/server"
	"github.com/betorvs/playbypost/core/sys/web/types"
	"github.com/betorvs/playbypost/core/utils"
	"golang.org/x/crypto/bcrypt"
)

type Session struct {
	admin  string
	logger *slog.Logger
	db     db.DBClient
	s      *server.SvrWeb
	ctx    context.Context
}

func NewSession(admin string, logger *slog.Logger, db db.DBClient, s *server.SvrWeb, ctx context.Context) *Session {
	return &Session{
		admin:  admin,
		logger: logger,
		db:     db,
		s:      s,
		ctx:    ctx,
	}
}

func (a *Session) AddAdminSession(admin, token string) {
	expiresAt := time.Now().Add(8760 * time.Hour)
	session := types.Session{
		Username:     admin,
		Token:        token,
		Expiry:       expiresAt,
		ClientType:   "admin-ctl",
		ClientInfo:   "local",
		IPAddress:    "localhost",
		UserAgent:    "internal",
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
		LastActivity: time.Now(),
	}
	err := a.db.CreateSession(a.ctx, session)
	if err != nil {
		a.logger.Error("failed to create admin session", "error", err)
	}
}

func (a *Session) CheckAuth(r *http.Request) bool {
	headerToken := r.Header.Get(types.HeaderToken)
	if headerToken == "" {
		return true
	}
	session, err := a.db.GetSessionByToken(a.ctx, headerToken)
	if err != nil {
		return true
	}
	if session.IsExpired() {
		return true
	}
	return false
}

func (a *Session) Signin(w http.ResponseWriter, r *http.Request) {
	var creds types.Credentials
	err := json.NewDecoder(r.Body).Decode(&creds)
	if err != nil {
		a.s.ErrJSON(w, http.StatusBadRequest, "json issue")
		return
	}

	user, err := a.db.GetWriterByUsername(a.ctx, creds.Username)
	if err != nil {
		a.s.ErrJSON(w, http.StatusBadRequest, "user not found")
		return
	}

	if err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(creds.Password)); err != nil {
		// If the two passwords don't match, return a 401 status
		a.s.ErrJSON(w, http.StatusUnauthorized, "username or password does not match")
		return
	}

	sessionToken := utils.RandomString(48) // uuid.NewString()
	expiresAt := time.Now().Add(3000 * time.Second)

	session := types.Session{
		Username: creds.Username,
		Token:    sessionToken,
		Expiry:   expiresAt,
		UserID:   user.ID,
		ClientType:   "web",
		ClientInfo:   r.RemoteAddr,
		IPAddress:    r.RemoteAddr,
		UserAgent:    r.UserAgent(),
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
		LastActivity: time.Now(),
	}
	err = a.db.CreateSession(a.ctx, session)
	if err != nil {
		a.logger.Error("failed to create session", "error", err)
		a.s.ErrJSON(w, http.StatusInternalServerError, "failed to create session")
		return
	}

	login := types.Login{
		Status:      "ok",
		Message:     "logged in",
		AccessToken: sessionToken,
		ExpireOn:    expiresAt,
		UserID:      user.ID,
	}
	a.s.JSON(w, login)
}

func (a *Session) Logout(w http.ResponseWriter, r *http.Request) {
	headerToken := r.Header.Get(types.HeaderToken)
	if headerToken == "" {
		a.s.ErrJSON(w, http.StatusUnauthorized, "missing token")
		return
	}
	err := a.db.DeleteSessionByToken(a.ctx, headerToken)
	if err != nil {
		a.logger.Error("failed to delete session", "error", err)
		a.s.ErrJSON(w, http.StatusInternalServerError, "failed to delete session")
		return
	}
	login := types.Login{
		Status:  "ok",
		Message: "logged out",
	}
	a.s.JSON(w, login)
}

func (a Session) ValidateSession(w http.ResponseWriter, r *http.Request) {
	if a.CheckAuth(r) {
		a.s.ErrJSON(w, http.StatusForbidden, "required authentication headers")
		return
	}
	a.s.JSON(w, types.Msg{Msg: "authenticated"})
}

func (s Session) Admin() string {
	return s.admin
}

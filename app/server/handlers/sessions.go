package handlers

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"sync"
	"time"

	"github.com/betorvs/playbypost/core/sys/db"
	"github.com/betorvs/playbypost/core/sys/web/server"
	"github.com/betorvs/playbypost/core/sys/web/types"
	"github.com/betorvs/playbypost/core/utils"
	"golang.org/x/crypto/bcrypt"
)

type Session struct {
	logger   *slog.Logger
	db       db.DBClient
	s        *server.SvrWeb
	ctx      context.Context
	Sessions Sessions
}

type Sessions struct {
	Current map[string]types.Session
	mu      *sync.Mutex
}

func NewSession(logger *slog.Logger, db db.DBClient, s *server.SvrWeb, ctx context.Context) *Session {
	return &Session{
		logger: logger,
		db:     db,
		s:      s,
		ctx:    ctx,
		Sessions: Sessions{
			Current: make(map[string]types.Session),
			mu:      &sync.Mutex{},
		},
	}
}

func (m *Sessions) Add(index string, value types.Session) {
	m.mu.Lock()
	m.Current[index] = value
	m.mu.Unlock()
}

func (m *Sessions) Remove(index string) {
	m.mu.Lock()
	delete(m.Current, index)
	m.mu.Unlock()
}

func (a *Session) AddAdminSession(admin, token string) {
	expiresAt := time.Now().Add(8760 * time.Hour)
	session := types.Session{
		Username: admin,
		Token:    token,
		Expiry:   expiresAt,
	}
	a.Sessions.Add(admin, session)
}

func (a *Session) CheckAuth(r *http.Request) bool {

	headerToken := r.Header.Get(types.HeaderToken)
	headerUsername := r.Header.Get(types.HeaderUsername)
	if headerToken != "" && headerUsername != "" {
		v, ok := a.Sessions.Current[headerUsername]
		if !ok || v.IsExpired() || headerToken != v.Token {
			return true
		}
		return false
	}

	return true
}

func (a Session) Signin(w http.ResponseWriter, r *http.Request) {
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

	s, ok := a.Sessions.Current[creds.Username]
	if ok && !s.IsExpired() {
		sessionToken = s.Token
		expiresAt = s.Expiry
		a.logger.Info("user already logged in")
	} else {
		a.logger.Info("login added")
		session := types.Session{
			Username: creds.Username,
			Token:    sessionToken,
			Expiry:   expiresAt,
			UserID:   user.ID,
			// EncodingKey: user.EncodingKey,
		}
		a.Sessions.Add(creds.Username, session)
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

func (a Session) Logout(w http.ResponseWriter, r *http.Request) {
	if a.CheckAuth(r) {
		a.s.ErrJSON(w, http.StatusForbidden, "required authentication headers")
		return
	}
	headerToken := r.Header.Get(types.HeaderToken)
	headerUsername := r.Header.Get(types.HeaderUsername)
	s, ok := a.Sessions.Current[headerUsername]
	if headerToken == s.Token && ok {
		a.Sessions.Remove(headerUsername)
	} else {
		a.s.ErrJSON(w, http.StatusBadRequest, "user not found")
		return
	}
	login := types.Login{
		Status:  "ok",
		Message: "logged off",
	}
	a.s.JSON(w, login)
}

func (a Session) Refresh(w http.ResponseWriter, r *http.Request) {
	if a.CheckAuth(r) {
		a.s.ErrJSON(w, http.StatusForbidden, "required authentication headers")
		return
	}
	headerUsername := r.Header.Get(types.HeaderUsername)
	s, ok := a.Sessions.Current[headerUsername]
	if !ok {
		a.s.ErrJSON(w, http.StatusUnauthorized, "require login first")
		return
	}
	userID := s.UserID
	if s.IsExpired() {
		a.Sessions.Remove(headerUsername)
		a.s.ErrJSON(w, http.StatusUnauthorized, "token expired, do log in again")
		return
	}
	newSessionToken := utils.RandomString(48) // uuid.NewString()
	expiresAt := time.Now().Add(300 * time.Second)

	session := types.Session{
		Username: headerUsername,
		Token:    newSessionToken,
		Expiry:   expiresAt,
		UserID:   userID,
	}
	a.Sessions.Add(headerUsername, session)

	login := types.Login{
		Status:      "ok",
		Message:     "refresh okay",
		AccessToken: newSessionToken,
		ExpireOn:    expiresAt,
		UserID:      userID,
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

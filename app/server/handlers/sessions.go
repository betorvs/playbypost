package handlers

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/betorvs/playbypost/core/sys/db"
	"github.com/betorvs/playbypost/core/sys/web/server"
	"github.com/betorvs/playbypost/core/sys/web/types"
	"github.com/betorvs/playbypost/core/utils"
	"golang.org/x/crypto/bcrypt"
)

type SessionChecker interface {
	CheckAuth(r *http.Request) bool
	AddAdminSession(admin, token string)
	Admin() string
	GetActiveSessions() map[string]types.Session
	GetAllSessions(w http.ResponseWriter, r *http.Request)
	Signin(w http.ResponseWriter, r *http.Request)
	Logout(w http.ResponseWriter, r *http.Request)
	ValidateSession(w http.ResponseWriter, r *http.Request)
}

type Session struct {
	admin    string
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

func NewSession(admin string, logger *slog.Logger, db db.DBClient, s *server.SvrWeb, ctx context.Context) *Session {
	return &Session{
		admin:  admin,
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

func (m *Sessions) AddToCache(index string, value types.Session) {
	m.mu.Lock()
	m.Current[index] = value
	m.mu.Unlock()
}

func (m *Sessions) RemoveFromCache(index string) {
	m.mu.Lock()
	delete(m.Current, index)
	m.mu.Unlock()
}

func (a Session) AddAdminSession(admin, token string) {
	expiresAt := time.Now().Add(8760 * time.Hour)
	session := types.Session{
		Username:     admin,
		Token:        token,
		Expiry:       expiresAt,
		ClientType:   "admin-ctl",
		ClientInfo:   "{\"kind\": \"admin-ctl\", \"source\": \"local\"}",
		IPAddress:    "127.0.0.1",
		UserAgent:    "internal",
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
		LastActivity: time.Now(),
	}
	// admin does not have writer id, so we don't need to create a session in the database
	a.Sessions.AddToCache(token, session)
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

	remoteAddr := r.RemoteAddr
	if strings.HasPrefix(remoteAddr, "[::1]") {
		remoteAddr = "127.0.0.1"
	}
	clientType, clientInfo := getClientContext(r)
	session := types.Session{
		Username:     creds.Username,
		Token:        sessionToken,
		Expiry:       expiresAt,
		UserID:       user.ID,
		ClientType:   clientType,
		ClientInfo:   clientInfo,
		IPAddress:    remoteAddr,
		UserAgent:    r.UserAgent(),
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
		LastActivity: time.Now(),
	}
	a.logger.Info("creating session", "session", session)
	err = a.db.CreateSession(a.ctx, session)
	if err != nil {
		a.logger.Error("failed to create session", "error", err)
		a.s.ErrJSON(w, http.StatusInternalServerError, "failed to create session")
		return
	}
	// add session to cache after creating it in the database
	a.Sessions.AddToCache(sessionToken, session)

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
	headerToken := r.Header.Get(types.HeaderToken)
	if headerToken == "" {
		a.s.ErrJSON(w, http.StatusUnauthorized, "missing token")
		return
	}
	a.Sessions.RemoveFromCache(headerToken)
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

func (s Session) GetActiveSessions() map[string]types.Session {
	// loop through the sessions and remove token from the map
	for k, v := range s.Sessions.Current {
		v.Token = "REDACTED"
		s.Sessions.Current[k] = v
	}
	return s.Sessions.Current
}

func (a Session) GetAllSessions(w http.ResponseWriter, r *http.Request) {
	if a.CheckAuth(r) {
		a.s.ErrJSON(w, http.StatusForbidden, "required authentication headers")
		return
	}
	sessions, err := a.db.GetAllSessions(a.ctx)
	if err != nil {
		a.logger.Error("failed to get all sessions", "error", err)
		a.s.ErrJSON(w, http.StatusInternalServerError, "failed to get all sessions")
		return
	}
	// loop through the sessions and remove token from the map
	for k, v := range sessions {
		v.Token = "REDACTED"
		sessions[k] = v
	}
	a.s.JSON(w, sessions)
}

// CheckAuth checks if the request is authenticated
// it returns true if the request is not authenticated
// it returns false if the request is authenticated
func (a Session) CheckAuth(r *http.Request) bool {
	headerToken := r.Header.Get(types.HeaderToken)
	// empty token means no authentication
	if headerToken == "" {
		a.logger.Debug("no token in request", "request", r)
		return true
	}
	// check if the token is in the cache
	v, ok := a.Sessions.Current[headerToken]
	if ok && !v.IsExpired() && headerToken == v.Token {
		a.logger.Debug("check auth in cache", "session", v)
		return false
	}
	// check if the token is in the database
	session, err := a.db.GetSessionByToken(a.ctx, headerToken)
	if err != nil {
		a.logger.Error("failed to get session by token", "error", err)
		return true
	}
	// check if the session is expired
	if session.IsExpired() {
		a.logger.Error("session expired", "session", session)
		return true
	}
	// add the session to the cache
	a.logger.Debug("adding session to cache", "session", session)
	a.Sessions.AddToCache(headerToken, session)
	return false
}

func getClientContext(r *http.Request) (clientType, clientInfo string) {
	userAgent := r.UserAgent()
	referer := r.Referer()

	switch {
	case strings.HasPrefix(userAgent, "Mozilla/5.0"):
		clientType = "web"
		clientInfo = `{"kind": "web", "source": "` + r.RemoteAddr + `", "referer": "` + referer + `"}`
	case strings.HasPrefix(userAgent, "Go-http-client"):
		clientType = "cli"
		clientInfo = `{"kind": "cli", "source": "local"}`
	default:
		clientType = "unknown"
		clientInfo = `{"kind": "unknown", "userAgent": "` + userAgent + `"}`
	}

	return clientType, clientInfo
}

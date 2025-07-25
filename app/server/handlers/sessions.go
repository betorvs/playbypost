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
	GetAllSessions(w http.ResponseWriter, r *http.Request)
	Signin(w http.ResponseWriter, r *http.Request)
	Logout(w http.ResponseWriter, r *http.Request)
	ValidateSession(w http.ResponseWriter, r *http.Request)
	CheckAuth(r *http.Request) bool
	AddAdminSession(admin, token string)
	Admin() string
	GetActiveSessions() map[string]types.Session
	DeleteSessionByID(ctx context.Context, sessionID int64) error
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

	remoteAddr := r.RemoteAddr
	if strings.HasPrefix(remoteAddr, "[::1]") {
		remoteAddr = "127.0.0.1"
	}

	user, err := a.db.GetWriterByUsername(a.ctx, creds.Username)
	if err != nil {
		// Log failed login attempt
		a.db.LogLoginAttempt(a.ctx, creds.Username, remoteAddr, r.UserAgent(), false)

		a.s.ErrJSON(w, http.StatusBadRequest, "user not found")
		return
	}

	if err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(creds.Password)); err != nil {
		// Log failed login attempt
		a.db.LogLoginAttempt(a.ctx, creds.Username, remoteAddr, r.UserAgent(), false)

		// If the two passwords don't match, return a 401 status
		a.s.ErrJSON(w, http.StatusUnauthorized, "username or password does not match")
		return
	}

	sessionToken := utils.RandomString(48) // uuid.NewString()
	expiresAt := time.Now().Add(3000 * time.Second)

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

	// Log successful login attempt
	a.db.LogLoginAttempt(a.ctx, creds.Username, remoteAddr, r.UserAgent(), true)

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

	// Get session details before deletion for event logging
	session, err := a.db.GetSessionByToken(a.ctx, headerToken)
	if err != nil {
		a.logger.Error("failed to get session for logout event", "error", err)
		// Continue with logout even if we can't get session details
	} else {
		// Log logout event
		a.db.LogLogout(a.ctx, session.ID, session.Username)
	}

	err = a.db.DeleteSessionByToken(a.ctx, headerToken)
	if err != nil {
		a.logger.Error("failed to delete session", "error", err)
		a.s.ErrJSON(w, http.StatusInternalServerError, "failed to delete session")
		return
	}
	a.Sessions.RemoveFromCache(headerToken)
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
		// Log successful session validation from cache
		session, err := a.db.GetSessionByToken(a.ctx, headerToken)
		if err == nil {
			a.db.LogSessionValidated(a.ctx, session.ID, v.Username)
		}
		return false
	}
	// check if the token is in the database
	session, err := a.db.GetSessionByToken(a.ctx, headerToken)
	if err != nil {
		a.logger.Error("failed to get session by token", "error", err)
		// Log session invalid event
		a.db.LogSessionInvalid(a.ctx, 0, "session_not_found")
		return true
	}
	// check if the session is expired
	if session.IsExpired() {
		a.logger.Error("session expired", "session", session)
		// Log session invalid event
		a.db.LogSessionInvalid(a.ctx, session.ID, "session_expired")
		return true
	}
	// add the session to the cache
	a.logger.Debug("adding session to cache", "session", session)
	a.Sessions.AddToCache(headerToken, session)

	// Log successful session validation from database
	a.db.LogSessionValidated(a.ctx, session.ID, session.Username)
	return false
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

func (a Session) DeleteSessionByID(ctx context.Context, sessionID int64) error {
	session, err := a.db.GetSessionByID(ctx, sessionID)
	if err != nil {
		a.logger.Error("failed to get session by ID", "error", err)
		return err
	}
	err = a.db.DeleteSessionByID(ctx, sessionID)
	if err != nil {
		a.logger.Error("failed to delete session by ID", "error", err)
		return err
	}
	a.Sessions.RemoveFromCache(session.Token)
	return nil
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

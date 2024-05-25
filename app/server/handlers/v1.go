package handlers

import (
	"context"
	"log/slog"
	"net/http"
	"sync"
	"time"

	"github.com/betorvs/playbypost/core/rpg"
	"github.com/betorvs/playbypost/core/sys/db"
	"github.com/betorvs/playbypost/core/sys/web/server"
	"github.com/betorvs/playbypost/core/sys/web/types"
)

type MainApi struct {
	Sessions Sessions
	logger   *slog.Logger
	s        *server.SvrWeb
	db       db.DBClient
	ctx      context.Context
	dice     rpg.Roll
	rpg      *rpg.RPGSystem
	soloRPG  *rpg.RPGSystem
	didactic *rpg.RPGSystem
}

func NewMainApi(ctx context.Context, dice rpg.Roll, db db.DBClient, l *slog.Logger, s *server.SvrWeb, rpgSystem *rpg.RPGSystem) *MainApi {
	session := Sessions{
		Current: map[string]types.Session{},
		mu:      &sync.Mutex{},
	}
	return &MainApi{
		Sessions: session,
		ctx:      ctx,
		dice:     dice,
		db:       db,
		logger:   l,
		s:        s,
		rpg:      rpgSystem,
		soloRPG:  rpg.LoadRPGSystemsDefault(rpg.Solo, l),
		didactic: rpg.LoadRPGSystemsDefault(rpg.Didactic, l),
	}
}

type Sessions struct {
	Current map[string]types.Session
	mu      *sync.Mutex
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

func (a *MainApi) AddAdminSession(admin, token string) {
	expiresAt := time.Now().Add(8760 * time.Hour)
	session := types.Session{
		Username: admin,
		Token:    token,
		Expiry:   expiresAt,
	}
	a.Sessions.Add(admin, session)
}

func (a *MainApi) checkAuth(r *http.Request) bool {

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

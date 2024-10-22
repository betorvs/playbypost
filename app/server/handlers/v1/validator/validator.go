package validator

import (
	"context"
	"log/slog"
	"sync"
	"time"

	"github.com/betorvs/playbypost/core/sys/db"
)

const (
	KindStory    = "story"
	KindStage    = "stage"
	KindAutoPlay = "autoplay"
)

type Worker interface {
	Execute()
}

type Validator struct {
	Request map[string]Request
	ctx     context.Context
	logger  *slog.Logger
	db      db.DBClient
	mu      sync.Mutex
}

type Request struct {
	ID        int       `json:"id"`
	Kind      string    `json:"kind"`
	Valid     bool      `json:"valid"`
	Output    string    `json:"output"`
	CheckSum  string    `json:"checksum"`
	UpdatedAt time.Time `json:"updated_at"`
	Analise   Analitycs `json:"analise"`
}

type Analitycs struct {
	Results []string `json:"results"`
}

func New(logger *slog.Logger, db db.DBClient, ctx context.Context) *Validator {
	return &Validator{
		Request: make(map[string]Request),
		ctx:     ctx,
		logger:  logger,
		db:      db,
	}
}

func (v *Validator) AddRequest(id int, hashID, kind string) {
	v.mu.Lock()
	defer v.mu.Unlock()
	if _, ok := v.Request[hashID]; ok {
		return
	} else {
		v.Request[hashID] = Request{
			ID:        id,
			Kind:      kind,
			UpdatedAt: time.Now(),
			Analise: Analitycs{
				Results: []string{},
			},
		}
	}
}

func (v *Validator) UpdateRequest(hashID string, valid bool, output string, a Analitycs) {
	v.mu.Lock()
	defer v.mu.Unlock()
	if _, ok := v.Request[hashID]; ok {
		v.Request[hashID] = Request{
			ID:        v.Request[hashID].ID,
			Kind:      v.Request[hashID].Kind,
			Valid:     valid,
			Output:    output,
			UpdatedAt: time.Now(),
			Analise:   a,
		}
	}
}

func (v *Validator) RemoveOldRequests() {
	v.mu.Lock()
	defer v.mu.Unlock()
	for k, r := range v.Request {
		if time.Since(r.UpdatedAt) > 24*time.Hour {
			delete(v.Request, k)
		}
	}
}

func (v *Validator) Execute() {
	v.logger.Info("starting validator worker api execution", "time", time.Now())
	v.RemoveOldRequests()
}

func (v *Validator) Slice() []Request {
	result := []Request{}
	for k, value := range v.Request {
		i, ok := containsResultIDKind(result, value)
		if !ok {
			result = append(result, value)
			// result = slices.Replace(result, i, i, v)
		} else {
			if result[i].UpdatedAt.Before(v.Request[k].UpdatedAt) {
				result[i] = value
			}
		}
	}
	return result
}

func containsResultIDKind(r []Request, value Request) (int, bool) {
	for k, v := range r {
		if v.ID == value.ID && v.Kind == value.Kind && v.Valid == value.Valid {
			return k, true
		}
	}
	return -1, false
}

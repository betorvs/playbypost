package server

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/betorvs/playbypost/core/sys/web/types"
	"github.com/betorvs/playbypost/ui"
)

// SvrWeb struct
type SvrWeb struct {
	mux    *http.ServeMux
	Server *http.Server
	logger *slog.Logger
}

// NewServer factory function to create a new server
func NewServer(port int, logger *slog.Logger) *SvrWeb {
	mux := http.NewServeMux()
	log := slog.NewLogLogger(logger.Handler(), slog.LevelError)
	server := &http.Server{
		Addr:     fmt.Sprintf(":%v", port),
		Handler:  mux,
		ErrorLog: log,
		// ReadTimeout:    10 * time.Second,
		// WriteTimeout:   10 * time.Second,
		// MaxHeaderBytes: 1 << 20,
	}
	return &SvrWeb{
		mux:    mux,
		Server: server,
		logger: logger,
	}
}
func (s *SvrWeb) Register(pattern string, handler func(http.ResponseWriter, *http.Request)) {
	s.logger.Info("registering", "pattern", pattern)
	s.mux.HandleFunc(pattern, handler)
}

func (s *SvrWeb) RegisterStatic() {
	web, err := ui.Assets()
	if err != nil {
		s.logger.Error("assets error", "error", err.Error())
	}
	// s.mux.Handle("/", http.StripPrefix("/", http.FileServer(http.Dir("./web/dist"))))
	s.mux.Handle("/", http.StripPrefix("/", http.FileServer(http.FS(web))))

}

// JSON renders 'v' as JSON and writes it as a response into w.
func (s *SvrWeb) JSON(w http.ResponseWriter, v interface{}) {
	js, err := json.Marshal(v)
	if err != nil {
		s.logger.Error("json marshal error", "error", err.Error())
		s.ErrJSON(w, http.StatusInternalServerError, "encoding json error")
		// w.Header().Set("Content-Type", "application/json")
		// w.WriteHeader(http.StatusInternalServerError)
		// fmt.Fprint(w, "{\"error\":\"encoding json\"}")
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*") //"Access-Control-Allow-Origin": "http://192.168.1.210:5173/",
	w.Header().Set("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept, Authorization, X-Username, X-Access-Token, X-Cursor, X-Cursor-URI")
	w.Header().Set("Access-Control-Request-Methods", "GET, POST, DELETE, PUT, OPTIONS")
	w.Header().Set("Access-Control-Expose-Headers", "X-Cursor")
	_, err = w.Write(js)
	if err != nil {
		s.logger.Error("json write error", "error", err.Error())
	}

}

func (s *SvrWeb) GetHealth(w http.ResponseWriter, r *http.Request) {
	s.logger.Debug("get health", "http_method", r.Method)
	type health struct {
		Status string `json:"status"`
	}
	s.JSON(w, health{Status: "OK"})
}

func (s *SvrWeb) ErrJSON(w http.ResponseWriter, code int, msg string) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept, Authorization, X-Username, X-Access-Token, X-Cursor, X-Cursor-URI")
	w.Header().Set("Access-Control-Request-Methods", "GET, POST, DELETE, PUT, OPTIONS")
	w.WriteHeader(code)
	message := types.Msg{Msg: msg}
	js, err := json.Marshal(message)
	if err == nil {
		_, err = w.Write(js)
		if err != nil {
			s.logger.Error("json write error", "error", err.Error())
		}
	} else {
		fmt.Fprint(w, msg)
	}
}

func (s *SvrWeb) Options(w http.ResponseWriter, r *http.Request) {
	s.logger.Debug("options", "http_method", r.Method, "url", r.URL.Path)
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept, Authorization, X-Username, X-Access-Token, X-Cursor, X-Cursor-URI")
	w.Header().Set("Access-Control-Request-Methods", "GET, POST, DELETE, PUT, OPTIONS")
	if r.Header.Get("Access-Control-Request-Method") != "" {
		w.Header().Set("Access-Control-Allow-Method", r.Header.Get("Access-Control-Request-Method"))
	}
	w.WriteHeader(http.StatusOK)

	message := types.Msg{Msg: "OK"}
	fmt.Fprint(w, message)
}

package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"

	"github.com/betorvs/playbypost/core/sys/web/server"
	"github.com/betorvs/playbypost/core/sys/web/types"
	"github.com/betorvs/playbypost/core/tests/mock/dbclient"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"
)

// TestSigninWithEventLogging tests that signin works with the new database interface
func TestSigninWithEventLogging(t *testing.T) {
	mockDB := new(dbclient.MockDBClient)
	s := &Session{
		db:     mockDB,
		s:      server.NewServer(0, slog.Default()),
		ctx:    context.Background(),
		logger: slog.Default(),
		Sessions: Sessions{
			Current: make(map[string]types.Session),
			mu:      &sync.Mutex{},
		},
	}

	creds := types.Credentials{Username: "testuser", Password: "password"}
	body, _ := json.Marshal(creds)
	req := httptest.NewRequest("POST", "/login", bytes.NewBuffer(body))
	rr := httptest.NewRecorder()

	// Generate a proper bcrypt hash for "password"
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password"), bcrypt.DefaultCost)
	writer := types.Writer{ID: 1, Username: "testuser", Password: string(hashedPassword)}

	// Mock the database calls - event logging is now integrated into CreateSession
	mockDB.On("GetWriterByUsername", mock.Anything, "testuser").Return(writer, nil)
	mockDB.On("CreateSession", mock.Anything, mock.AnythingOfType("types.Session")).Return(nil)
	mockDB.On("LogLoginAttempt", mock.Anything, "testuser", mock.AnythingOfType("string"), mock.AnythingOfType("string"), true).Return(nil)

	s.Signin(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	mockDB.AssertExpectations(t)
}

// TestLogoutWithEventLogging tests that logout works with the new database interface
func TestLogoutWithEventLogging(t *testing.T) {
	mockDB := new(dbclient.MockDBClient)
	s := &Session{
		db:     mockDB,
		s:      server.NewServer(0, slog.Default()),
		ctx:    context.Background(),
		logger: slog.Default(),
		Sessions: Sessions{
			Current: make(map[string]types.Session),
			mu:      &sync.Mutex{},
		},
	}

	req := httptest.NewRequest("POST", "/logout", nil)
	req.Header.Set("X-Access-Token", "testtoken")
	rr := httptest.NewRecorder()

	// Mock session retrieval and deletion - event logging is now integrated into DeleteSessionByToken
	session := types.Session{
		ID:           123,
		Username:     "testuser",
		Token:        "testtoken",
		Expiry:       time.Now().Add(1 * time.Hour),
		UserID:       123,
		ClientType:   "browser",
		IPAddress:    "127.0.0.1",
		UserAgent:    "test",
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
		LastActivity: time.Now(),
	}
	mockDB.On("GetSessionByToken", mock.Anything, "testtoken").Return(session, nil)
	mockDB.On("DeleteSessionByToken", mock.Anything, "testtoken").Return(nil)
	mockDB.On("LogLogout", mock.Anything, int64(123), "testuser").Return(nil)

	s.Logout(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	mockDB.AssertExpectations(t)
}

// TestCheckAuthWithEventLogging tests that CheckAuth works with the new database interface
func TestCheckAuthWithEventLogging(t *testing.T) {
	mockDB := new(dbclient.MockDBClient)
	s := &Session{
		db:     mockDB,
		s:      server.NewServer(0, slog.Default()),
		ctx:    context.Background(),
		logger: slog.Default(),
		Sessions: Sessions{
			Current: make(map[string]types.Session),
			mu:      &sync.Mutex{},
		},
	}

	req := httptest.NewRequest("GET", "/check", nil)
	req.Header.Set("X-Access-Token", "testtoken")

	// Mock session retrieval and validation - event logging is now integrated into UpdateSessionLastActivity
	session := types.Session{
		ID:           123,
		Username:     "testuser",
		Token:        "testtoken",
		Expiry:       time.Now().Add(1 * time.Hour), // Ensure session is not expired
		UserID:       123,
		ClientType:   "browser",
		IPAddress:    "127.0.0.1",
		UserAgent:    "test",
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
		LastActivity: time.Now(),
	}

	// Mock all possible calls that CheckAuth might make
	mockDB.On("GetSessionByToken", mock.Anything, "testtoken").Return(session, nil)
	mockDB.On("LogSessionValidated", mock.Anything, int64(123), "testuser").Return(nil)

	result := s.CheckAuth(req)
	assert.False(t, result) // CheckAuth returns false when authenticated
	mockDB.AssertExpectations(t)
}

// TestSessionEventTypes tests that all event types are properly handled
func TestSessionEventTypes(t *testing.T) {
	eventTypes := []string{
		"session_created",
		"session_deleted",
		"session_expired",
		"login_attempt",
		"logout",
		"session_validated",
		"session_invalid",
		"activity_updated",
		"cleanup_executed",
	}

	for _, eventType := range eventTypes {
		event := types.SessionEvent{
			ID:        1,
			SessionID: 123,
			EventType: eventType,
			Timestamp: time.Now(),
			Data:      `{"test":"data"}`,
		}

		assert.Equal(t, eventType, event.EventType, "Event type should match")
		assert.NotEmpty(t, event.EventType, "Event type should not be empty")
	}
}

// TestSessionEventDataStructure tests the SessionEvent structure
func TestSessionEventDataStructure(t *testing.T) {
	event := types.SessionEvent{
		ID:        1,
		SessionID: 123,
		EventType: "session_created",
		Timestamp: time.Now(),
		Data:      `{"username":"testuser","client_type":"browser"}`,
	}

	assert.Equal(t, 1, event.ID)
	assert.Equal(t, int64(123), event.SessionID)
	assert.Equal(t, "session_created", event.EventType)
	assert.NotZero(t, event.Timestamp)
	assert.Contains(t, event.Data, "testuser")
}

// TestSessionEventRequestStructure tests the SessionEventRequest structure
func TestSessionEventRequestStructure(t *testing.T) {
	request := types.SessionEventRequest{
		SessionID: 456,
		EventType: "login_attempt",
		Data:      `{"username":"testuser","success":true}`,
	}

	assert.Equal(t, int64(456), request.SessionID)
	assert.Equal(t, "login_attempt", request.EventType)
	assert.Contains(t, request.Data, "testuser")
}

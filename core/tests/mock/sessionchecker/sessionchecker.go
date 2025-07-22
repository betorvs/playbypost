package sessionchecker

import (
	"net/http"

	"github.com/betorvs/playbypost/core/sys/web/types"
	"github.com/stretchr/testify/mock"
)

type MockSessionChecker struct {
	mock.Mock
}

func (m *MockSessionChecker) CheckAuth(r *http.Request) bool {
	args := m.Called(r)
	return args.Bool(0)
}

func (m *MockSessionChecker) AddAdminSession(admin, token string) {
	m.Called(admin, token)
}

func (m *MockSessionChecker) Admin() string {
	args := m.Called()
	return args.String(0)
}

func (m *MockSessionChecker) GetActiveSessions() map[string]types.Session {
	args := m.Called()
	return args.Get(0).(map[string]types.Session)
}

func (m *MockSessionChecker) Logout(w http.ResponseWriter, r *http.Request) {
	m.Called(w, r)
}

func (m *MockSessionChecker) ValidateSession(w http.ResponseWriter, r *http.Request) {
	m.Called(w, r)
}

func (m *MockSessionChecker) Signin(w http.ResponseWriter, r *http.Request) {
	m.Called(w, r)
}

func (m *MockSessionChecker) GetSessionEvents(w http.ResponseWriter, r *http.Request) {
	m.Called(w, r)
}

func (m *MockSessionChecker) GetAllSessions(w http.ResponseWriter, r *http.Request) {
	m.Called(w, r)
}

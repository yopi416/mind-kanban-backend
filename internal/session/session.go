package session

import (
	"sync"
	"time"
)

type SessionData struct {
	UserID    int64 // DB上のユーザーID
	ExpiresAt time.Time
}

type SessionManager struct {
	mu   sync.Mutex
	data map[string]SessionData
	ttl  time.Duration
}

func NewSessionManager(ttl time.Duration) *SessionManager {
	return &SessionManager{
		data: make(map[string]SessionData),
		ttl:  ttl,
	}
}

func (sm *SessionManager) CreateSession(sessID string, userID int64) {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	sm.data[sessID] = SessionData{
		UserID:    userID,
		ExpiresAt: time.Now().Add(sm.ttl),
	}
}

func (sm *SessionManager) GetSession(sessID string) (int64, bool) {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	sessData, ok := sm.data[sessID]

	// sessionがない or 有効期限切れの場合はsession無しでreturn
	if !ok || time.Now().After(sessData.ExpiresAt) {
		// 有効期限切れの場合は削除
		delete(sm.data, sessID)
		return 0, false
	}

	return sessData.UserID, true
}

func (sm *SessionManager) DeleteSession(sessID string) {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	delete(sm.data, sessID)
}

func (sm *SessionManager) GetTTL() time.Duration {
	return sm.ttl
}

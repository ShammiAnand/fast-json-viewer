package services

import (
	"errors"
	"log"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/shammianand/fast-json-viewer/internal/models"
)

type SessionManager struct {
	sessions map[string]*sessionData
	mu       sync.RWMutex
}

type sessionData struct {
	trie      *models.JSONTrie
	createdAt time.Time
}

func NewSessionManager() *SessionManager {
	sm := &SessionManager{
		sessions: make(map[string]*sessionData),
	}
	go sm.cleanupRoutine()
	return sm
}

func (sm *SessionManager) CreateSession(trie *models.JSONTrie) string {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	id := uuid.New().String()
	sm.sessions[id] = &sessionData{
		trie:      trie,
		createdAt: time.Now(),
	}
	return id
}

func (sm *SessionManager) GetSession(id string) (*models.JSONTrie, error) {
	sm.mu.RLock()
	defer sm.mu.RUnlock()

	if session, exists := sm.sessions[id]; exists {
		return session.trie, nil
	}
	log.Printf("session not found for session id: %v", id)
	return nil, errors.New("session not found")
}

func (sm *SessionManager) cleanupRoutine() {
	ticker := time.NewTicker(10 * time.Minute)
	for range ticker.C {
		sm.mu.Lock()
		for id, session := range sm.sessions {
			if time.Since(session.createdAt) > 1*time.Hour {
				delete(sm.sessions, id)
			}
		}
		sm.mu.Unlock()
	}
}

package application

import (
	"sync"

	"github.com/mpraes/tabyte/internal/domain"
)

type SessionStore struct {
	mu   sync.RWMutex
	byID map[string]domain.AnalysisSession
}

func NewSessionStore() *SessionStore {
	return &SessionStore{byID: make(map[string]domain.AnalysisSession)}
}

func (s *SessionStore) Save(session domain.AnalysisSession) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.byID[session.ID] = session
}

func (s *SessionStore) Get(id string) (domain.AnalysisSession, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	session, ok := s.byID[id]
	return session, ok
}

func (s *SessionStore) List() []domain.AnalysisSession {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make([]domain.AnalysisSession, 0, len(s.byID))
	for _, session := range s.byID {
		out = append(out, session)
	}
	return out
}

func (s *SessionStore) Delete(id string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.byID[id]; !ok {
		return false
	}
	delete(s.byID, id)
	return true
}
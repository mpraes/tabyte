package application

import (
	"fmt"
	"sync"

	"github.com/mpraes/tabyte/internal/domain"
)

type SessionStore struct {
	mu   sync.RWMutex
	byID map[string]domain.AnalysisSession
	repo SessionRepository
}

func NewSessionStore(repo SessionRepository) *SessionStore {
	return &SessionStore{
		byID: make(map[string]domain.AnalysisSession),
		repo: repo,
	}
}

func (s *SessionStore) Save(session domain.AnalysisSession) {
	s.mu.Lock()
	s.byID[session.ID] = session
	repo := s.repo
	s.mu.Unlock()

	if repo != nil {
		if err := repo.UpsertSession(session); err != nil {
			fmt.Printf("persist upsert session %s: %v\n", session.ID, err)
		}
	}
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
	if _, ok := s.byID[id]; !ok {
		s.mu.Unlock()
		return false
	}
	delete(s.byID, id)
	repo := s.repo
	s.mu.Unlock()

	if repo != nil {
		if err := repo.DeleteSession(id); err != nil {
			fmt.Printf("persist delete session %s: %v\n", id, err)
		}
	}
	return true
}

func (s *SessionStore) LoadFromRepo() error {
	if s.repo == nil {
		return nil
	}
	persisted, err := s.repo.LoadAll()
	if err != nil {
		return err
	}
	for _, p := range persisted {
		session, err := RebuildSession(p)
		if err != nil {
			fmt.Printf("persist hydrate session %s: %v\n", p.ID, err)
			continue
		}
		s.mu.Lock()
		s.byID[session.ID] = session
		s.mu.Unlock()
	}
	return nil
}

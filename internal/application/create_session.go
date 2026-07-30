package application

import (
	"errors"
	"strings"
	"fmt"
	"time"
	"github.com/mpraes/tabyte/internal/domain"
)

var (
	ErrInvalidDDL    = errors.New("ddl_text is required")
	ErrInvalidEngine = errors.New("engine must be sqlserver or postgres")
)

type CreateSessionInput struct {
	Engine     string
	SourceName string
	DDLText    string
}

func CreateSession(store *SessionStore, in CreateSessionInput) (domain.AnalysisSession, error) {
	ddl := strings.TrimSpace(in.DDLText)
	if ddl == "" {
		return domain.AnalysisSession{}, ErrInvalidDDL
	}

	engine, ok := domain.ParseEngine(strings.ToLower(strings.TrimSpace(in.Engine)))
	if !ok {
		return domain.AnalysisSession{}, ErrInvalidEngine
	}

	session := domain.AnalysisSession{
		ID:         fmt.Sprintf("as_%d", time.Now().UnixNano()),
		Engine:     engine,
		SourceName: in.SourceName,
		DDLText:    ddl,
		Status:     "created",
	}
	store.Save(session)
	return session, nil
}

func GetSession(store *SessionStore, id string) (domain.AnalysisSession, bool) {
	return store.Get(id)
}

func ListSessions(store *SessionStore) []domain.AnalysisSession {
	return store.List()
}

func DeleteSession(store *SessionStore, id string) bool {
	return store.Delete(id)
}
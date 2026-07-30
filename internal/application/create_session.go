package application

import (
	"errors"
	"strings"
	"fmt"
	"time"

	"github.com/mpraes/tabyte/internal/domain"
	"github.com/mpraes/tabyte/internal/parser"
)

var (
	ErrInvalidDDL    = errors.New("ddl_text is required")
	ErrInvalidEngine = errors.New("engine must be sqlserver or postgres")
	ErrNoTablesFound = errors.New("no tables found in DDL")
	ErrSessionNotFound = errors.New("session not found")
	ErrTableNotFound   = errors.New("table not found")
	ErrInvalidRowCount = errors.New("assumed_row_count must be > 0")
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

	eng, ok := domain.ParseEngine(strings.ToLower(strings.TrimSpace(in.Engine)))
	if !ok {
		return domain.AnalysisSession{}, ErrInvalidEngine
	}

	tables := parser.ParseTables(ddl)
	tables = enrichTables(eng, tables)
	if len(tables) == 0 {
		return domain.AnalysisSession{}, ErrNoTablesFound
	}
	total := sumSchemaBytes(tables) 
	warnings := collectWarnings(tables)

	session := domain.AnalysisSession{
		ID:                  fmt.Sprintf("as_%d", time.Now().UnixNano()),
		Engine:              eng,
		SourceName:          in.SourceName,
		DDLText:             ddl,
		Status:              "created",
		Tables:              tables,
		EstimatedTotalBytes: &total,
		Warnings:            warnings,
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

func UpdateTableRowCount(store *SessionStore, sessionID, tableName string, rows int64) (domain.AnalysisSession, error) {
	if rows <= 0 {
		return domain.AnalysisSession{}, ErrInvalidRowCount
	}

	session, ok := store.Get(sessionID)
	if !ok {
		return domain.AnalysisSession{}, ErrSessionNotFound
	}

	found := false
	for i, t := range session.Tables {
		if t.Name == tableName {
			session.Tables[i].AssumedRowCount = rows
			session.Tables[i] = estimateTableVolume(session.Tables[i])
			found = true
			break
		}
	}
	if !found {
		return domain.AnalysisSession{}, ErrTableNotFound
	}

	total := sumSchemaBytes(session.Tables)
	session.EstimatedTotalBytes = &total

	store.Save(session)
	return session, nil
}
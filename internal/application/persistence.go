package application

import "github.com/mpraes/tabyte/internal/domain"

type Setting struct {
	Key       string
	Value     string
	ValueType string
}

type SettingsRepository interface {
	ListSettings() ([]Setting, error)
	UpsertSetting(key, value, valueType string) error
}

type PersistedTable struct {
	Name                string
	AssumedRowCount     int64
	GrowthRowsPerPeriod int64
	GrowthPeriod        string
	GrowthHorizon       int64
}

type PersistedSession struct {
	ID         string
	Engine     string
	SourceName string
	DDLText    string
	Status     string
	Tables     []PersistedTable
}

type SessionRepository interface {
	UpsertSession(session domain.AnalysisSession) error
	DeleteSession(id string) error
	LoadAll() ([]PersistedSession, error)
}

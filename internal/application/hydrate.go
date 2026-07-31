package application

import (
	"strings"

	"github.com/mpraes/tabyte/internal/domain"
	"github.com/mpraes/tabyte/internal/parser"
)

func RebuildSession(p PersistedSession) (domain.AnalysisSession, error) {
	ddl := strings.TrimSpace(p.DDLText)
	if ddl == "" {
		return domain.AnalysisSession{}, ErrInvalidDDL
	}
	eng, ok := domain.ParseEngine(strings.ToLower(strings.TrimSpace(p.Engine)))
	if !ok {
		return domain.AnalysisSession{}, ErrInvalidEngine
	}

	tables := parser.ParseTables(ddl)
	tables = enrichTables(eng, tables)
	if len(tables) == 0 {
		return domain.AnalysisSession{}, ErrNoTablesFound
	}

	overrides := map[string]PersistedTable{}
	for _, t := range p.Tables {
		overrides[strings.ToLower(t.Name)] = t
	}

	for i, t := range tables {
		o, ok := overrides[strings.ToLower(t.Name)]
		if !ok {
			continue
		}
		if o.AssumedRowCount > 0 {
			t.AssumedRowCount = o.AssumedRowCount
			t = estimateTableVolume(t)
			t = estimateIndexes(eng, t)
		}
		if o.GrowthRowsPerPeriod > 0 && o.GrowthPeriod != "" && o.GrowthHorizon > 0 {
			if grown, err := applyTableGrowth(t, o.GrowthRowsPerPeriod, o.GrowthPeriod, o.GrowthHorizon); err == nil {
				t = grown
			}
		}
		tables[i] = t
	}

	total := sumSchemaBytes(tables)
	session := domain.AnalysisSession{
		ID:                  p.ID,
		Engine:              eng,
		SourceName:          p.SourceName,
		DDLText:             ddl,
		Status:              p.Status,
		Tables:              tables,
		EstimatedTotalBytes: &total,
		Warnings:            collectWarnings(tables),
		Signals:             collectSignals(tables),
	}
	session.ProjectedTotalBytes = sumProjectedBytes(tables)
	return session, nil
}

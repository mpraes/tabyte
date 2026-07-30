package application

import (
	"errors"
	"strings"

	"github.com/mpraes/tabyte/internal/domain"
)

var (
	ErrInvalidGrowth = errors.New("rows_per_period and horizon must be > 0; period must be hour, day, or month")
)

func applyTableGrowth(table domain.Table, rowsPerPeriod int64, period string, horizon int64) (domain.Table, error) {
	period = strings.ToLower(strings.TrimSpace(period))
	switch period {
	case "hour", "day", "month":
	default:
		return table, ErrInvalidGrowth
	}
	if rowsPerPeriod <= 0 || horizon <= 0 {
		return table, ErrInvalidGrowth
	}
	if table.EstimatedRowBytes == nil {
		return table, ErrInvalidGrowth
	}

	projectedRows := table.AssumedRowCount + rowsPerPeriod*horizon
	projectedBytes := *table.EstimatedRowBytes * projectedRows

	table.GrowthRowsPerPeriod = rowsPerPeriod
	table.GrowthPeriod = period
	table.GrowthHorizon = horizon
	table.ProjectedRowCount = &projectedRows
	table.ProjectedTableBytes = &projectedBytes
	return table, nil
}

func sumProjectedBytes(tables []domain.Table) *int64 {
	var sum int64
	any := false
	for _, t := range tables {
		if t.ProjectedTableBytes != nil {
			sum += *t.ProjectedTableBytes
			any = true
		}
	}
	if !any {
		return nil
	}
	return &sum
}

func UpdateTableGrowth(store *SessionStore, sessionID, tableName string, rowsPerPeriod int64, period string, horizon int64) (domain.AnalysisSession, error) {
	session, ok := store.Get(sessionID)
	if !ok {
		return domain.AnalysisSession{}, ErrSessionNotFound
	}

	found := false
	for i, t := range session.Tables {
		if t.Name == tableName {
			updated, err := applyTableGrowth(t, rowsPerPeriod, period, horizon)
			if err != nil {
				return domain.AnalysisSession{}, err
			}
			session.Tables[i] = updated
			found = true
			break
		}
	}
	if !found {
		return domain.AnalysisSession{}, ErrTableNotFound
	}

	session.ProjectedTotalBytes = sumProjectedBytes(session.Tables)
	store.Save(session)
	return session, nil
}
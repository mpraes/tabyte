package application

import (
	"github.com/mpraes/tabyte/internal/domain"
	"github.com/mpraes/tabyte/internal/engine/postgres"
	"github.com/mpraes/tabyte/internal/engine/sqlserver"
)

func enrichTables(eng domain.Engine, tables []domain.Table) []domain.Table {
	out := make([]domain.Table, len(tables))
	for i, t := range tables {
		cols := make([]domain.Column, len(t.Columns))
		for j, c := range t.Columns {
			switch eng {
			case domain.EnginePostgres:
				cols[j] = postgres.NormalizeColumn(c)
				cols[j] = postgres.EstimateColumn(cols[j])
			case domain.EngineSQLServer:
				cols[j] = sqlserver.NormalizeColumn(c)
				cols[j] = sqlserver.EstimateColumn(cols[j])
			default:
				c.NormalizedType = "unknown"
				cols[j] = c
			}
		}
		t.Columns = cols

		switch eng {
		case domain.EnginePostgres:
			t = postgres.EstimateRow(t)
		case domain.EngineSQLServer:
			t = sqlserver.EstimateRow(t)
		}
		t = estimateTableVolume(t)
		out[i] = t
	}
	return out
}
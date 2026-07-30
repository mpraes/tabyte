package application

import (
	"strings"

	"github.com/mpraes/tabyte/internal/domain"
)

func estimateIndexes(eng domain.Engine, table domain.Table) domain.Table {
	colsByName := map[string]domain.Column{}
	for _, c := range table.Columns {
		colsByName[strings.ToLower(c.Name)] = c
	}

	var indexTotal int64
	for i, idx := range table.Indexes {
		// SQL Server clustered PK: no extra heap/index volume in v0
		if eng == domain.EngineSQLServer && idx.Kind == "primary_key" {
			zero := int64(0)
			table.Indexes[i].EstimatedBytes = &zero
			continue
		}

		var payload int64
		for _, name := range idx.Columns {
			if c, ok := colsByName[strings.ToLower(name)]; ok && c.EstimatedBytes != nil {
				payload += *c.EstimatedBytes
			}
		}
		total := payload * table.AssumedRowCount
		table.Indexes[i].EstimatedBytes = &total
		indexTotal += total
	}

	if table.Calculation != nil {
		table.Calculation.IndexBytes = indexTotal
	}
	return table
}
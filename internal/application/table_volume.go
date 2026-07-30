package application

import "github.com/mpraes/tabyte/internal/domain"

const DefaultAssumedRowCount int64 = 1000

func estimateTableVolume(table domain.Table) domain.Table {
	if table.AssumedRowCount <= 0 {
		table.AssumedRowCount = DefaultAssumedRowCount
	}
	if table.EstimatedRowBytes == nil {
		return table
	}
	total := *table.EstimatedRowBytes * table.AssumedRowCount
	table.EstimatedTableBytes = &total
	return table
}
package application

import (
	"fmt"

	"github.com/mpraes/tabyte/internal/domain"
)

const (
	wideRowBytesThreshold = int64(100)
	manyColumnsThreshold  = 20
	variableHeavyThreshold = 5
)

func collectSignals(tables []domain.Table) []domain.Signal {
	var out []domain.Signal
	for _, t := range tables {
		if t.EstimatedRowBytes != nil && *t.EstimatedRowBytes >= wideRowBytesThreshold {
			out = append(out, domain.Signal{
				Code:    "WIDE_ROW",
				Message: fmt.Sprintf("estimated row is %d bytes; wide rows can increase I/O and memory pressure", *t.EstimatedRowBytes),
				Table:   t.Name,
			})
		}
		if len(t.Columns) >= manyColumnsThreshold {
			out = append(out, domain.Signal{
				Code:    "MANY_COLUMNS",
				Message: fmt.Sprintf("table has %d columns; wide schemas can slow scans and maintenance", len(t.Columns)),
				Table:   t.Name,
			})
		}
		varCount := 0
		for _, c := range t.Columns {
			switch c.NormalizedType {
			case "varchar", "nvarchar", "text", "ntext":
				varCount++
			}
		}
		if varCount >= variableHeavyThreshold {
			out = append(out, domain.Signal{
				Code:    "VARIABLE_HEAVY",
				Message: fmt.Sprintf("table has %d variable-length columns; updates and sorting may be more expensive", varCount),
				Table:   t.Name,
			})
		}
	}
	if out == nil {
		return []domain.Signal{}
	}
	return out
}
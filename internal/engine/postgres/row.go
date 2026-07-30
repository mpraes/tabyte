package postgres

import "github.com/mpraes/tabyte/internal/domain"

// Heap tuple header is commonly modeled as 23 bytes; null bitmap is ceil(n/8).
// Alignment/TOAST ignored in v0.
func EstimateRow(table domain.Table) domain.Table {
	var payload int64
	for _, c := range table.Columns {
		if c.EstimatedBytes != nil {
			payload += *c.EstimatedBytes
		}
	}
	nullBitmap := int64((len(table.Columns) + 7) / 8)
	total := int64(23) + nullBitmap + payload
	table.EstimatedRowBytes = &total
	return table
}
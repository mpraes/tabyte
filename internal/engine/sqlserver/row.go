package sqlserver

import "github.com/mpraes/tabyte/internal/domain"

// v0: 4-byte row header + null bitmap + column payload.
// Bit packing / page density / variable-column offsets deferred.
func EstimateRow(table domain.Table) domain.Table {
	var payload int64
	for _, c := range table.Columns {
		if c.EstimatedBytes != nil {
			payload += *c.EstimatedBytes
		}
	}
	nullBitmap := int64((len(table.Columns) + 7) / 8)
	total := int64(4) + nullBitmap + payload
	table.EstimatedRowBytes = &total
	return table
}
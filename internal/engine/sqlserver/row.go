package sqlserver

import (
	"github.com/mpraes/tabyte/internal/domain"
)

// v0: 4-byte row header + null bitmap + column payload (23 bytes).
// Bit packing / page density / variable-column offsets deferred.
func EstimateRow(table domain.Table) domain.Table {
	var payload int64
	for _, c := range table.Columns {
		if c.EstimatedBytes != nil {
			payload += *c.EstimatedBytes
		}
	}
	nullBitmap := int64((len(table.Columns) + 7) / 8)
	header := int64(4)
	total := header + nullBitmap + payload

	table.EstimatedRowBytes = &total
	table.Calculation = &domain.RowCalculation{
		ColumnPayloadBytes: payload,
		RowHeaderBytes:     header,
		NullBitmapBytes:    nullBitmap,
		EstimatedRowBytes:  total,
		IndexBytes:         0,
	}
	return table
}
package postgres

import (
	"testing"

	"github.com/mpraes/tabyte/internal/domain"
)

func TestEstimateRow(t *testing.T) {
	b4 := int64(4)
	table := domain.Table{
		Name: "a",
		Columns: []domain.Column{
			{Name: "id", EstimatedBytes: &b4},
		},
	}
	table = EstimateRow(table)
	// 23 + ceil(1/8)=1 + 4 = 28
	if table.EstimatedRowBytes == nil || *table.EstimatedRowBytes != 28 {
		t.Fatalf("row: %v", table.EstimatedRowBytes)
	}
}

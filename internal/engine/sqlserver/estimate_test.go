package sqlserver

import (
	"testing"

	"github.com/mpraes/tabyte/internal/domain"
)

func TestEstimateDateTime2(t *testing.T) {
	col := NormalizeColumn(domain.Column{OriginalType: "datetime2(3)"})
	col = EstimateColumn(col)

	if col.NormalizedType != "datetime2" {
		t.Fatalf("normalized: %q", col.NormalizedType)
	}
	if col.EstimatedBytes == nil || *col.EstimatedBytes != 7 {
		t.Fatalf("bytes: %v", col.EstimatedBytes)
	}
}

func TestEstimateNVarchar(t *testing.T) {
	col := NormalizeColumn(domain.Column{OriginalType: "NVARCHAR(50)"})
	col = EstimateColumn(col)
	if col.NormalizedType != "nvarchar" {
		t.Fatalf("type: %q", col.NormalizedType)
	}
	// Length 50 < fallback 64 → assumed = 50 → 50*2+2 = 102
	if col.EstimatedBytes == nil || *col.EstimatedBytes != 102 {
		t.Fatalf("bytes: %v", *col.EstimatedBytes)
	}
}
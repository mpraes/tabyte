package postgres

import (
	"testing"

	"github.com/mpraes/tabyte/internal/domain"
)

func TestNormalizeColumnVarchar(t *testing.T) {
	col := NormalizeColumn(domain.Column{OriginalType: "VARCHAR(100)"})
	if col.NormalizedType != "varchar" {
		t.Fatalf("type: %q", col.NormalizedType)
	}
	if col.Length == nil || *col.Length != 100 {
		t.Fatalf("length: %v", col.Length)
	}
}
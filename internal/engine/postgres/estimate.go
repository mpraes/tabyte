package postgres

import (
	"github.com/mpraes/tabyte/internal/domain"
)

func EstimateColumn(col domain.Column) domain.Column {
	var n int64

	switch col.NormalizedType {
	case "smallint":
		n = 2
	case "int":
		n = 4
	case "bigint":
		n = 8
	case "boolean":
		n = 1
	case "uuid":
		n = 16
	case "date":
		n = 4
	case "timestamp":
		n = 8
	case "float":
		n = 8
	case "char":
		if col.Length != nil {
			n = pgVarlenaSize(int64(*col.Length), true)
		}
	case "varchar":
		avg := assumedLength(col, 64)
		col.AssumedAvgLength = &avg
		n = pgVarlenaSize(int64(avg), false)
	case "text":
		n = pgVarlenaSize(256, false)
	case "numeric":
		n = estimatePostgresNumeric(col.Precision, col.Scale)
	default:
		return col
	}

	col.EstimatedBytes = &n
	return col
}

func pgVarlenaSize(dataBytes int64, blankPadded bool) int64 {
	if dataBytes < 127 {
		return dataBytes + 1
	}
	return dataBytes + 4
}

func assumedLength(col domain.Column, fallback int) int {
	if col.AssumedAvgLength != nil && *col.AssumedAvgLength > 0 {
		return *col.AssumedAvgLength
	}
	if col.Length != nil && *col.Length > 0 {
		if *col.Length < fallback {
			return *col.Length
		}
		return *col.Length / 2
	}
	return fallback
}

func estimatePostgresNumeric(precision, scale *int) int64 {
	p := 18
	if precision != nil && *precision > 0 {
		p = *precision
	}
	digitsGroups := (p + 3) / 4
	return int64(digitsGroups*2 + 4)
}
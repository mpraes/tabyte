package sqlserver

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
	case "tinyint":
		n = 1
	case "boolean":
		n = 1
	case "uuid":
		n = 16
	case "date":
		n = 3
	case "datetime":
		n = 8
	case "smalldatetime":
		n = 4
	case "datetime2":
		n = estimateSQLServerDateTime2(col.Scale)
	case "timestamp":
		n = 8
	case "float":
		n = 8
	case "char":
		if col.Length != nil {
			n = int64(*col.Length)
		}
	case "nchar":
		if col.Length != nil {
			n = int64(*col.Length) * 2
		}
	case "varchar":
		avg := assumedLength(col, 64)
		col.AssumedAvgLength = &avg
		n = int64(avg) + 2
	case "nvarchar":
		avg := assumedLength(col, 64)
		col.AssumedAvgLength = &avg
		n = int64(avg)*2 + 2
	case "text":
		n = 256 + 2
	case "numeric":
		n = estimateSQLServerNumeric(col.Precision)
	default:
		return col
	}

	col.EstimatedBytes = &n
	return col
}

func estimateSQLServerNumeric(precision *int) int64 {
	p := 18
	if precision != nil && *precision > 0 {
		p = *precision
	}
	switch {
	case p >= 1 && p <= 9:
		return 5
	case p >= 10 && p <= 19:
		return 9
	case p >= 20 && p <= 28:
		return 13
	default:
		return 17
	}
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

func estimateSQLServerDateTime2(scale *int) int64 {
	s := 7
	if scale != nil && *scale >= 0 {
		s = *scale
	}
	switch {
	case s <= 2:
		return 6
	case s <= 4:
		return 7
	default:
		return 8
	}
}
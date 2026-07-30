package postgres

import (
	"github.com/mpraes/tabyte/internal/domain"
	"github.com/mpraes/tabyte/internal/engine"
)

func NormalizeColumn(col domain.Column) domain.Column {
	info := engine.ParseTypeInfo(col.OriginalType)
	switch info.Base {
	case "int", "integer", "int4":
		col.NormalizedType = "int"
	case "bigint", "int8":
		col.NormalizedType = "bigint"
	case "smallint", "int2":
		col.NormalizedType = "smallint"
	case "varchar", "character varying":
		col.NormalizedType = "varchar"
		col.Length = info.Length
	case "char", "character":
		col.NormalizedType = "char"
		col.Length = info.Length
	case "text":
		col.NormalizedType = "text"
	case "numeric", "decimal":
		col.NormalizedType = "numeric"
		col.Precision = info.Precision
		col.Scale = info.Scale
	case "bool", "boolean":
		col.NormalizedType = "boolean"
	case "uuid":
		col.NormalizedType = "uuid"
	case "timestamp", "timestamptz", "timestamp without time zone", "timestamp with time zone":
		col.NormalizedType = "timestamp"
	case "date":
		col.NormalizedType = "date"
	default:
		col.NormalizedType = "unknown"
	}
	return col
}
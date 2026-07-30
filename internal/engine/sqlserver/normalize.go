package sqlserver

import (
	"github.com/mpraes/tabyte/internal/domain"
	"github.com/mpraes/tabyte/internal/engine"
)

func NormalizeColumn(col domain.Column) domain.Column {
	info := engine.ParseTypeInfo(col.OriginalType)
	switch info.Base {
	case "int", "integer":
		col.NormalizedType = "int"
	case "bigint":
		col.NormalizedType = "bigint"
	case "smallint":
		col.NormalizedType = "smallint"
	case "tinyint":
		col.NormalizedType = "tinyint"
	case "varchar", "nvarchar":
		col.NormalizedType = "varchar"
		col.Length = info.Length
	case "char", "nchar":
		col.NormalizedType = "char"
		col.Length = info.Length
	case "text", "ntext":
		col.NormalizedType = "text"
	case "numeric", "decimal":
		col.NormalizedType = "numeric"
		col.Precision = info.Precision
		col.Scale = info.Scale
	case "bit":
		col.NormalizedType = "boolean"
	case "uniqueidentifier":
		col.NormalizedType = "uuid"
	case "datetime", "datetime2", "smalldatetime":
		col.NormalizedType = "timestamp"
	case "date":
		col.NormalizedType = "date"
	case "float", "real":
		col.NormalizedType = "float"
	default:
		col.NormalizedType = "unknown"
	}
	return col
}
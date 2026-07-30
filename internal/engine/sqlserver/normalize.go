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
	case "varchar":
		col.NormalizedType = "varchar"
		col.Length = info.Length
	case "nvarchar":
		col.NormalizedType = "nvarchar"
		col.Length = info.Length
	case "char":
		col.NormalizedType = "char"
		col.Length = info.Length
	case "nchar":
		col.NormalizedType = "nchar"
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
	case "datetime":
		col.NormalizedType = "datetime"
	case "datetime2":
		col.NormalizedType = "datetime2"
		if info.Length != nil {
			col.Scale = info.Length
		} else if info.Precision != nil {
			col.Scale = info.Precision
		}
	case "smalldatetime":
		col.NormalizedType = "smalldatetime"
	case "date":
		col.NormalizedType = "date"
	case "float", "real":
		col.NormalizedType = "float"
	default:
		col.NormalizedType = "unknown"
	}
	return col
}
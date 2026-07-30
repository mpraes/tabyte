package application

import (
	"fmt"

	"github.com/mpraes/tabyte/internal/domain"
)

const wideLengthThreshold = 255

func collectWarnings(tables []domain.Table) []domain.Warning {
	var out []domain.Warning
	for _, t := range tables {
		for _, c := range t.Columns {
			nt := c.NormalizedType
			switch nt {
			case "varchar", "nvarchar":
				if c.Length != nil && *c.Length >= wideLengthThreshold {
					out = append(out, domain.Warning{
						Code:    "WIDE_VARCHAR",
						Message: fmt.Sprintf("%s(%d) is wide; storage and index cost may grow quickly", nt, *c.Length),
						Table:   t.Name,
						Column:  c.Name,
					})
				}
				if c.Length == nil {
					out = append(out, domain.Warning{
						Code:    "GENERIC_TYPE",
						Message: fmt.Sprintf("%s without explicit length is treated as unbounded/generic", nt),
						Table:   t.Name,
						Column:  c.Name,
					})
				}
			case "text", "ntext":
				out = append(out, domain.Warning{
					Code:    "GENERIC_TYPE",
					Message: fmt.Sprintf("%s is a large/generic text type", nt),
					Table:   t.Name,
					Column:  c.Name,
				})
			}
		}
	}
	if out == nil {
		return []domain.Warning{}
	}
	return out
}
package engine

import (
	"regexp"
	"strconv"
	"strings"
)

type TypeInfo struct {
	Base      string // lowercased, e.g. "varchar", "character varying"
	Length    *int
	Precision *int
	Scale     *int
}

var typeRe = regexp.MustCompile(`(?i)^([a-zA-Z][a-zA-Z0-9 ]*)\s*(?:\(([^)]*)\))?`)

func ParseTypeInfo(original string) TypeInfo {
	s := strings.TrimSpace(original)
	// cut common trailing constraints for v0
	upper := strings.ToUpper(s)
	for _, stop := range []string{" NOT NULL", " NULL", " PRIMARY", " UNIQUE", " DEFAULT", " COLLATE", " IDENTITY", " CONSTRAINT"} {
		if i := strings.Index(upper, stop); i > 0 {
			s = strings.TrimSpace(s[:i])
			upper = strings.ToUpper(s)
		}
	}

	m := typeRe.FindStringSubmatch(s)
	if m == nil {
		return TypeInfo{Base: strings.ToLower(strings.TrimSpace(s))}
	}
	base := strings.ToLower(strings.Join(strings.Fields(m[1]), " ")) // "character varying"
	info := TypeInfo{Base: base}
	if m[2] == "" {
		return info
	}
	parts := strings.Split(m[2], ",")
	if len(parts) == 1 {
		if n, err := strconv.Atoi(strings.TrimSpace(parts[0])); err == nil {
			info.Length = &n
			info.Precision = &n // numeric(p) sometimes; engines decide
		}
	}
	if len(parts) >= 2 {
		if p, err := strconv.Atoi(strings.TrimSpace(parts[0])); err == nil {
			info.Precision = &p
		}
		if sc, err := strconv.Atoi(strings.TrimSpace(parts[1])); err == nil {
			info.Scale = &sc
		}
	}
	return info
}
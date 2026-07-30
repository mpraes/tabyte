package parser

import (
	"regexp"
	"strings"

	"github.com/mpraes/tabyte/internal/domain"
)

var createIndexRe = regexp.MustCompile(
	`(?i)CREATE\s+(UNIQUE\s+)?INDEX\s+(?:IF\s+NOT\s+EXISTS\s+)?\[?["']?([\w]+)["']?\]?\s+ON\s+(?:\[?[\w]+\]?\.)?\[?["']?([\w]+)["']?\]?\s*\(([^)]*)\)`,
)

func parseIndexesFromTableBody(tableName, body string) []domain.Index {
	var out []domain.Index
	parts := splitTopLevel(body, ',')
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}
		u := strings.ToUpper(part)

		// column-level: name type ... PRIMARY KEY
		if !isConstraint(part) {
			name, typ := splitNameAndType(part)
			if name != "" && strings.Contains(strings.ToUpper(typ), "PRIMARY KEY") {
				out = append(out, domain.Index{
					Name:    "",
					Table:   tableName,
					Columns: []string{name},
					Kind:    "primary_key",
				})
			}
			continue
		}

		// table-level PRIMARY KEY / CONSTRAINT ... PRIMARY KEY
		if strings.Contains(u, "PRIMARY KEY") {
			cols := parseIdentListInParens(part)
			if len(cols) == 0 {
				continue
			}
			out = append(out, domain.Index{
				Name:    constraintName(part),
				Table:   tableName,
				Columns: cols,
				Kind:    "primary_key",
			})
			continue
		}

		// table-level UNIQUE (cols) — skip UNIQUE as column attribute alone
		if strings.HasPrefix(u, "UNIQUE") || (strings.HasPrefix(u, "CONSTRAINT") && strings.Contains(u, "UNIQUE")) {
			cols := parseIdentListInParens(part)
			if len(cols) == 0 {
				continue
			}
			out = append(out, domain.Index{
				Name:    constraintName(part),
				Table:   tableName,
				Columns: cols,
				Kind:    "unique",
			})
		}
	}
	return out
}

func ParseIndexes(ddl string) []domain.Index {
	matches := createIndexRe.FindAllStringSubmatch(ddl, -1)
	out := make([]domain.Index, 0, len(matches))
	for _, m := range matches {
		unique := strings.TrimSpace(m[1]) != ""
		name := trimIdent(m[2])
		table := trimIdent(m[3])
		cols := splitIdentList(m[4])
		kind := "index"
		if unique {
			kind = "unique"
		}
		out = append(out, domain.Index{
			Name:    name,
			Table:   table,
			Columns: cols,
			Kind:    kind,
		})
	}
	return out
}

func parseIdentListInParens(s string) []string {
	start := strings.Index(s, "(")
	end := strings.LastIndex(s, ")")
	if start < 0 || end <= start {
		return nil
	}
	return splitIdentList(s[start+1 : end])
}

func splitIdentList(s string) []string {
	parts := splitTopLevel(s, ',')
	var out []string
	for _, p := range parts {
		id := trimIdent(p)
		// drop ASC/DESC
		fields := strings.Fields(id)
		if len(fields) == 0 {
			continue
		}
		out = append(out, trimIdent(fields[0]))
	}
	return out
}

func constraintName(part string) string {
	u := strings.ToUpper(strings.TrimSpace(part))
	if !strings.HasPrefix(u, "CONSTRAINT") {
		return ""
	}
	rest := strings.TrimSpace(part[len("CONSTRAINT"):])
	name, _ := splitNameAndType(rest)
	return name
}
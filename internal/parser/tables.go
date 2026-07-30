package parser

import (
	"regexp"
	"strings"
	"unicode"

	"github.com/mpraes/tabyte/internal/domain"
)

var createTableRe = regexp.MustCompile(
	`(?i)CREATE\s+TABLE\s+(?:IF\s+NOT\s+EXISTS\s+)?(?:\[?[\w]+\]?\.)?\[?["']?([\w]+)["']?\]?`,
)

func ParseTables(ddl string) []domain.Table {
	matches := createTableRe.FindAllStringSubmatchIndex(ddl, -1)
	out := make([]domain.Table, 0, len(matches))
	seen := map[string]struct{}{}

	for _, loc := range matches {
		// loc[2]:loc[3] = table name capture
		name := strings.TrimSpace(ddl[loc[2]:loc[3]])
		if name == "" {
			continue
		}
		key := strings.ToLower(name)
		if _, ok := seen[key]; ok {
			continue
		}
		seen[key] = struct{}{}

		bodyStart := findOpenParen(ddl, loc[1]) // after full CREATE TABLE match
		var cols []domain.Column
		var indexes []domain.Index
		if bodyStart >= 0 {
			if body, ok := extractParenBody(ddl, bodyStart); ok {
				cols = parseColumns(body)
				indexes = parseIndexesFromTableBody(name, body)
			}
		}

		out = append(out, domain.Table{
			Name:    name,
			Columns: cols,
			Indexes: indexes,
		})
	}
	standalone := ParseIndexes(ddl)
	if len(standalone) > 0 {
		byTable := map[string]int{}
		for i, t := range out {
			byTable[strings.ToLower(t.Name)] = i
		}
		for _, idx := range standalone {
			if i, ok := byTable[strings.ToLower(idx.Table)]; ok {
				out[i].Indexes = append(out[i].Indexes, idx)
			}
		}
	}
	return out
}

func findOpenParen(s string, from int) int {
	for i := from; i < len(s); i++ {
		if s[i] == '(' {
			return i
		}
		// stop if another statement starts before '('
		if s[i] == ';' {
			return -1
		}
	}
	return -1
}

func extractParenBody(s string, openIdx int) (string, bool) {
	if openIdx < 0 || openIdx >= len(s) || s[openIdx] != '(' {
		return "", false
	}
	depth := 0
	for i := openIdx; i < len(s); i++ {
		switch s[i] {
		case '(':
			depth++
		case ')':
			depth--
			if depth == 0 {
				return s[openIdx+1 : i], true
			}
		}
	}
	return "", false
}

func parseColumns(body string) []domain.Column {
	parts := splitTopLevel(body, ',')
	var cols []domain.Column
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part == "" || isConstraint(part) {
			continue
		}
		name, typ := splitNameAndType(part)
		if name == "" {
			continue
		}
		cols = append(cols, domain.Column{Name: name, OriginalType: typ})
	}
	return cols
}

func splitTopLevel(s string, sep rune) []string {
	var parts []string
	depth := 0
	start := 0
	for i, r := range s {
		switch r {
		case '(':
			depth++
		case ')':
			if depth > 0 {
				depth--
			}
		default:
			if r == sep && depth == 0 {
				parts = append(parts, s[start:i])
				start = i + 1
			}
		}
	}
	parts = append(parts, s[start:])
	return parts
}

func splitNameAndType(part string) (name, typ string) {
	part = strings.TrimSpace(part)
	if part == "" {
		return "", ""
	}
	// first token = name
	i := 0
	for i < len(part) && !unicode.IsSpace(rune(part[i])) {
		i++
	}
	name = trimIdent(part[:i])
	typ = strings.TrimSpace(part[i:])
	return name, typ
}

func trimIdent(s string) string {
	s = strings.TrimSpace(s)
	s = strings.Trim(s, `"'[]`)
	return s
}

func isConstraint(s string) bool {
	u := strings.ToUpper(strings.TrimSpace(s))
	return strings.HasPrefix(u, "PRIMARY KEY") ||
		strings.HasPrefix(u, "UNIQUE") ||
		strings.HasPrefix(u, "CONSTRAINT") ||
		strings.HasPrefix(u, "FOREIGN KEY") ||
		strings.HasPrefix(u, "CHECK") ||
		strings.HasPrefix(u, "INDEX") ||
		strings.HasPrefix(u, "KEY ")
}
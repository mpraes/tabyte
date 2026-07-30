package parser

import (
	"regexp"
	"strings"

	"github.com/mpraes/tabyte/internal/domain"
)

// Matches: CREATE TABLE [IF NOT EXISTS] name / "name" / [name] / schema.name
var createTableRe = regexp.MustCompile(
	`(?i)CREATE\s+TABLE\s+(?:IF\s+NOT\s+EXISTS\s+)?(?:\[?[\w]+\]?\.)?\[?["']?([\w]+)["']?\]?`,
)

func ParseTableNames(ddl string) []domain.Table {
	matches := createTableRe.FindAllStringSubmatch(ddl, -1)
	out := make([]domain.Table, 0, len(matches))
	seen := map[string]struct{}{}
	for _, m := range matches {
		name := strings.TrimSpace(m[1])
		if name == "" {
			continue
		}
		key := strings.ToLower(name)
		if _, ok := seen[key]; ok {
			continue
		}
		seen[key] = struct{}{}
		out = append(out, domain.Table{Name: name})
	}
	return out
}
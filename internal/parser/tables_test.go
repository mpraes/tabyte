package parser

import (
	"strings"
	"testing"
)

func TestParseTableNames(t *testing.T) {
	ddl := `
CREATE TABLE users (id INT);
CREATE TABLE IF NOT EXISTS orders (id INT);
create table "products" (id int);
`
	got := ParseTables(ddl)
	if len(got) != 3 {
		t.Fatalf("want 3 tables, got %d: %+v", len(got), got)
	}
}

func TestParseTablesColumns(t *testing.T) {
	ddl := `CREATE TABLE users (
  id INT,
  name VARCHAR(100),
  PRIMARY KEY (id)
);`
	got := ParseTables(ddl)
	if len(got) != 1 {
		t.Fatalf("tables: %d", len(got))
	}
	if len(got[0].Columns) != 2 {
		t.Fatalf("columns: %+v", got[0].Columns)
	}
	if got[0].Columns[0].Name != "id" || !strings.Contains(strings.ToUpper(got[0].Columns[0].OriginalType), "INT") {
		t.Fatalf("col0: %+v", got[0].Columns[0])
	}
	if got[0].Columns[1].Name != "name" || !strings.Contains(strings.ToUpper(got[0].Columns[1].OriginalType), "VARCHAR") {
		t.Fatalf("col1: %+v", got[0].Columns[1])
	}
}
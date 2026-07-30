package parser

import "testing"

func TestParseTableNames(t *testing.T) {
	ddl := `
CREATE TABLE users (id INT);
CREATE TABLE IF NOT EXISTS orders (id INT);
create table "products" (id int);
`
	got := ParseTableNames(ddl)
	if len(got) != 3 {
		t.Fatalf("want 3 tables, got %d: %+v", len(got), got)
	}
}
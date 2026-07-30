package domain

type Index struct {
	Name    string   // may be empty for inline PK
	Table   string
	Columns []string
	Kind    string   // primary_key | unique | index
}
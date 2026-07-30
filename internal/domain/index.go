package domain

type Index struct {
	Name           string
	Table          string
	Columns        []string
	Kind           string
	EstimatedBytes *int64
}
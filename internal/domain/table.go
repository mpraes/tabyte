package domain

type Table struct {
	Name              string
	Columns           []Column
	EstimatedRowBytes *int64
}
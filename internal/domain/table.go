package domain

type Table struct {
	Name                string
	Columns             []Column
	EstimatedRowBytes   *int64
	AssumedRowCount     int64
	EstimatedTableBytes *int64
}
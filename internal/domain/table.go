package domain

type RowCalculation struct {
	ColumnPayloadBytes int64
	RowHeaderBytes     int64
	NullBitmapBytes    int64
	EstimatedRowBytes  int64
}

type Table struct {
	Name                string
	Columns             []Column
	EstimatedRowBytes   *int64
	AssumedRowCount     int64
	EstimatedTableBytes *int64
	Calculation         *RowCalculation
}
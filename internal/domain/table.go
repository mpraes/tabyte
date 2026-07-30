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

	GrowthRowsPerPeriod int64
	GrowthPeriod        string // hour|day|month
	GrowthHorizon       int64
	ProjectedRowCount   *int64
	ProjectedTableBytes *int64
	Indexes             []Index
}
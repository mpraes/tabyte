package domain

type AnalysisSession struct {
	ID         string
	Engine     Engine
	SourceName string
	DDLText    string
	Status     string
	Tables     []Table
	EstimatedTotalBytes *int64
}
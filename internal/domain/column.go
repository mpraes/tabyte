package domain

type Column struct {
	Name             string
	OriginalType     string
	NormalizedType   string
	Length           *int
	Precision        *int
	Scale            *int
	AssumedAvgLength *int
	EstimatedBytes   *int64
}
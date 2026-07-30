package domain

type Column struct {
	Name            string
	OriginalType    string
	NormalizedType  string
	Length          *int // VARCHAR(n), CHAR(n)
	Precision       *int // NUMERIC(p,s)
	Scale           *int
}
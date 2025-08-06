package models

type Request struct {
	ApiKey             *string
	ProductName        *string
	ProductDescription *string
	ProductImage       *string
	StoreName          *string
	Price              *int
	Stock              *int
	PopularityScore    *int
	UrgencyScore       *int
}

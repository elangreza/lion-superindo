package domain

type Product struct {
	ID          int
	Name        string
	Price       int
	ProductType ProductTypes
	BaseDate
}

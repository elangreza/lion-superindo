package domain

type Product struct {
	ID          int
	Name        string
	Quantity    int
	Price       int
	ProductType ProductTypes
	BaseDate
}

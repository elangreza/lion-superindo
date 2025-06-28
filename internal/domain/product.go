package domain

import "time"

type Product struct {
	ID          int
	Name        string
	Price       int
	ProductType ProductType
	CreatedAt   time.Time
}

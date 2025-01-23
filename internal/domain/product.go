package domain

import (
	"encoding/json"
)

type Product struct {
	ID          int
	Name        string
	Quantity    int
	Price       int
	ProductType ProductTypes
	BaseDate
}

func (p *Product) MarshalBinary() ([]byte, error) {
	return json.Marshal(p)
}

func (p *Product) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, p)
}

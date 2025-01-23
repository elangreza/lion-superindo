package domain

import "encoding/json"

type ProductTypes struct {
	Name string
	BaseDate
}

func (p *ProductTypes) MarshalBinary() ([]byte, error) {
	return json.Marshal(p)
}

func (p *ProductTypes) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, p)
}

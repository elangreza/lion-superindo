package domain

import (
	"encoding/json"
	"time"
)

type BaseDate struct {
	CreatedAt time.Time
	UpdatedAt *time.Time
}

func (p *BaseDate) MarshalBinary() ([]byte, error) {
	return json.Marshal(p)
}

func (p *BaseDate) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, p)
}

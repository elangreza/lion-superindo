package domain

import (
	"time"
)

type BaseDate struct {
	CreatedAt time.Time
	UpdatedAt *time.Time
}

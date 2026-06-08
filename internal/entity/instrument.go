package entity

import "time"

type Instrument struct {
	ID        int
	Name      string
	Picture   string
	CreatedAt time.Time
	UpdatedAt time.Time
}

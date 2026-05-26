package entity

import "time"

type Course struct {
	ID        int
	Name      string
	Year      int
	CreatedAt time.Time
	UpdatedAt time.Time
}

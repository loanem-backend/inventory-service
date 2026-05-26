package entity

import "time"

type Toolkit struct {
	ID              int
	KitName         string
	TotalCount      int
	OutOfOrderCount int
	Course          Course
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

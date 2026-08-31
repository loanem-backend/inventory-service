package entity

import (
	"time"

	"github.com/oklog/ulid/v2"
)

type Loan struct {
	ID            ulid.ULID
	Toolkit       Toolkit
	Team          Team
	Submitter     Participant
	Date          time.Time
	SessionNumber int
	Status        LoanStatus
	Note          string
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

type LoanStatus string

const (
	LoanStatusUpcoming  = "UPCOMING"
	LoanStatusOngoing   = "ONGOING"
	LoanStatusDone      = "DONE"
	LoanStatusCancelled = "CANCELLED"
	LoanStatusExpired   = "EXPIRED"
)

type Team struct {
	ID        string
	Number    int
	ClassID   int
	ClassName string
}

type Participant struct {
	ID   string
	Nim  string
	Name string
	// Team Team
}

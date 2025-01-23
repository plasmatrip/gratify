package models

import (
	"database/sql/driver"
	"time"
)

type Status int

const (
	StatusNew Status = iota
	StatusRegistered
	StatusProcessing
	StatusProcessed
	StatusInvalid
	StatusUnknown
)

type Order struct {
	Number  int64     `json:"number"`
	UserID  int32     `json:"-"`
	Status  Status    `json:"status"`
	Accrual float32   `json:"accrual,omitempty"`
	Sum     float32   `json:"-"`
	Date    time.Time `json:"uploaded_at"`
}

type ResponseOrder struct {
	Number  string  `json:"number"`
	Status  string  `json:"status"`
	Accrual float32 `json:"accrual,omitempty"`
	Date    string  `json:"uploaded_at"`
}

type AccrualResponse struct {
	Order   string  `json:"order"`
	Status  string  `json:"status"`
	Accrual float32 `accrual:"accrual"`
}

// deserialize from DB
func (s *Status) Scan(value interface{}) error {
	*s = toStatus(value.(string))
	return nil
}

// serialize to DB
func (s Status) Value() (driver.Value, error) {
	return s.String(), nil
}

func toStatus(s string) Status {
	switch s {
	case "NEW":
		return StatusNew
	case "REGISTERED":
		return StatusRegistered
	case "PROCESSING":
		return StatusProcessing
	case "PROCESSED":
		return StatusProcessed
	case "INVALID":
		return StatusInvalid
	}
	return StatusUnknown
}

func (s Status) String() string {
	switch s {
	case StatusNew:
		return "NEW"
	case StatusRegistered:
		return "REGISTERED"
	case StatusProcessing:
		return "PROCESSING"
	case StatusProcessed:
		return "PROCESSED"
	case StatusInvalid:
		return "INVALID"
	}
	return "unknown"
}

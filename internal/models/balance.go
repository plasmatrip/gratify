package models

import "time"

type Balance struct {
	Current   float32 `json:"current"`
	Withdrawn float32 `json:"withdrawn"`
}

type Withdraw struct {
	Order string    `json:"order"`
	Sum   float32   `json:"sum"`
	Date  time.Time `json:"processed_at,omitempty"`
}

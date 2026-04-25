package models

import "time"

type Slot struct {
	ID        string    `json:"id"`
	VetID     string    `json:"vet_id"`
	StartTime time.Time `json:"start_time"`
	EndTime   time.Time `json:"end_time"`
	SlotLimit int       `json:"slot_limit"`
}

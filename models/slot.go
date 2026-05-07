package models

type Slot struct {
	ID         string `json:"id"`
	VetID      string `json:"vet_id"`
	VetName    string `json:"vet_name"`
	Date       string `json:"date"`
	TimePeriod string `json:"time_period"`
	SlotLimit  int    `json:"slot_limit"`
}

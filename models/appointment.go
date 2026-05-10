package models

import "time"

type Appointment struct {
	ID              string    `json:"id"`
	SlotID          string    `json:"slot_id"`
	PetName         string    `json:"pet_name"`
	PetType         string    `json:"pet_type"`
	ClientName      string    `json:"client_name"`
	Reason          string    `json:"reason"`
	Status          string    `json:"status"`
	AppointmentTime time.Time `json:"appointment_time"`
}

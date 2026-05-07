package services

import (
	"errors"
	"regexp"
	"testing"
	"veterinary-api/repositories"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestGetAvailableSlots_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	rows := sqlmock.NewRows([]string{"id", "vet_id", "vet_name", "date", "time_period", "slot_limit"}).
		AddRow("S-001", "U-002", "สพ.ญ.สมศรี", "2026-05-01", "09:00-10:00", 1)

	query := regexp.QuoteMeta("SELECT slots.id, slots.vet_id, users.name, slots.date, slots.time_period, slots.slot_limit FROM slots JOIN users ON slots.vet_id = users.id WHERE slots.vet_id = ? AND slots.slot_limit > 0")

	mock.ExpectQuery(query).WithArgs("U-002").WillReturnRows(rows)

	repo := repositories.NewSlotRepository(db)
	service := NewSlotService(repo)

	slots, err := service.GetAvailableSlots("U-002")

	assert.NoError(t, err)
	assert.Len(t, slots, 1)
	assert.Equal(t, "S-001", slots[0].ID)
	assert.Equal(t, "สพ.ญ.สมศรี", slots[0].VetName)
}

func TestGetAvailableSlots_Error(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	query := regexp.QuoteMeta("SELECT slots.id, slots.vet_id, users.name, slots.date, slots.time_period, slots.slot_limit FROM slots JOIN users ON slots.vet_id = users.id WHERE slots.vet_id = ? AND slots.slot_limit > 0")

	mock.ExpectQuery(query).WithArgs("U-002").WillReturnError(errors.New("db query error"))

	repo := repositories.NewSlotRepository(db)
	service := NewSlotService(repo)

	slots, err := service.GetAvailableSlots("U-002")

	assert.Error(t, err)
	assert.Nil(t, slots)
	assert.Equal(t, "db query error", err.Error())
}

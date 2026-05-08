package services

import (
	"errors"
	"regexp"
	"testing"
	"veterinary-api/repositories"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

// =====================================================================
// Function GetAvailableSlots
// =====================================================================
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

// =====================================================================
// Function GetAllAvailableSlots
// =====================================================================
func TestGetAllAvailableSlots_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	rows := sqlmock.NewRows([]string{"id", "vet_id", "vet_name", "date", "time_period", "slot_limit"}).
		AddRow("S-001", "U-001", "นสพ.สมชาย", "2026-05-01", "09:00-10:00", 1).
		AddRow("S-003", "U-001", "นสพ.สมชาย", "2026-05-01", "13:00-14:00", 2)

	query := regexp.QuoteMeta("SELECT slots.id, slots.vet_id, users.name, slots.date, slots.time_period, slots.slot_limit FROM slots JOIN users ON slots.vet_id = users.id WHERE slots.slot_limit > 0")

	mock.ExpectQuery(query).WillReturnRows(rows)

	repo := repositories.NewSlotRepository(db)
	service := NewSlotService(repo)

	slots, err := service.GetAllAvailableSlots()

	assert.NoError(t, err)
	assert.Len(t, slots, 2)
}

func TestGetAllAvailableSlots_NotFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	rows := sqlmock.NewRows([]string{"id", "vet_id", "vet_name", "date", "time_period", "slot_limit"})
	query := regexp.QuoteMeta("SELECT slots.id, slots.vet_id, users.name, slots.date, slots.time_period, slots.slot_limit FROM slots JOIN users ON slots.vet_id = users.id WHERE slots.slot_limit > 0")

	mock.ExpectQuery(query).WillReturnRows(rows)

	repo := repositories.NewSlotRepository(db)
	service := NewSlotService(repo)

	_, err = service.GetAllAvailableSlots()

	assert.Error(t, err)
	assert.Equal(t, repositories.ErrNotFound, err)
}

func TestGetAllAvailableSlots_DatabaseError(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	query := regexp.QuoteMeta("SELECT slots.id, slots.vet_id, users.name, slots.date, slots.time_period, slots.slot_limit FROM slots JOIN users ON slots.vet_id = users.id WHERE slots.slot_limit > 0")

	mock.ExpectQuery(query).WillReturnError(errors.New("database connection failed"))

	repo := repositories.NewSlotRepository(db)
	service := NewSlotService(repo)

	slots, err := service.GetAllAvailableSlots()

	assert.Error(t, err)
	assert.Nil(t, slots)
	assert.Equal(t, "database connection failed", err.Error())
}

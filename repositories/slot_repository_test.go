package repositories

import (
	"errors"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

// ทดสอบว่า SQL Query ทำงานถูกต้องและดึงข้อมูลมาใส่ Slice ได้ครบ
func TestGetAllAvailableSlots_Repo_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	// Mock ข้อมูลที่คาดว่าจะได้จาก DB
	rows := sqlmock.NewRows([]string{"id", "vet_id", "name", "date", "time_period", "slot_limit"}).
		AddRow("S-001", "U-001", "นสพ.สมชาย", "2026-05-01", "09:00-10:00", 1).
		AddRow("S-002", "U-002", "นสพ.หญิง", "2026-05-01", "10:00-11:00", 5)

	// ตรวจสอบ SQL Statement
	query := regexp.QuoteMeta("SELECT slots.id, slots.vet_id, users.name, slots.date, slots.time_period, slots.slot_limit FROM slots JOIN users ON slots.vet_id = users.id WHERE slots.slot_limit > 0")

	mock.ExpectQuery(query).WillReturnRows(rows)

	repo := NewSlotRepository(db)
	slots, err := repo.GetAllAvailableSlots()

	assert.NoError(t, err)
	assert.Len(t, slots, 2)
	assert.Equal(t, "นสพ.สมชาย", slots[0].VetName)
}

// ทดสอบ Logic 'if len(slots) == 0'
func TestGetAllAvailableSlots_Repo_NotFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	// จำลองกรณีไม่มีข้อมูลส่งกลับมาเลย
	rows := sqlmock.NewRows([]string{"id", "vet_id", "name", "date", "time_period", "slot_limit"})

	query := regexp.QuoteMeta("SELECT slots.id, slots.vet_id, users.name, slots.date, slots.time_period, slots.slot_limit FROM slots JOIN users ON slots.vet_id = users.id WHERE slots.slot_limit > 0")

	mock.ExpectQuery(query).WillReturnRows(rows)

	repo := NewSlotRepository(db)
	slots, err := repo.GetAllAvailableSlots()

	// Error ErrNotFound
	assert.Error(t, err)
	assert.Equal(t, ErrNotFound, err)
	assert.Nil(t, slots)
}

// r.DB.Query พัง
func TestGetAllAvailableSlots_Repo_QueryError(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	// DB พ่น Error ทันที่ Query
	mock.ExpectQuery("SELECT").WillReturnError(errors.New("db connection failed"))

	repo := NewSlotRepository(db)
	_, err = repo.GetAllAvailableSlots()

	assert.Error(t, err)
}

// Scan ข้อมูลไผิดพลาด
func TestGetAllAvailableSlots_Repo_ScanError(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	// ส่ง Column มาไม่ครบ
	rows := sqlmock.NewRows([]string{"id"}).AddRow(1)

	mock.ExpectQuery("SELECT").WillReturnRows(rows)

	repo := NewSlotRepository(db)
	_, err = repo.GetAllAvailableSlots()

	assert.Error(t, err)
}

package repositories

import (
	"errors"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

// =====================================================================
// Test: GetAllAvailableSlots Repository
// =====================================================================

// ดึงข้อมูลสำเร็จ (Success)
func TestGetAllAvailableSlots_Repo_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	rows := sqlmock.NewRows([]string{"id", "vet_id", "name", "date", "time_period", "slot_limit"}).
		AddRow("S-001", "U-001", "นสพ.สมชาย", "2026-05-01", "09:00-10:00", 1).
		AddRow("S-002", "U-002", "นสพ.หญิง", "2026-05-01", "10:00-11:00", 5)

	query := regexp.QuoteMeta("SELECT slots.id, slots.vet_id, users.name, slots.date, slots.time_period, slots.slot_limit FROM slots JOIN users ON slots.vet_id = users.id WHERE slots.slot_limit > 0")

	mock.ExpectQuery(query).WillReturnRows(rows)

	repo := NewSlotRepository(db)
	slots, err := repo.GetAllAvailableSlots()

	assert.NoError(t, err)
	assert.Len(t, slots, 2)
	assert.Equal(t, "นสพ.สมชาย", slots[0].VetName)
}

// ไม่มีคิวว่างเลยในระบบ
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

	assert.Error(t, err)
	assert.Equal(t, ErrNotFound, err)
	assert.Nil(t, slots)
}

// ดาต้าเบสพังตั้งแต่ตอน Query
func TestGetAllAvailableSlots_Repo_QueryError(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	mock.ExpectQuery("SELECT").WillReturnError(errors.New("db connection failed"))

	repo := NewSlotRepository(db)
	_, err = repo.GetAllAvailableSlots()

	assert.Error(t, err)
}

// Scan ข้อมูลผิดพลาด
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

// =====================================================================
// Test: GetAvailableSlots Repository
// =====================================================================

// ดึงข้อมูลคิวว่างของหมอ "(vetID)" สำเร็จ
func TestGetAvailableSlots_Repo_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	vetID := "V-001"
	rows := sqlmock.NewRows([]string{"id", "vet_id", "name", "date", "time_period", "slot_limit"}).
		AddRow("S-001", vetID, "นสพ.สมชาย", "2026-05-01", "09:00", 1)

	query := regexp.QuoteMeta("SELECT slots.id, slots.vet_id, users.name, slots.date, slots.time_period, slots.slot_limit FROM slots JOIN users ON slots.vet_id = users.id WHERE slots.vet_id = ? AND slots.slot_limit > 0")

	mock.ExpectQuery(query).WithArgs(vetID).WillReturnRows(rows)

	repo := NewSlotRepository(db)
	slots, err := repo.GetAvailableSlots(vetID)

	assert.NoError(t, err)
	assert.Len(t, slots, 1)
	assert.Equal(t, vetID, slots[0].VetID)
}

// ค้นหาหมอด้วยรหัสนี้แล้ว แต่ไม่พบคิวว่าง
func TestGetAvailableSlots_Repo_NotFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	vetID := "V-001"
	rows := sqlmock.NewRows([]string{"id", "vet_id", "name", "date", "time_period", "slot_limit"})

	query := regexp.QuoteMeta("SELECT slots.id, slots.vet_id, users.name, slots.date, slots.time_period, slots.slot_limit FROM slots JOIN users ON slots.vet_id = users.id WHERE slots.vet_id = ? AND slots.slot_limit > 0")
	mock.ExpectQuery(query).WithArgs(vetID).WillReturnRows(rows)

	repo := NewSlotRepository(db)
	slots, err := repo.GetAvailableSlots(vetID)

	assert.Error(t, err)
	assert.Nil(t, slots)
	assert.Equal(t, ErrNotFound, err)
}

// Database พังระหว่างดึงข้อมูล
func TestGetAvailableSlots_Repo_QueryError(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	vetID := "V-001"

	query := regexp.QuoteMeta("SELECT")
	mock.ExpectQuery(query).WithArgs(vetID).WillReturnError(errors.New("db connection failed"))

	repo := NewSlotRepository(db)
	_, err = repo.GetAvailableSlots(vetID)

	assert.Error(t, err)
}

// คิวว่างของหมอที่ส่งกลับมามีโครงสร้างผิดปกติจน scan ไม่ได้
func TestGetAvailableSlots_Repo_ScanError(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	vetID := "V-001"

	rows := sqlmock.NewRows([]string{"id"}).AddRow("S-001")

	query := regexp.QuoteMeta("SELECT slots.id, slots.vet_id, users.name, slots.date, slots.time_period, slots.slot_limit FROM slots JOIN users ON slots.vet_id = users.id WHERE slots.vet_id = ? AND slots.slot_limit > 0")
	mock.ExpectQuery(query).WithArgs(vetID).WillReturnRows(rows)

	repo := NewSlotRepository(db)
	_, err = repo.GetAvailableSlots(vetID)

	assert.Error(t, err)
}

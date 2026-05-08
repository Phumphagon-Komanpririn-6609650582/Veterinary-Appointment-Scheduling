package repositories

import (
	"errors"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

// =====================================================================
// Test: GetAllVets
// =====================================================================

// ดึงข้อมูลสำเร็จ (Success)
func TestGetAllVets_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	rows := sqlmock.NewRows([]string{"id", "username", "name", "role"}).
		AddRow("V-001", "dr_somchai", "นสพ.สมชาย", "vet").
		AddRow("V-002", "dr_ying", "นสพ.หญิง", "vet")

	query := regexp.QuoteMeta("SELECT id, username, name, role FROM users WHERE role = 'vet'")
	mock.ExpectQuery(query).WillReturnRows(rows)

	repo := NewVetRepository(db)
	vets, err := repo.GetAllVets()

	assert.NoError(t, err)
	assert.Len(t, vets, 2)
	assert.Equal(t, "V-001", vets[0].ID)
	assert.Equal(t, "นสพ.สมชาย", vets[0].Name)
}

// ดาต้าเบสพังตั้งแต่ตอน Query
func TestGetAllVets_QueryError(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	query := regexp.QuoteMeta("SELECT id, username, name, role FROM users WHERE role = 'vet'")
	mock.ExpectQuery(query).WillReturnError(errors.New("db connection timeout"))

	repo := NewVetRepository(db)
	vets, err := repo.GetAllVets()

	assert.Error(t, err)
	assert.Nil(t, vets)
}

// ดึงข้อมูลได้ แต่โครงสร้างพังทำให้ Scan Error
func TestGetAllVets_ScanError(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	rows := sqlmock.NewRows([]string{"id", "username"}).
		AddRow("V-001", "dr_somchai")

	query := regexp.QuoteMeta("SELECT id, username, name, role FROM users WHERE role = 'vet'")
	mock.ExpectQuery(query).WillReturnRows(rows)

	repo := NewVetRepository(db)
	vets, err := repo.GetAllVets()

	assert.Error(t, err)
	assert.Nil(t, vets)
}

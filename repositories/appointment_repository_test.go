package repositories

import (
	"database/sql"
	"testing"

	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
)

func TestAppointmentRepository_AllCases(t *testing.T) {
	db, _ := sql.Open("sqlite3", "../veterinary.db")
	repo := NewAppointmentRepository(db)

	// --- เคสที่ 1: ดึงข้อมูลสำเร็จ ---
	t.Run("GetByID_Success", func(t *testing.T) {
		app, err := repo.GetByID("A-001")
		assert.NoError(t, err)
		assert.NotNil(t, app)
		assert.Equal(t, "A-001", app.ID)
	})

	// --- เคสที่ 2: ยกเลิกสำเร็จ ---
	t.Run("CancelAppointment_Success", func(t *testing.T) {
		err := repo.CancelAppointment("A-001")
		assert.NoError(t, err)

		app, _ := repo.GetByID("A-001")
		assert.Equal(t, "cancelled", app.Status)
	})

	// --- เคสที่ 3: ยกเลิก ID ที่ไม่มีอยู่จริง (เก็บตก RowsAffected == 0) ---
	t.Run("CancelAppointment_NotFound", func(t *testing.T) {
		err := repo.CancelAppointment("ID-HAVE-NOT")
		assert.Error(t, err)
		// เช็คข้อความให้ตรงกับในโค้ด Repo (ภาษาอังกฤษ)
		assert.Equal(t, "No appointment information was found.", err.Error())
	})
}

// --- เคสที่ 4: บังคับ Error เพื่อเก็บตกบรรทัด return err ---
func TestAppointmentRepository_DB_Errors(t *testing.T) {
	db, _ := sql.Open("sqlite3", ":memory:")
	repo := NewAppointmentRepository(db)
	db.Close() // บังคับปิดการเชื่อมต่อ

	t.Run("GetByID_DB_Closed", func(t *testing.T) {
		_, err := repo.GetByID("any-id")
		assert.Error(t, err)
	})

	t.Run("CancelAppointment_DB_Closed", func(t *testing.T) {
		err := repo.CancelAppointment("any-id")
		assert.Error(t, err)
	})
}

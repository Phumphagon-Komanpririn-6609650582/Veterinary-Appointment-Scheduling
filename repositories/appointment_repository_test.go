package repositories

import (
	"database/sql"
	"testing"
	"veterinary-api/models"

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

	// ==========================================
	// 👩‍💻 เทสของนุช (UpdateAppointment)
	// ==========================================
	t.Run("UpdateAppointment_Success", func(t *testing.T) {
		// จำลองข้อมูลที่จะใช้ในการอัปเดต
		appToUpdate := models.Appointment{
			ID:         "A-001", // ต้องเป็น ID ที่มีอยู่จริงใน veterinary.db
			SlotID:     "S-001",
			PetName:    "น้องเหมียวอัปเดต",
			PetType:    "Cat",
			ClientName: "คุณนุชชี่",
			Reason:     "ตรวจสุขภาพประจำปี",
		}

		err := repo.UpdateAppointment(appToUpdate)
		assert.NoError(t, err) // ต้องไม่มี Error
	})

	t.Run("UpdateAppointment_NotFound", func(t *testing.T) {
		appToUpdate := models.Appointment{
			ID: "ID-HAVE-NOT", // ID มั่ว
		}

		err := repo.UpdateAppointment(appToUpdate)
		assert.Error(t, err) // ต้องเกิด Error
		// เช็กว่าข้อความ Error ตรงกับที่นุชเขียนดักไว้ (RowsAffected == 0)
		assert.Equal(t, "appointment not found", err.Error())
	})
	// ==========================================

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

	// 👩‍💻 เทสของนุช: กรณี Database มีปัญหา (บรรทัด if err != nil)
	t.Run("UpdateAppointment_DB_Closed", func(t *testing.T) {
		err := repo.UpdateAppointment(models.Appointment{ID: "any-id"})
		assert.Error(t, err)
	})

	t.Run("CancelAppointment_DB_Closed", func(t *testing.T) {
		err := repo.CancelAppointment("any-id")
		assert.Error(t, err)
	})
}
// --- เคสที่ 5: เก็บตกบรรทัดสีแดง (time.ParseInLocation Error) ---
func TestAppointmentRepository_GetByID_TimeParseError(t *testing.T) {
	// สร้าง Database จำลองใน Memory
	db, _ := sql.Open("sqlite3", ":memory:")
	defer db.Close()
	repo := NewAppointmentRepository(db)

	// สร้างตารางชั่วคราว
	db.Exec(`CREATE TABLE slots (id TEXT PRIMARY KEY, date TEXT, time_period TEXT);`)
	db.Exec(`CREATE TABLE appointments (id TEXT PRIMARY KEY, slot_id TEXT, pet_name TEXT, pet_type TEXT, client_name TEXT, reason TEXT, status TEXT);`)

	// 📌 จุดสำคัญ: แกล้งใส่ข้อมูลวันที่แบบพังๆ ("NOT-A-DATE") ลงไป
	db.Exec(`INSERT INTO slots (id, date, time_period) VALUES ('S-ERR', 'NOT-A-DATE', 'XX:XX');`)
	db.Exec(`INSERT INTO appointments (id, slot_id, pet_name, pet_type, client_name, reason, status) 
			 VALUES ('A-ERR', 'S-ERR', 'หมา', 'Dog', 'คุณนุช', 'ป่วย', 'pending');`)

	// เรียกใช้งาน GetByID
	_, err := repo.GetByID("A-ERR")

	// ต้องเกิด Error ตอนแปลงเวลา
	assert.Error(t, err)
}
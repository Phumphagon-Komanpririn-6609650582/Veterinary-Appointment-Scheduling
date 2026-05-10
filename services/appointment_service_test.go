package services

import (
	"database/sql"
	"testing"
	"veterinary-api/models"
	"veterinary-api/repositories"

	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
)

// =====================================================================
// 👨‍💻 พื้นที่ของ: เอลฟ์ (Test CancelAppointment)
// =====================================================================
func TestAppointmentService_CancelAppointment(t *testing.T) {
	db, _ := sql.Open("sqlite3", ":memory:")
	defer db.Close()

	db.Exec(`CREATE TABLE slots (id TEXT PRIMARY KEY, date TEXT, time_period TEXT);`)
	db.Exec(`CREATE TABLE appointments (id TEXT PRIMARY KEY, slot_id TEXT, pet_name TEXT, pet_type TEXT, client_name TEXT, reason TEXT, status TEXT);`)

	// ใส่ข้อมูลให้ครบทุกคอลัมน์ กัน Error NULL ตอนดึงข้อมูล (GetByID)
	db.Exec(`INSERT INTO slots (id, date, time_period) VALUES ('S-003', '2099-12-31', '10:00-11:00');`)
	db.Exec(`INSERT INTO appointments (id, slot_id, pet_name, pet_type, client_name, reason, status) VALUES ('A-003', 'S-003', '-', '-', '-', '-', 'pending');`)

	db.Exec(`INSERT INTO slots (id, date, time_period) VALUES ('S-002', '2020-01-01', '10:00-11:00');`)
	db.Exec(`INSERT INTO appointments (id, slot_id, pet_name, pet_type, client_name, reason, status) VALUES ('A-002', 'S-002', '-', '-', '-', '-', 'pending');`)

	db.Exec(`INSERT INTO slots (id, date, time_period) VALUES ('S-001', '2099-12-31', '10:00-11:00');`)
	db.Exec(`INSERT INTO appointments (id, slot_id, pet_name, pet_type, client_name, reason, status) VALUES ('A-001', 'S-001', '-', '-', '-', '-', 'cancelled');`)

	repo := repositories.NewAppointmentRepository(db)
	service := NewAppointmentService(repo)

	t.Run("Cancel_Success", func(t *testing.T) {
		err := service.CancelAppointment("A-003")
		assert.NoError(t, err)
	})

	t.Run("Cancel_TooLate_Error", func(t *testing.T) {
		err := service.CancelAppointment("A-002")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "at least 2 hours in advance")
	})

	t.Run("Error_NotFound", func(t *testing.T) {
		err := service.CancelAppointment("NON-EXISTENT-ID")
		assert.Error(t, err)
		assert.Equal(t, "This appointment was not found in the system.", err.Error())
	})

	t.Run("Error_AlreadyCancelled", func(t *testing.T) {
		err := service.CancelAppointment("A-001")
		assert.Error(t, err)
		assert.Equal(t, "This appointment has been cancelled.", err.Error())
	})
}

// =====================================================================
// 👩‍💻 พื้นที่ของ: นุช (Test UpdateAppointment)
// =====================================================================
func TestAppointmentService_UpdateAppointment(t *testing.T) {
	db, _ := sql.Open("sqlite3", ":memory:")
	defer db.Close()

	db.Exec(`CREATE TABLE slots (id TEXT PRIMARY KEY, date TEXT, time_period TEXT);`)
	db.Exec(`CREATE TABLE appointments (id TEXT PRIMARY KEY, slot_id TEXT, pet_name TEXT, pet_type TEXT, client_name TEXT, reason TEXT, status TEXT);`)

	db.Exec(`INSERT INTO slots (id, date, time_period) VALUES ('S-001', '2099-12-31', '10:00-11:00');`)
	
	// เคสปกติ (ใส่ให้ครบทุกคอลัมน์)
	db.Exec(`INSERT INTO appointments (id, slot_id, pet_name, pet_type, client_name, reason, status) VALUES ('A-004', 'S-001', '-', '-', '-', '-', 'pending');`) 
	// เคสโดนยกเลิก (ใส่ให้ครบทุกคอลัมน์)
	db.Exec(`INSERT INTO appointments (id, slot_id, pet_name, pet_type, client_name, reason, status) VALUES ('A-001', 'S-001', '-', '-', '-', '-', 'cancelled');`) 

	repo := repositories.NewAppointmentRepository(db)
	service := NewAppointmentService(repo)

	// เคสที่ 1: อัปเดตสำเร็จ
	t.Run("Update_Success", func(t *testing.T) {
		appToUpdate := models.Appointment{
			ID:         "A-004", 
			SlotID:     "S-002",
			PetName:    "น้องด่าง",
			PetType:    "Dog",
			ClientName: "คุณนุช",
			Reason:     "อาบน้ำตัดขน",
		}
		err := service.UpdateAppointment(appToUpdate)
		assert.NoError(t, err)
	})

	// เคสที่ 2: หา ID ไม่เจอในระบบ
	t.Run("Update_Fail_NotFound", func(t *testing.T) {
		appToUpdate := models.Appointment{ID: "ID-HAVE-NOT"} 
		err := service.UpdateAppointment(appToUpdate)
		
		assert.Error(t, err)
		assert.Equal(t, "appointment not found", err.Error())
	})

	// เคสที่ 3: ห้ามแก้ไขคิวที่โดนยกเลิกไปแล้ว
	t.Run("Update_Fail_Cancelled", func(t *testing.T) {
		appToUpdate := models.Appointment{ID: "A-001"} 
		err := service.UpdateAppointment(appToUpdate)
		
		assert.Error(t, err)
		assert.Equal(t, "cannot edit cancelled appointment", err.Error())
	})
}
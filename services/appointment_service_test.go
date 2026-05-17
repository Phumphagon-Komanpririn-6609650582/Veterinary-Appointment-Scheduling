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

// =====================================================================
// 👩‍💻 พื้นที่ของ: ปลา (Test CreateAppointment)
// =====================================================================
func TestAppointmentService_CreateAppointment(t *testing.T) {

	db, _ := sql.Open("sqlite3", ":memory:")
	defer db.Close()

	// สร้าง table slots
	db.Exec(`
		CREATE TABLE slots (
			id TEXT PRIMARY KEY,
			slot_limit INTEGER
		);
	`)

	// สร้าง table appointments
	db.Exec(`
		CREATE TABLE appointments (
			id TEXT PRIMARY KEY,
			slot_id TEXT,
			pet_name TEXT,
			pet_type TEXT,
			client_name TEXT,
			reason TEXT,
			status TEXT
		);
	`)

	// slot ปกติ
	db.Exec(`
		INSERT INTO slots (id, slot_limit)
		VALUES ('S-001', 1);
	`)

	// slot เต็ม
	db.Exec(`
		INSERT INTO slots (id, slot_limit)
		VALUES ('S-002', 0);
	`)

	repo := repositories.NewAppointmentRepository(db)
	service := NewAppointmentService(repo)

	// =====================================================
	// success
	// =====================================================
	t.Run("Create_Success", func(t *testing.T) {

		app := &models.Appointment{
			SlotID:     "S-001",
			PetName:    "Lucky",
			PetType:    "Dog",
			ClientName: "ปลา",
			Reason:     "ตรวจสุขภาพ",
		}

		err := service.CreateAppointment(app)

		assert.NoError(t, err)
		assert.Equal(t, "in-progress", app.Status)
	})

	// =====================================================
	// pet name required
	// =====================================================
	t.Run("PetName_Required", func(t *testing.T) {

		app := &models.Appointment{
			SlotID:     "S-001",
			PetType:    "Dog",
			ClientName: "ปลา",
			Reason:     "ตรวจสุขภาพ",
		}

		err := service.CreateAppointment(app)

		assert.Error(t, err)
		assert.Equal(t, "pet name is required", err.Error())
	})

	// =====================================================
	// pet type required
	// =====================================================
	t.Run("PetType_Required", func(t *testing.T) {

		app := &models.Appointment{
			SlotID:     "S-001",
			PetName:    "Lucky",
			ClientName: "ปลา",
			Reason:     "ตรวจสุขภาพ",
		}

		err := service.CreateAppointment(app)

		assert.Error(t, err)
		assert.Equal(t, "pet type is required", err.Error())
	})

	// =====================================================
	// client name required
	// =====================================================
	t.Run("ClientName_Required", func(t *testing.T) {

		app := &models.Appointment{
			SlotID:  "S-001",
			PetName: "Lucky",
			PetType: "Dog",
			Reason:  "ตรวจสุขภาพ",
		}

		err := service.CreateAppointment(app)

		assert.Error(t, err)
		assert.Equal(t, "client name is required", err.Error())
	})

	// =====================================================
	// reason required
	// =====================================================
	t.Run("Reason_Required", func(t *testing.T) {

		app := &models.Appointment{
			SlotID:     "S-001",
			PetName:    "Lucky",
			PetType:    "Dog",
			ClientName: "ปลา",
		}

		err := service.CreateAppointment(app)

		assert.Error(t, err)
		assert.Equal(t, "reason is required", err.Error())
	})

	// =====================================================
	// slot id required
	// =====================================================
	t.Run("SlotID_Required", func(t *testing.T) {

		app := &models.Appointment{
			PetName:    "Lucky",
			PetType:    "Dog",
			ClientName: "ปลา",
			Reason:     "ตรวจสุขภาพ",
		}

		err := service.CreateAppointment(app)

		assert.Error(t, err)
		assert.Equal(t, "slot id is required", err.Error())
	})

	// =====================================================
	// slot not found
	// =====================================================
	t.Run("Slot_NotFound", func(t *testing.T) {

		app := &models.Appointment{
			SlotID:     "S-999",
			PetName:    "Lucky",
			PetType:    "Dog",
			ClientName: "ปลา",
			Reason:     "ตรวจสุขภาพ",
		}

		err := service.CreateAppointment(app)

		assert.Error(t, err)
		assert.Equal(t, "slot not found", err.Error())
	})

	// =====================================================
	// slot full
	// =====================================================
	t.Run("Slot_Full", func(t *testing.T) {

		app := &models.Appointment{
			SlotID:     "S-002",
			PetName:    "Lucky",
			PetType:    "Dog",
			ClientName: "ปลา",
			Reason:     "ตรวจสุขภาพ",
		}

		err := service.CreateAppointment(app)

		assert.Error(t, err)
		assert.Equal(t, "slot is full", err.Error())
	})
}

// =====================================================================
// 👩‍💻 พื้นที่ของ: ปลา (Cover Interface Functions)
// =====================================================================

func TestAppointmentService_InterfaceFunctions(t *testing.T) {

	db, _ := sql.Open("sqlite3", ":memory:")
	defer db.Close()

	// สร้าง table ให้ครบ
	db.Exec(`
		CREATE TABLE slots (
			id TEXT PRIMARY KEY,
			slot_limit INTEGER
		);
	`)

	db.Exec(`
		CREATE TABLE appointments (
			id TEXT PRIMARY KEY,
			slot_id TEXT,
			pet_name TEXT,
			pet_type TEXT,
			client_name TEXT,
			reason TEXT,
			status TEXT
		);
	`)

	// เพิ่ม slot
	db.Exec(`
		INSERT INTO slots (id, slot_limit)
		VALUES ('S-001', 1);
	`)

	// เพิ่ม appointment เพื่อให้ UpdateStatus หาเจอ
	db.Exec(`
		INSERT INTO appointments (id, slot_id, pet_name, pet_type, client_name, reason, status)
		VALUES ('A-001', 'S-001', 'Lucky', 'Dog', 'ปลา', 'ตรวจสุขภาพ', 'in-progress');
	`)

	repo := repositories.NewAppointmentRepository(db)
	service := NewAppointmentService(repo)

	// =====================================================
	// GetAppointments
	// =====================================================
	t.Run("GetAppointments", func(t *testing.T) {

		result, err := service.GetAppointments()

		assert.NoError(t, err)
		assert.Nil(t, result)
	})

	// =====================================================
	// UpdateStatus
	// =====================================================
	t.Run("UpdateStatus", func(t *testing.T) {

		err := service.UpdateStatus("A-001", "done")

		assert.NoError(t, err)
	})

	// =====================================================
	// CreateAppointment Success
	// =====================================================
	t.Run("CreateAppointment_Success", func(t *testing.T) {

		app := &models.Appointment{
			SlotID:     "S-001",
			PetName:    "Lucky",
			PetType:    "Dog",
			ClientName: "ปลา",
			Reason:     "ตรวจสุขภาพ",
		}

		err := service.CreateAppointment(app)

		assert.NoError(t, err)
		assert.Equal(t, "in-progress", app.Status)
	})

	// =====================================================
	// CreateAppointment Insert Error
	// =====================================================
	t.Run("CreateAppointment_Insert_Error", func(t *testing.T) {

		errorDB, _ := sql.Open("sqlite3", ":memory:")
		defer errorDB.Close()

		// มีแค่ slots
		errorDB.Exec(`
			CREATE TABLE slots (
				id TEXT PRIMARY KEY,
				slot_limit INTEGER
			);
		`)

		errorDB.Exec(`
			INSERT INTO slots (id, slot_limit)
			VALUES ('S-003', 1);
		`)

		errorRepo := repositories.NewAppointmentRepository(errorDB)
		errorService := NewAppointmentService(errorRepo)

		app := &models.Appointment{
			SlotID:     "S-003",
			PetName:    "Lucky",
			PetType:    "Dog",
			ClientName: "ปลา",
			Reason:     "ตรวจสุขภาพ",
		}

		err := errorService.CreateAppointment(app)

		assert.Error(t, err)
		assert.Equal(t, "failed to create appointment", err.Error())
	})
}

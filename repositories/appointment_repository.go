package repositories

import (
	"database/sql"
	"errors"
	"time"
	"veterinary-api/models"
)

type AppointmentRepository struct {
	DB *sql.DB
}

func NewAppointmentRepository(db *sql.DB) *AppointmentRepository {
	return &AppointmentRepository{DB: db}
}

func (r *AppointmentRepository) GetByID(id string) (*models.Appointment, error) {
	var app models.Appointment
	var dateStr, timeStr string

	// ใช้ JOIN เพื่อไปเอาวันที่และเวลาจากตาราง slots
	// SUBSTR(s.time_period, 1, 5) จะดึงแค่ "09:00" ออกมาจาก "09:00-10:00"
	query := `
		SELECT 
			a.id, a.slot_id, a.pet_name, a.pet_type, a.client_name, a.reason, a.status,
			s.date, SUBSTR(s.time_period, 1, 5) as start_time
		FROM appointments a
		JOIN slots s ON a.slot_id = s.id
		WHERE a.id = ?`

	err := r.DB.QueryRow(query, id).Scan(
		&app.ID, &app.SlotID, &app.PetName, &app.PetType,
		&app.ClientName, &app.Reason, &app.Status,
		&dateStr, &timeStr,
	)
	if err != nil {
		return nil, err
	}

	// รวมร่าง String เป็นฟอร์แมตที่ Go เข้าใจ: "2026-05-10 09:00"
	fullDateTimeStr := dateStr + " " + timeStr
	layout := "2006-01-02 15:04"

	// แปลงเป็น time.Time เพื่อส่งไปให้ Service คำนวณ
	parsedTime, err := time.ParseInLocation(layout, fullDateTimeStr, time.Local)
	if err != nil {
		return nil, err
	}

	app.AppointmentTime = parsedTime
	return &app, nil
}

// =====================================================================
// 👩‍💻 พื้นที่ของ: ปลา (Insert นัดใหม่)
// =====================================================================
func (r *AppointmentRepository) CreateAppointment() {

}

// =====================================================================
// 👩‍💻 พื้นที่ของ: นุช (อัปเดตข้อมูลนัดหมาย)
// =====================================================================
func (r *AppointmentRepository) UpdateAppointment() {

}

// =====================================================================
// 👨‍💻 พื้นที่ของ: เอลฟ์ (ลบข้อมูลนัดหมาย)
// =====================================================================
func (r *AppointmentRepository) CancelAppointment(id string) error {
	query := "UPDATE appointments SET status = 'cancelled' WHERE id = ?"
	result, err := r.DB.Exec(query, id)
	if err != nil {
		return err
	}
	// เช็คว่ามีแถวไหนโดนแก้ไหม (ถ้าเป็น 0 แปลว่าหา ID นี้ไม่เจอ)
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return errors.New("No appointment information was found.")
	}

	return nil
}

// =====================================================================
// 👨‍💻 พื้นที่ของ: พี่สิรภพ (ดึงรายการนัดพร้อม Filter)
// =====================================================================
func (r *AppointmentRepository) GetAppointments() {

}

// =====================================================================
// 👨‍💻 พื้นที่ของ: พี่อิทธิเชษฐ์ (อัปเดตสถานะการรักษา)
// =====================================================================
func (r *AppointmentRepository) UpdateStatus() {

}

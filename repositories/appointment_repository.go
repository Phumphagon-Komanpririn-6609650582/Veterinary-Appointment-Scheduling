package repositories

import (
	"database/sql"
	"errors"
	"fmt"
	"strconv"
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
func (r *AppointmentRepository) CreateAppointment(app *models.Appointment) error {

	// สร้าง ID ใหม่
	var lastID string

	queryLastID := `
		SELECT id
		FROM appointments
		ORDER BY id DESC
		LIMIT 1
	`

	r.DB.QueryRow(queryLastID).Scan(&lastID)

	newID := "A-001"

	if lastID != "" {
		numberPart := lastID[2:]
		number, _ := strconv.Atoi(numberPart)
		number++

		newID = fmt.Sprintf("A-%03d", number)
	}

	app.ID = newID

	query := `
		INSERT INTO appointments (
			id,
			slot_id,
			pet_name,
			pet_type,
			client_name,
			reason,
			status
		)
		VALUES (?, ?, ?, ?, ?, ?, ?)
	`

	_, err := r.DB.Exec(
		query,
		app.ID,
		app.SlotID,
		app.PetName,
		app.PetType,
		app.ClientName,
		app.Reason,
		app.Status,
	)

	if err != nil {
		return err
	}

	return nil
}

// เช็กว่า slot มีอยู่จริงไหม
func (r *AppointmentRepository) CheckSlotExists(slotID string) bool {

	query := `SELECT COUNT(*) FROM slots WHERE id = ?`

	var count int

	err := r.DB.QueryRow(query, slotID).Scan(&count)

	if err != nil {
		return false
	}

	return count > 0
}

// ลดจำนวน slot หลังจองสำเร็จ
func (r *AppointmentRepository) DecreaseSlotLimit(slotID string) error {

	query := `
		UPDATE slots
		SET slot_limit = slot_limit - 1
		WHERE id = ? AND slot_limit > 0
	`

	result, err := r.DB.Exec(query, slotID)

	if err != nil {
		return err
	}

	rows, _ := result.RowsAffected()

	if rows == 0 {
		return errors.New("slot is full")
	}

	return nil
}

// =====================================================================
// 👩‍💻 พื้นที่ของ: นุช (อัปเดตข้อมูลนัดหมาย)
// =====================================================================
func (r *AppointmentRepository) UpdateAppointment(app models.Appointment) error {

	query := `
		UPDATE appointments
SET
    slot_id = ?,
    pet_name = ?,
    pet_type = ?,
    client_name = ?,
    reason = ?
WHERE id = ?`

	result, err := r.DB.Exec(
		query,
		app.SlotID,
		app.PetName,
		app.PetType,
		app.ClientName,
		app.Reason,
		app.ID,
	)

	if err != nil {
		return err
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return errors.New("appointment not found")
	}

	return nil
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
func (r *AppointmentRepository) GetAllAppointments() ([]models.Appointment, error) {
	query := `
		SELECT COALESCE(id, ''),slot_id,pet_name,pet_type,client_name,reason,status
		FROM appointments
	`
	rows, err := r.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var appointments []models.Appointment

	for rows.Next() {
		var appointment models.Appointment
		if err := rows.Scan(
			&appointment.ID,
			&appointment.SlotID,
			&appointment.PetName,
			&appointment.PetType,
			&appointment.ClientName,
			&appointment.Reason,
			&appointment.Status); err != nil {
			return nil, err
		}
		appointments = append(appointments, appointment)
	}

	if len(appointments) == 0 {
		return nil, ErrNotFound
	}

	return appointments, nil

}
func (r *AppointmentRepository) GetAppointmentsByVet(vetID string) ([]models.Appointment, error) {
	query := `
		SELECT COALESCE(appointments.id, ''),appointments.slot_id,appointments.pet_name,appointments.pet_type,appointments.client_name,appointments.reason,appointments.status
		FROM appointments
		INNER JOIN slots ON appointments.slot_id = slots.id
		WHERE slots.vet_id = ?
	`
	rows, err := r.DB.Query(query, vetID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var appointments []models.Appointment
	for rows.Next() {
		var appointment models.Appointment
		if err := rows.Scan(
			&appointment.ID,
			&appointment.SlotID,
			&appointment.PetName,
			&appointment.PetType,
			&appointment.ClientName,
			&appointment.Reason,
			&appointment.Status,
		); err != nil {
			return nil, err
		}
		appointments = append(appointments, appointment)
	}

	if len(appointments) == 0 {
		return nil, ErrNotFound
	}

	return appointments, nil
}

func (r *AppointmentRepository) GetAppointmentsByDate(Date string) ([]models.Appointment, error) {
	query := `
		SELECT  COALESCE(appointments.id, ''),appointments.slot_id,appointments.pet_name,appointments.pet_type,appointments.client_name,appointments.reason,appointments.status
		FROM appointments
		INNER JOIN slots ON appointments.slot_id = slots.id
		WHERE slots.date = ?
	`
	rows, err := r.DB.Query(query, Date)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var appointments []models.Appointment
	for rows.Next() {
		var appointment models.Appointment
		if err := rows.Scan(
			&appointment.ID,
			&appointment.SlotID,
			&appointment.PetName,
			&appointment.PetType,
			&appointment.ClientName,
			&appointment.Reason,
			&appointment.Status,
		); err != nil {
			return nil, err
		}
		appointments = append(appointments, appointment)
	}

	if len(appointments) == 0 {
		return nil, ErrNotFound
	}

	return appointments, nil
}


// =====================================================================
// 👨‍💻 พื้นที่ของ: พี่อิทธิเชษฐ์ (อัปเดตสถานะการรักษา)
// =====================================================================
func (r *AppointmentRepository) UpdateStatus(id string, status string) error {
	query := `UPDATE appointments SET status = ? WHERE id = ?`
	result, err := r.DB.Exec(query, status, id)
	if err != nil {
		return err
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return errors.New("appointment not found")
	}

	return nil
}

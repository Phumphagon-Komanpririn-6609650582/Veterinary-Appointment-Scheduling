package services

import (
	"errors"
	"time"
	"veterinary-api/repositories"
)

type AppointmentService struct {
	Repo *repositories.AppointmentRepository
}

func NewAppointmentService(repo *repositories.AppointmentRepository) *AppointmentService {
	return &AppointmentService{Repo: repo}
}

// =====================================================================
// 👩‍💻 พื้นที่ของ: ปลา (เช็ก Slot Limit ก่อนสร้างนัด)
// =====================================================================
func (s *AppointmentService) CreateAppointment() {

}

// =====================================================================
// 👩‍💻 พื้นที่ของ: นุช (ตรวจสอบความถูกต้องก่อนแก้ข้อมูล)
// =====================================================================
func (s *AppointmentService) UpdateAppointment() {

}

// =====================================================================
// 👨‍💻 พื้นที่ของ: เอลฟ์ (เช็กเงื่อนไข ต้องยกเลิกก่อน 2 ชั่วโมง)
// =====================================================================
func (s *AppointmentService) CancelAppointment(id string) error {
	// 1. ดึงข้อมูลมาเช็คเวลา
	app, err := s.Repo.GetByID(id)
	if err != nil {
		return errors.New("This appointment was not found in the system.")
	}

	// 2. Logic เช็ค 2 ชม.
	diff := time.Until(app.AppointmentTime)
	if diff < 2*time.Hour {
		return errors.New("This cannot be canceled as it must be canceled at least 2 hours in advance.")
	}

	if app.Status == "cancelled" {
		return errors.New("This appointment has been cancelled.")
	}

	// 3. ถ้าผ่านเงื่อนไข ค่อยสั่ง Repo ให้ทำงาน
	return s.Repo.CancelAppointment(id) // เรียกชื่อเดียวกันแต่คนละหน้าที่
}

// =====================================================================
// 👨‍💻 พื้นที่ของ: พี่สิรภพ (กรองรายการนัดหมายตามวัน/ชื่อหมอ)
// =====================================================================
func (s *AppointmentService) GetAppointments() {

}

// =====================================================================
// 👨‍💻 พื้นที่ของ: พี่อิทธิเชษฐ์ (อัปเดตสถานะ done/in-progress/cancelled)
// =====================================================================
func (s *AppointmentService) UpdateStatus() {

}

package services

import (
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
func (s *AppointmentService) CancelAppointment() {

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

package services

import (
	"errors"
	"time"
	"veterinary-api/models"
	"veterinary-api/repositories"
)

type IAppointmentService interface {
	CreateAppointment(app *models.Appointment) error
	UpdateAppointment(app models.Appointment) error
	CancelAppointment(id string) error
	GetAppointments() ([]models.Appointment, error)
	UpdateStatus(id string, status string) error
}

type AppointmentService struct {
	Repo *repositories.AppointmentRepository
}

func NewAppointmentService(repo *repositories.AppointmentRepository) IAppointmentService {
	return &AppointmentService{Repo: repo}
}

// =====================================================================
// 👩‍💻 พื้นที่ของ: ปลา
// =====================================================================
func (s *AppointmentService) CreateAppointment(app *models.Appointment) error {

	// เช็กชื่อสัตว์
	if app.PetName == "" {
		return errors.New("pet name is required")
	}

	// เช็กประเภทสัตว์
	if app.PetType == "" {
		return errors.New("pet type is required")
	}

	// เช็กชื่อเจ้าของ
	if app.ClientName == "" {
		return errors.New("client name is required")
	}

	// เช็กเหตุผลการนัด
	if app.Reason == "" {
		return errors.New("reason is required")
	}

	// เช็ก slot
	if app.SlotID == "" {
		return errors.New("slot id is required")
	}

	// เช็กว่า slot มีจริงไหม
	if !s.Repo.CheckSlotExists(app.SlotID) {
		return errors.New("slot not found")
	}

	// ถ้า slot เต็ม
	err := s.Repo.DecreaseSlotLimit(app.SlotID)
	if err != nil {
		return err
	}

	// ถ้าไม่ส่ง status มา
	if app.Status == "" {
		app.Status = "in-progress"
	}

	// บันทึกลง database
	err = s.Repo.CreateAppointment(app)

	if err != nil {
		return errors.New("failed to create appointment")
	}

	return nil
}

// =====================================================================
// 👩‍💻 พื้นที่ของ: นุช
// =====================================================================
func (s *AppointmentService) UpdateAppointment(app models.Appointment) error {
	oldApp, err := s.Repo.GetByID(app.ID)
	if err != nil {
		return errors.New("appointment not found")
	}

	if oldApp.Status == "cancelled" {
		return errors.New("cannot edit cancelled appointment")
	}

	return s.Repo.UpdateAppointment(app)
}

// =====================================================================
// 👨‍💻 พื้นที่ของ: เอลฟ์
// =====================================================================
func (s *AppointmentService) CancelAppointment(id string) error {
	app, err := s.Repo.GetByID(id)
	if err != nil {
		return errors.New("This appointment was not found in the system.")
	}

	diff := time.Until(app.AppointmentTime)
	if diff < 2*time.Hour {
		return errors.New("This cannot be canceled as it must be canceled at least 2 hours in advance.")
	}

	if app.Status == "cancelled" {
		return errors.New("This appointment has been cancelled.")
	}

	return s.Repo.CancelAppointment(id)
}

// =====================================================================
// 👨‍💻 พื้นที่ของ: พี่สิรภพ
// =====================================================================
func (s *AppointmentService) GetAppointments() ([]models.Appointment, error) {
	// ต้องคืนค่าให้ตรงกับ Interface
	return nil, nil
}

// =====================================================================
// 👨‍💻 พื้นที่ของ: พี่อิทธิเชษฐ์
// =====================================================================
func (s *AppointmentService) UpdateStatus(id string, status string) error {
	// ตรวจสอบสถานะที่อนุญาต
	validStatuses := map[string]bool{
		"done":        true,
		"in-progress": true,
		"cancelled":   true,
	}

	if !validStatuses[status] {
		return errors.New("invalid status: must be 'done', 'in-progress', or 'cancelled'")
	}

	return s.Repo.UpdateStatus(id, status)
}

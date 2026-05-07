package services

import (
	"veterinary-api/models"
	"veterinary-api/repositories"
)

type SlotService struct {
	Repo *repositories.SlotRepository
}

func NewSlotService(repo *repositories.SlotRepository) *SlotService {
	return &SlotService{Repo: repo}
}

// =====================================================================
// 👨‍💻 พื้นที่ของ: ภูมิ (กรองและแสดงเฉพาะเวลาที่ว่าง)
// =====================================================================
func (s *SlotService) GetAvailableSlots(vetID string) ([]models.Slot, error) {
	return s.Repo.GetAvailableSlots(vetID)
}

// =====================================================================
// 👩‍💻 พื้นที่ของ: ไตเติ้ล (ตรวจสอบสิทธิ์ผู้ช่วยก่อนเพิ่มเวลา)
// =====================================================================
func (s *SlotService) AddSlot() {

}

package services

import (
	"veterinary-api/models"
	"veterinary-api/repositories"

	"github.com/google/uuid" //ถ้าขึ้นสีแดงให้ใช้คำสั่ง go get github.com/google/uuid ในcmd
)

type SlotService struct {
	Repo *repositories.SlotRepository
}

func NewSlotService(repo *repositories.SlotRepository) *SlotService {
	return &SlotService{Repo: repo}
}

// =====================================================================
// 👨‍💻 พื้นที่ของ: ภูมิ ( ดูช่วงเวลาว่างของสัตวแพทย์แต่ละท่าน)
// =====================================================================
func (s *SlotService) GetAvailableSlots(vetID string) ([]models.Slot, error) {
	return s.Repo.GetAvailableSlots(vetID)
}

// =====================================================================
// 👩‍💻 พื้นที่ของ: ไตเติ้ล (ตรวจสอบสิทธิ์ผู้ช่วยก่อนเพิ่มเวลา)
// =====================================================================
func (s *SlotService) AddSlot(slot models.Slot) error {
	// ถ้าหน้าบ้านไม่ได้ส่ง ID ของ Slot มา ให้เราสร้างให้ใหม่
	if slot.ID == "" {
		slot.ID = "S-" + uuid.New().String()[:6] // สุ่ม ID เช่น S-a1b2c3
	}

	// ส่งไปบันทึกลงฐานข้อมูล
	return s.Repo.AddSlot(slot)
}

// =====================================================================
// 👨‍💻 พื้นที่ของ: ภูมิ ( ดูช่วงเวลาว่างของสัตวแพทย์ทุกคน)
// =====================================================================
func (s *SlotService) GetAllAvailableSlots() ([]models.Slot, error) {
	return s.Repo.GetAllAvailableSlots()
}

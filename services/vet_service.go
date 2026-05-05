package services

import (
	"veterinary-api/repositories"
)

type VetService struct {
	Repo *repositories.VetRepository
}

func NewVetService(repo *repositories.VetRepository) *VetService {
	return &VetService{Repo: repo}
}

// =====================================================================
// 👨‍💻 พื้นที่ของ: ภูมิ (ดึงข้อมูลหมอทั้งหมดไปทำตัวเลือก)
// =====================================================================
func (s *VetService) GetAllVets() {

}

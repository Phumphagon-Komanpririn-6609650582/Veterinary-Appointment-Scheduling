package services

import (
	"veterinary-api/repositories"
)

type AuthService struct {
	Repo *repositories.AuthRepository
}

func NewAuthService(repo *repositories.AuthRepository) *AuthService {
	return &AuthService{Repo: repo}
}

// =====================================================================
// 👨‍💻 พื้นที่ของ: เอลฟ์ (ตรวจสอบรหัสผ่านและสร้าง JWT Token)
// =====================================================================
func (s *AuthService) Login() {

}

// =====================================================================
// 👨‍💻 พื้นที่ของ: พี่อิทธิเชษฐ์ (จัดการเคลียร์ Session)
// =====================================================================
func (s *AuthService) Logout() {

}

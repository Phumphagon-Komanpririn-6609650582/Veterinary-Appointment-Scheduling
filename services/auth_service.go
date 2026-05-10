package services

import (
	"errors"
	"time"
	"veterinary-api/repositories"

	"github.com/golang-jwt/jwt/v5"
)

type AuthService struct {
	Repo *repositories.AuthRepository
}

func NewAuthService(repo *repositories.AuthRepository) *AuthService {
	return &AuthService{Repo: repo}
}

// กำหนด Key สำหรับเซ็นชื่อ (แนะนำให้ดึงจาก os.Getenv("JWT_SECRET"))
var jwtSecret = []byte("your_super_secret_key_2026")

// =====================================================================
// 👨‍💻 พื้นที่ของ: เอลฟ์ (ตรวจสอบรหัสผ่านและสร้าง JWT Token)
// =====================================================================
func (s *AuthService) Login(username, password string) (string, error) {
	// 1. ค้นหา User จากฐานข้อมูลผ่าน Repo
	user, err := s.Repo.FindByUsername(username)
	if err != nil {
		return "", errors.New("This username was not found in the system.")
	}

	// 2. ตรวจสอบรหัสผ่าน
	// แก้ไข: เนื่องจากใน DB ของเอลฟ์เก็บเป็นค่าธรรมเนียม (123, 456)
	// จึงต้องเทียบกันตรงๆ แบบนี้ไปก่อนครับ
	if user.Password != password {
		return "", errors.New("The password is incorrect.")
	}

	// 3. เมื่อรหัสผ่านถูกต้อง -> สร้างบัตรพนักงาน (JWT Claims)
	claims := jwt.MapClaims{
		"user_id": user.ID,
		"role":    user.Role,
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
		"iat":     time.Now().Unix(),
	}

	// 4. ทำการเซ็นชื่อลงบนบัตร
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString(jwtSecret)
	if err != nil {
		return "", err
	}

	return signedToken, nil
}

// =====================================================================
// 👨‍💻 พื้นที่ของ: พี่อิทธิเชษฐ์ (จัดการเคลียร์ Session)
// =====================================================================
func (s *AuthService) Logout() error {
	return s.Repo.Logout()
}

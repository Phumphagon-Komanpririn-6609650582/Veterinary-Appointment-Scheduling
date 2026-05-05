package repositories

import (
	"database/sql"
)

type AuthRepository struct {
	DB *sql.DB
}

func NewAuthRepository(db *sql.DB) *AuthRepository {
	return &AuthRepository{DB: db}
}

// =====================================================================
// 👨‍💻 พื้นที่ของ: เอลฟ์ (ค้นหา User ตอน Login)
// =====================================================================
func (r *AuthRepository) Login() {

}

// =====================================================================
// 👨‍💻 พื้นที่ของ: พี่อิทธิเชษฐ์ (จัดการเรื่อง Logout เช่น อัปเดตสถานะ Token ใน DB ถ้ามี)
// =====================================================================
func (r *AuthRepository) Logout() {

}

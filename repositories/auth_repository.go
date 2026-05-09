package repositories

import (
	"database/sql"
	"errors"
	"veterinary-api/models"
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
func (r *AuthRepository) FindByUsername(username string) (*models.User, error) {
	var user models.User

	// 1. เปลี่ยนชื่อคอลัมน์จาก password เป็น password_hash ให้ตรงกับในรูป DB
	// 2. เพิ่ม name เข้าไปด้วยเพื่อให้ครบตามโครงสร้างตารางในรูป
	query := "SELECT id, username, password_hash, role, name FROM users WHERE username = ? LIMIT 1"

	err := r.DB.QueryRow(query, username).Scan(
		&user.ID,
		&user.Username,
		&user.Password, // ข้อมูลจาก password_hash จะมาลงที่นี่
		&user.Role,
		&user.Name, // อย่าลืม Scan ชื่อมาด้วยนะ
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("This username was not found.")
		}
		return nil, err
	}

	return &user, nil
}

// =====================================================================
// 👨‍💻 พื้นที่ของ: พี่อิทธิเชษฐ์ (จัดการเรื่อง Logout เช่น อัปเดตสถานะ Token ใน DB ถ้ามี)
// =====================================================================
func (r *AuthRepository) Logout() {

}

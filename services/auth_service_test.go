package services

import (
	"database/sql"
	"testing"
	"veterinary-api/repositories" // import repo เข้ามาด้วย

	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
)

func TestAuthService_Login_Success(t *testing.T) {
	// 1. จำลองการต่อ Database (Path ต้องถอยหลัง 1 ก้าวเพื่อไปหา .db ที่อยู่ root)
	db, _ := sql.Open("sqlite3", "../veterinary.db")
	defer db.Close()

	// 2. สร้าง Repo และส่งเข้าไปใน Service
	repo := repositories.NewAuthRepository(db)
	service := NewAuthService(repo) // ใช้ NewAuthService แทนการสร้าง struct เปล่าๆ

	// 3. ทดสอบ (ตรวจสอบว่าใน DB มี user: admin รหัส: 1234 จริงๆ ไหม)
	token, err := service.Login("ass01", "789")

	assert.NoError(t, err)
	assert.NotEmpty(t, token)
}

func TestAuthService_Login_WrongPassword(t *testing.T) {
	db, _ := sql.Open("sqlite3", "../veterinary.db")
	defer db.Close()

	repo := repositories.NewAuthRepository(db)
	service := NewAuthService(repo)

	_, err := service.Login("ass01", "wrong")

	assert.Error(t, err)
	// ตรงนี้ actual จะกลายเป็น "รหัสผ่านไม่ถูกต้อง" แล้ว Coverage จะพุ่งครับ!
	assert.Equal(t, "The password is incorrect.", err.Error())
}

func TestAuthService_Login_UserNotFound(t *testing.T) {
	// --- 1. ต้องประกาศพวกนี้ใหม่ในทุกฟังก์ชันที่ต้องการใช้ service ---
	db, _ := sql.Open("sqlite3", "../veterinary.db")
	defer db.Close()

	repo := repositories.NewAuthRepository(db)
	service := NewAuthService(repo) // บรรทัดนี้แหละที่จะทำให้ตัวแดงหายไป!
	// -----------------------------------------------------------

	_, err := service.Login("user_67667", "123")

	assert.Error(t, err)
	assert.Equal(t, "This username was not found in the system.", err.Error())
}

package repositories

import (
	"database/sql"
	"testing"

	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
)

func TestAuthRepository_FindByUsername(t *testing.T) {
	// 1. เชื่อมต่อ DB (ถอย 1 ชั้นไปหาไฟล์ .db)
	db, err := sql.Open("sqlite3", "../veterinary.db")
	assert.NoError(t, err)
	defer db.Close()

	repo := NewAuthRepository(db)

	// 2. เคสที่เจอ User (ใช้ข้อมูลจริงจากตารางเช่น ass01)
	user, err := repo.FindByUsername("ass01")
	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, "ass01", user.Username)
	assert.NotEmpty(t, user.Password) // ต้องดึงค่าจาก password_hash มาได้

	// 3. เคสที่ไม่เจอ User
	userNotFound, err := repo.FindByUsername("user_ไม่มีจริง")
	assert.Error(t, err)
	assert.Nil(t, userNotFound)
	assert.Equal(t, "ไม่พบชื่อผู้ใช้งานนี้", err.Error())
}

func TestAuthRepository_DatabaseError(t *testing.T) {
	db, _ := sql.Open("sqlite3", "../veterinary.db")
	repo := NewAuthRepository(db)

	db.Close() // ปิดมันซะ! เพื่อให้ QueryRow พัง

	_, err := repo.FindByUsername("any")
	assert.Error(t, err)
	assert.NotEqual(t, "ไม่พบชื่อผู้ใช้งานนี้", err.Error())
}

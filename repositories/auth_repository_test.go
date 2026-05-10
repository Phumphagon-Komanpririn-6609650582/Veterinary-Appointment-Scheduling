package repositories

import (
	"database/sql"
	"testing"

	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
)

func TestAuthRepository_FindByUsername(t *testing.T) {

	// 1. เชื่อมต่อ DB
	db, err := sql.Open("sqlite3", "../veterinary.db")
	assert.NoError(t, err)
	defer db.Close()

	repo := NewAuthRepository(db)

	// =====================================================
	// เจอ user
	// =====================================================
	user, err := repo.FindByUsername("ass01")

	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, "ass01", user.Username)
	assert.NotEmpty(t, user.Password)

	// =====================================================
	// ไม่เจอ user
	// =====================================================
	userNotFound, err := repo.FindByUsername("user_ไม่มีจริง")

	assert.Error(t, err)
	assert.Nil(t, userNotFound)
	assert.Equal(t, "This username was not found.", err.Error())
}

func TestAuthRepository_DatabaseError(t *testing.T) {

	db, _ := sql.Open("sqlite3", "../veterinary.db")

	repo := NewAuthRepository(db)

	// ปิด db ก่อน
	db.Close()

	_, err := repo.FindByUsername("any")

	assert.Error(t, err)
	assert.Equal(t, "sql: database is closed", err.Error())
}

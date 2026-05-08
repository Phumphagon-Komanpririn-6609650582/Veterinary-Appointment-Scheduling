package controllers

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"veterinary-api/repositories"
	"veterinary-api/services"

	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3" // อย่าลืมตัวนี้เพื่อเปิด driver sqlite
	"github.com/stretchr/testify/assert"
)

func TestLoginController_Success(t *testing.T) {
	// 1. Setup
	gin.SetMode(gin.TestMode)
	r := gin.New()

	// --- 1. เตรียม Database (ดึงไฟล์จากระดับโฟลเดอร์นอก) ---
	db, err := sql.Open("sqlite3", "../veterinary.db")
	if err != nil {
		t.Fatalf("เปิด DB ไม่สำเร็จ: %v", err)
	}
	defer db.Close()

	// --- 2. ประกอบร่างพนักงาน (Dependency Injection) ---
	authRepo := repositories.NewAuthRepository(db)
	authService := services.NewAuthService(authRepo)
	authController := NewAuthController(authService) // สร้าง controller ตัวจริงขึ้นมา

	// --- 3. ลงทะเบียนเส้นทาง (Route) ---
	r.POST("/login", authController.Login)

	// 2. จำลองข้อมูล (ใช้ข้อมูลจริงใน DB)
	loginData := map[string]string{
		"username": "ass01",
		"password": "789",
	}
	body, _ := json.Marshal(loginData)

	// 3. ยิง Request
	req, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	// 4. ตรวจสอบผลลัพธ์
	assert.Equal(t, http.StatusOK, w.Code)

	// (แถม) เช็คว่าได้ Token กลับมาจริงไหม
	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.NotNil(t, response["token"])
}

// 1. เทสกรณีส่งข้อมูลมาไม่ครบ (ทำให้บรรทัด 30-33 เขียว)
func TestLoginController_BadRequest(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	// ... setup controller เหมือนเดิม ...
	db, _ := sql.Open("sqlite3", "../veterinary.db")
	defer db.Close()
	authRepo := repositories.NewAuthRepository(db)
	authService := services.NewAuthService(authRepo)
	authController := NewAuthController(authService)
	// --------------------------------------------------------

	r.POST("/login", authController.Login)

	body, _ := json.Marshal(map[string]string{"username": "ass01"}) // ไม่ส่ง password
	req, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(body))
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, 400, w.Code)
}

// 2. เทสกรณีรหัสผ่านผิด (ทำให้บรรทัด 40-43 เขียว)
func TestLoginController_Unauthorized(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	// ... setup controller เหมือนเดิม ...
	db, _ := sql.Open("sqlite3", "../veterinary.db")
	defer db.Close()
	authRepo := repositories.NewAuthRepository(db)
	authService := services.NewAuthService(authRepo)
	authController := NewAuthController(authService)
	// --------------------------------------------------------

	r.POST("/login", authController.Login)

	body, _ := json.Marshal(map[string]string{"username": "ass01", "password": "wrongpassword"})
	req, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(body))
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, 401, w.Code)
}

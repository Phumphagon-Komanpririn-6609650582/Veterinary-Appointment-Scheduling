package middlewares

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
)

func TestRequireAuth_ValidToken(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()

	// สร้าง Route จำลองที่ใช้ Middleware ของเรา
	r.GET("/test", RequireAuth, func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	// --- ขั้นตอนสำคัญ: สร้าง Token ปลอมที่ "ถูกต้อง" ---
	claims := jwt.MapClaims{
		"user_id": 1,
		"role":    "admin",
		"exp":     time.Now().Add(time.Hour).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// ต้องใช้ jwtSecret ตัวเดียวกับใน auth_middleware.go นะเอล์ฟ!
	validToken, _ := token.SignedString(jwtSecret)

	// จำลอง Request
	req, _ := http.NewRequest("GET", "/test", nil)
	// ต้องมีคำว่า "Bearer " และเว้นวรรค 1 ทีนะ
	req.Header.Set("Authorization", "Bearer "+validToken)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	// ถ้ากุญแจตรงกัน และ Format ถูก ต้องได้ 200 OK
	assert.Equal(t, http.StatusOK, w.Code)
}

// 1. เทสกรณีไม่ส่ง Token มาเลย
func TestRequireAuth_NoToken(t *testing.T) {
	r := gin.New()
	r.GET("/test", RequireAuth, func(c *gin.Context) { c.Status(200) })

	req, _ := http.NewRequest("GET", "/test", nil) // ไม่ใส่ Header Authorization
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, 401, w.Code)
}

// 2. เทสกรณีส่งรูปแบบผิด (ไม่มี Bearer)
func TestRequireAuth_WrongFormat(t *testing.T) {
	r := gin.New()
	r.GET("/test", RequireAuth, func(c *gin.Context) { c.Status(200) })

	req, _ := http.NewRequest("GET", "/test", nil)
	req.Header.Set("Authorization", "JustToken12345") // ผิด Format
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, 401, w.Code)
}

func TestRequireAuth_WrongSigningMethod(t *testing.T) {
	// สร้าง Token ที่ใช้การเซ็นแบบอื่นที่ไม่ใช่ HS256 (เช่น RS256 หรือมั่วๆ)
	token := jwt.New(jwt.SigningMethodNone)
	tokenString, _ := token.SignedString(jwt.UnsafeAllowNoneSignatureType)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("GET", "/", nil)
	c.Request.Header.Set("Authorization", "Bearer "+tokenString)

	RequireAuth(c)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

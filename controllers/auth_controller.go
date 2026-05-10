package controllers

import (
	"net/http"
	"veterinary-api/services"

	"github.com/gin-gonic/gin"
)

type AuthController struct {
	Service *services.AuthService
}

func NewAuthController(service *services.AuthService) *AuthController {
	return &AuthController{Service: service}
}

// =====================================================================
// 👨‍💻 พื้นที่ของ: เอล์ฟ (POST /api/auth/login)
// =====================================================================
func (c *AuthController) Login(ctx *gin.Context) {
	var loginReq struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	// 1. ตรวจสอบว่าส่งข้อมูลมาครบไหม
	if err := ctx.ShouldBindJSON(&loginReq); err != nil {
		ctx.JSON(400, gin.H{"error": "Please enter your Username and Password."})
		return
	}

	// 2. เรียก Service เพื่อยืนยันตัวตนและขอ Token
	// (สมมติว่าใน AuthService มีฟังก์ชัน Login ที่คืนค่า token และ error)
	token, err := c.Service.Login(loginReq.Username, loginReq.Password)
	if err != nil {
		// ถ้ารหัสผิด หรือหา user ไม่เจอ ให้ตอบกลับแบบกั๊กๆ ไว้เพื่อความปลอดภัย
		ctx.JSON(401, gin.H{"error": "The username or password is incorrect."})
		return
	}

	// 3. ส่ง "บัตรพนักงาน" กลับไปให้ User
	ctx.JSON(200, gin.H{
		"message": "Login successful",
		"token":   token,
	})
}

// =====================================================================
// 👨‍💻 พื้นที่ของ: พี่อิทธิเชษฐ์ (POST /api/auth/logout)
// =====================================================================
func (c *AuthController) Logout(ctx *gin.Context) {
	// เรียก Service เพื่อจัดการเรื่องการ Logout (เช่น เคลียร์ Session หรือ Token ใน DataBase)
	err := c.Service.Logout()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to logout"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"message": "Logout successful. Please clear your token.",
	})
}

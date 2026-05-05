package controllers

import (
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

}

// =====================================================================
// 👨‍💻 พื้นที่ของ: พี่อิทธิเชษฐ์ (POST /api/auth/logout)
// =====================================================================
func (c *AuthController) Logout(ctx *gin.Context) {

}

package middlewares

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// ใช้ Secret Key เดียวกันกับใน AuthService
var jwtSecret = []byte("cs367secretkey")

func RequireAuth(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "ไม่มีสิทธิ์เข้าถึง: กรุณาส่ง Token มาด้วย"})
		c.Abort()
		return
	}

	// 2. ตรวจสอบรูปแบบ "Bearer <token>"
	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "รูปแบบ Token ไม่ถูกต้อง (ต้องเป็น Bearer <token>)"})
		c.Abort()
		return
	}

	tokenString := parts[1]

	// 3. แกะและตรวจสอบความถูกต้องของ Token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// ตรวจสอบว่าใช้วิธีการ Signing แบบที่คาดไว้ (HS256) ไหม
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return jwtSecret, nil
	})

	// 4. ถ้า Token ปลอมหรือหมดอายุ
	if err != nil || !token.Valid {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "บัตรพนักงานปลอมหรือหมดอายุแล้ว"})
		c.Abort()
		return
	}

	// 5. ดึงข้อมูลใน Payload (Claims) ออกมาแปะไว้ใน Context
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		c.Set("user_id", claims["user_id"])
		c.Set("role", claims["role"])
	}

	// ผ่านด่าน! ให้ไปทำ Logic ใน Controller ต่อได้
	c.Next()
}

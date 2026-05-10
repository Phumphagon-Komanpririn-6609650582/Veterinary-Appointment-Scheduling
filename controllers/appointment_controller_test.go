package controllers

import (
	"database/sql"
	"net/http"
	"net/http/httptest"
	"testing"
	"veterinary-api/repositories"
	"veterinary-api/services"

	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3" // อย่าลืมตัวนี้เพื่อเปิด driver sqlite
	"github.com/stretchr/testify/assert"
)

func TestAppointmentController_CancelAppointment(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db, _ := sql.Open("sqlite3", "../veterinary.db")
	repo := repositories.NewAppointmentRepository(db)
	service := services.NewAppointmentService(repo)
	ctrl := NewAppointmentController(service)

	t.Run("Success_Response", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		// จำลองการส่ง ID ผ่าน Param
		c.Params = []gin.Param{{Key: "id", Value: "A-001"}}

		ctrl.CancelAppointment(c)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "สำเร็จ")
	})
}

// เพิ่มใน controllers/appointment_controller_test.go

func TestAppointmentController_CancelAppointment_AllCases(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db, _ := sql.Open("sqlite3", "../veterinary.db")
	repo := repositories.NewAppointmentRepository(db)
	service := services.NewAppointmentService(repo)
	ctrl := NewAppointmentController(service)

	// --- เคสที่ 1: ยกเลิกสำเร็จ (HTTP 200) ---
	t.Run("Success_Response", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		// มั่นใจว่า ID A-001 ใน DB คือ "ปีหน้า" และ "ยังไม่โดนยกเลิก"
		c.Params = []gin.Param{{Key: "id", Value: "A-001"}}

		ctrl.CancelAppointment(c)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "The appointment has been cancelled.")
	})

	// --- เคสที่ 2: ยกเลิกไม่ทัน/ไม่พบข้อมูล (HTTP 400) ---
	t.Run("BadRequest_Error_Response", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		// ใช้ ID ที่เวลากระชั้นชิด หรือ ID ที่ไม่มีจริง
		c.Params = []gin.Param{{Key: "id", Value: "A-002"}}

		ctrl.CancelAppointment(c)

		// เมื่อ Service คืนค่า error ตัว Controller จะต้องพ่น 400 Bad Request
		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "error")
	})
}

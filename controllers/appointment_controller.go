package controllers

import (
	"net/http"
	"veterinary-api/services"

	"github.com/gin-gonic/gin"
)

type AppointmentController struct {
	Service *services.AppointmentService
}

func NewAppointmentController(service *services.AppointmentService) *AppointmentController {
	return &AppointmentController{Service: service}
}

// =====================================================================
// 👩‍💻 พื้นที่ของ: ปลา (POST /api/appointments)
// =====================================================================
func (c *AppointmentController) CreateAppointment(ctx *gin.Context) {

}

// =====================================================================
// 👩‍💻 พื้นที่ของ: นุช (PUT /api/appointments/:id)
// =====================================================================
func (c *AppointmentController) UpdateAppointment(ctx *gin.Context) {

}

// =====================================================================
// 👨‍💻 พื้นที่ของ: เอลฟ์ (DELETE /api/appointments/:id)
// =====================================================================
func (ctrl *AppointmentController) CancelAppointment(c *gin.Context) {
	// 1. ดึง ID จาก URL (เช่น /api/appointments/A-001/cancel)
	id := c.Param("id")

	// 2. ส่ง ID ไปให้ Service จัดการ Logic (เช็ค 2 ชม.)
	err := ctrl.Service.CancelAppointment(id)

	// 3. ถ้า Service บอกว่ามี Error (เช่น ไม่ถึง 2 ชม. หรือหาไม่เจอ)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	// 4. ถ้าทุกอย่างผ่านฉลุย
	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "The appointment has been cancelled.",
	})
}

// =====================================================================
// 👨‍💻 พื้นที่ของ: พี่สิรภพ (GET /api/appointments)
// =====================================================================
func (c *AppointmentController) GetAppointments(ctx *gin.Context) {

}

// =====================================================================
// 👨‍💻 พื้นที่ของ: พี่อิทธิเชษฐ์ (PATCH /api/appointments/:id/status)
// =====================================================================
func (c *AppointmentController) UpdateStatus(ctx *gin.Context) {

}

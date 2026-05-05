package controllers

import (
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
func (c *AppointmentController) DeleteAppointment(ctx *gin.Context) {

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

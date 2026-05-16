package controllers

import (
	"net/http"
	"veterinary-api/models"
	"veterinary-api/services"

	"github.com/gin-gonic/gin"
)

type AppointmentController struct {
	// เปลี่ยนจากเจาะจง Struct เป็นการใช้ Interface เพื่อให้ส่ง Mock เข้ามาเทสได้
	Service services.IAppointmentService
}

func NewAppointmentController(service services.IAppointmentService) *AppointmentController {
	return &AppointmentController{Service: service}
}

// =====================================================================
// 👩‍💻 พื้นที่ของ: ปลา (POST /api/appointments)
// =====================================================================
func (c *AppointmentController) CreateAppointment(ctx *gin.Context) {

	var appointment models.Appointment

	// รับ JSON
	if err := ctx.ShouldBindJSON(&appointment); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "invalid request data",
		})
		return
	}

	// เรียก service
	err := c.Service.CreateAppointment(&appointment)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	// success
	ctx.JSON(http.StatusCreated, gin.H{
		"status":  "success",
		"message": "appointment created successfully",
		"data":    appointment,
	})
}

// =====================================================================
// 👩‍💻 พื้นที่ของ: นุช (PUT /api/appointments/:id)
// =====================================================================
func (c *AppointmentController) UpdateAppointment(ctx *gin.Context) {
	id := ctx.Param("id")
	var app models.Appointment

	if err := ctx.ShouldBindJSON(&app); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Failed to update appointment",
		})
		return
	}

	app.ID = id
	err := c.Service.UpdateAppointment(app)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": err.Error(), // เปลี่ยนเป็น err.Error() เพื่อให้เห็น Error จาก Mock/Service
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "appointment updated successfully",
	})
}

// =====================================================================
// 👨‍💻 พื้นที่ของ: เอลฟ์ (DELETE /api/appointments/:id)
// =====================================================================
func (ctrl *AppointmentController) CancelAppointment(c *gin.Context) {
	id := c.Param("id")

	err := ctrl.Service.CancelAppointment(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "The appointment has been cancelled.",
	})
}

// (GET /api/appointments)
func (c *AppointmentController) GetAppointments(ctx *gin.Context) {
	date := ctx.Query("date")
	vetID := ctx.Query("vet_id")

	appointments, err := c.Service.GetAppointments(date, vetID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "failed to retrieve appointments",
		})
		return
	}

	if appointments == nil {
		appointments = []models.Appointment{}
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   appointments,
	})
}

// =====================================================================
// 👨‍💻 พื้นที่ของ: พี่อิทธิเชษฐ์ (PATCH /api/appointments/:id/status)
// =====================================================================
func (c *AppointmentController) UpdateStatus(ctx *gin.Context) {
	id := ctx.Param("id")
	var statusReq struct {
		Status string `json:"status" binding:"required"`
	}

	if err := ctx.ShouldBindJSON(&statusReq); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "invalid request data: status is required",
		})
		return
	}

	err := c.Service.UpdateStatus(id, statusReq.Status)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "appointment status updated successfully",
	})
}

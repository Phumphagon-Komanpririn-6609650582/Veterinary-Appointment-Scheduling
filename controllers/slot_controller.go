package controllers

import (
	"net/http"
	"veterinary-api/models"
	"veterinary-api/repositories"

	"github.com/gin-gonic/gin"
)

type ISlotService interface {
	GetAvailableSlots(vetID string) ([]models.Slot, error)
	GetAllAvailableSlots() ([]models.Slot, error)
	AddSlot(slot models.Slot) error
}

type SlotController struct {
	Service ISlotService
}

func NewSlotController(service ISlotService) *SlotController {
	return &SlotController{Service: service}
}

func respondWithError(c *gin.Context, code int, message string) {
	c.JSON(code, gin.H{"error": message})
}

// =====================================================================
// 👨‍💻 พื้นที่ของ: ภูมิ (GET /api/vets/:id/slots)
// =====================================================================
func (c *SlotController) GetAvailableSlots(ctx *gin.Context) {
	vetID := ctx.Param("id")

	slots, err := c.Service.GetAvailableSlots(vetID)
	if err != nil {

		if err == repositories.ErrNotFound {
			respondWithError(ctx, http.StatusNotFound, "Vet or slots not found")
			return
		}

		respondWithError(ctx, http.StatusInternalServerError, "Failed to retrieve available slots")
		return
	}

	ctx.JSON(http.StatusOK, slots)
}

// =====================================================================
// 👩‍💻 พื้นที่ของ: ไตเติ้ล (POST /api/vets/:id/slots)
// =====================================================================
func (c *SlotController) AddSlot(ctx *gin.Context) {
	// 1. เช็คสิทธิ์ว่าเป็นผู้ช่วย (assistant) หรือไม่
	role, exists := ctx.Get("role")
	if !exists || role != "assistant" {
		respondWithError(ctx, http.StatusForbidden, "Forbidden: Only assistants can add slots")
		return
	}

	// 2. รับข้อมูล JSON ที่ส่งมา (เช่น วันที่, เวลา, จำนวนคิว)
	var newSlot models.Slot
	if err := ctx.ShouldBindJSON(&newSlot); err != nil {
		respondWithError(ctx, http.StatusBadRequest, "Invalid request body")
		return
	}

	// 3. ดึงรหัสหมอ (VetID) จาก URL param (เช่น /api/vets/U-001/slots -> ได้ "U-001")
	newSlot.VetID = ctx.Param("id")

	// 4. ส่งข้อมูลให้ Service ไปบันทึกลงฐานข้อมูล
	if err := c.Service.AddSlot(newSlot); err != nil {
		respondWithError(ctx, http.StatusInternalServerError, "Failed to create slot")
		return
	}

	// 5. คืนค่าผลลัพธ์หน้าบ้านว่าเพิ่มสำเร็จ!
	ctx.JSON(http.StatusCreated, gin.H{
		"message": "Slot created successfully",
		"data":    newSlot,
	})
}

// =====================================================================
// 👨‍💻 พื้นที่ของ: ภูมิ (GET /api/slots)
// =====================================================================
func (c *SlotController) GetAllAvailableSlots(ctx *gin.Context) {
	// 📌 เพิ่มส่วนนี้: เช็คสิทธิ์ว่าเป็นผู้ช่วย (assistant) หรือไม่
	role, exists := ctx.Get("role")
	if !exists || role != "assistant" {
		respondWithError(ctx, http.StatusForbidden, "Forbidden: Only assistants can view all slots")
		return
	}
	// 📌 สิ้นสุดส่วนที่เพิ่ม

	slots, err := c.Service.GetAllAvailableSlots()
	if err != nil {
		if err == repositories.ErrNotFound {
			respondWithError(ctx, http.StatusNotFound, "No available slots found")
			return
		}
		respondWithError(ctx, http.StatusInternalServerError, "Failed to retrieve all available slots")
		return
	}

	ctx.JSON(http.StatusOK, slots)
}

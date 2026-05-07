package controllers

import (
	"net/http"
	"veterinary-api/models" // 🔥 ต้อง import models เพราะเราจะทำ Interface
	"veterinary-api/repositories"

	"github.com/gin-gonic/gin"
)

// =====================================================================
// 🔥 ทริค Tech Lead: สร้าง Interface ฝั่ง Controller เองเลย!
// =====================================================================
type ISlotService interface {
	GetAvailableSlots(vetID string) ([]models.Slot, error)
	GetAllAvailableSlots() ([]models.Slot, error)
}

type SlotController struct {
	Service ISlotService // 🔥 เปลี่ยนมารับ Interface ของเราเอง
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

}

// =====================================================================
// 👨‍💻 พื้นที่ของ: ภูมิ (GET /api/vets/:id/slots)
// =====================================================================
func (c *SlotController) GetAllAvailableSlots(ctx *gin.Context) {
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

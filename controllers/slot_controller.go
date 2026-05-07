package controllers

import (
	"veterinary-api/services"

	"github.com/gin-gonic/gin"
)

type SlotController struct {
	Service *services.SlotService
}

func NewSlotController(service *services.SlotService) *SlotController {
	return &SlotController{Service: service}
}

// =====================================================================
// 👨‍💻 พื้นที่ของ: ภูมิ (GET /api/vets/:id/slots)
// =====================================================================
func (c *SlotController) GetAvailableSlots(ctx *gin.Context) {
	vetID := ctx.Param("id")

	slots, err := c.Service.GetAvailableSlots(vetID)
	if err != nil {
		ctx.JSON(500, gin.H{"error": "Failed to fetch available slots"})
		return
	}

	ctx.JSON(200, gin.H{
		"message": "Success",
		"data":    slots,
	})
}

// =====================================================================
// 👩‍💻 พื้นที่ของ: ไตเติ้ล (POST /api/vets/:id/slots)
// =====================================================================
func (c *SlotController) AddSlot(ctx *gin.Context) {

}

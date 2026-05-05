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

}

// =====================================================================
// 👩‍💻 พื้นที่ของ: ไตเติ้ล (POST /api/vets/:id/slots)
// =====================================================================
func (c *SlotController) AddSlot(ctx *gin.Context) {

}

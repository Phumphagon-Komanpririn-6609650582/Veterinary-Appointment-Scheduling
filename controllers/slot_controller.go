package controllers

import (
	"net/http"
	"veterinary-api/repositories"
	"veterinary-api/services"

	"github.com/gin-gonic/gin"
)

type SlotController struct {
	Service *services.SlotService
}

func NewSlotController(service *services.SlotService) *SlotController {
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
	// (เว้นไว้ให้ไตเติ้ลทำต่อ)
}

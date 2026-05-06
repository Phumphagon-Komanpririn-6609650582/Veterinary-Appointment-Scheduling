package controllers

import (
	"veterinary-api/services"

	"github.com/gin-gonic/gin"
)

type VetController struct {
	Service *services.VetService
}

func NewVetController(service *services.VetService) *VetController {
	return &VetController{Service: service}
}

// =====================================================================
// 👨‍💻 พื้นที่ของ: ภูมิ (GET /api/vets)
// =====================================================================
func (c *VetController) GetAllVets(ctx *gin.Context) {
	vets, err := c.Service.GetAllVets()
	if err != nil {
		ctx.JSON(500, gin.H{"error": "Failed to fetch veterinarians"})
		return
	}

	ctx.JSON(200, gin.H{
		"message": "Success",
		"data":    vets,
	})
}

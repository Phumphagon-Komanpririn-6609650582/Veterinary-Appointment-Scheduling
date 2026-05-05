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

}

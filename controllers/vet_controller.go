package controllers

import (
	"net/http"
	"veterinary-api/repositories"
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

		if err == repositories.ErrNotFound {
			respondWithError(ctx, http.StatusNotFound, "Veterinarians not found")
			return
		}

		respondWithError(ctx, http.StatusInternalServerError, "Failed to retrieve veterinarians")
		return
	}

	ctx.JSON(http.StatusOK, vets)
}

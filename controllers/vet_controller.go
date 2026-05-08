package controllers

import (
	"errors"
	"net/http"
	"veterinary-api/models"
	"veterinary-api/repositories"

	"github.com/gin-gonic/gin"
)

type IVetService interface {
	GetAllVets() ([]models.User, error)
}

type VetController struct {
	Service IVetService
}

func NewVetController(service IVetService) *VetController {
	return &VetController{Service: service}
}

// =====================================================================
// 👨‍💻 พื้นที่ของ: ภูมิ (GET /api/vets)
// =====================================================================
func (c *VetController) GetAllVets(ctx *gin.Context) {
	vets, err := c.Service.GetAllVets()
	if err != nil {

		if errors.Is(err, repositories.ErrNotFound) {
			respondWithError(ctx, http.StatusNotFound, "Veterinarians not found")
			return
		}

		respondWithError(ctx, http.StatusInternalServerError, "Failed to retrieve veterinarians")
		return
	}

	ctx.JSON(http.StatusOK, vets)
}

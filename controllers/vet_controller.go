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
		// เช็กกรณีหาข้อมูลไม่เจอ (สมมติว่าไม่มีหมอในระบบเลย)
		if err == repositories.ErrNotFound {
			respondWithError(ctx, http.StatusNotFound, "Veterinarians not found")
			return
		}

		// กรณีพังจาก Database หรือระบบภายใน ใช้ Error กลางๆ
		respondWithError(ctx, http.StatusInternalServerError, "Failed to retrieve veterinarians")
		return
	}

	// ตอบกลับสำเร็จด้วย HTTP 200 OK และส่งข้อมูล array ออกไปตรงๆ (ไม่หุ้ม message/data)
	ctx.JSON(http.StatusOK, vets)
}

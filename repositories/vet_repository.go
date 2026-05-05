package repositories

import (
	"database/sql"
)

type VetRepository struct {
	DB *sql.DB
}

func NewVetRepository(db *sql.DB) *VetRepository {
	return &VetRepository{DB: db}
}

// =====================================================================
// 👨‍💻 พื้นที่ของ: ภูมิ (ค้นหารายชื่อสัตวแพทย์)
// =====================================================================
func (r *VetRepository) GetAllVets() {

}

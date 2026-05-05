package repositories

import (
	"database/sql"
)

type SlotRepository struct {
	DB *sql.DB
}

func NewSlotRepository(db *sql.DB) *SlotRepository {
	return &SlotRepository{DB: db}
}

// =====================================================================
// 👨‍💻 พื้นที่ของ: ภูมิ (ค้นหาเวลาว่าง)
// =====================================================================
func (r *SlotRepository) GetAvailableSlots() {

}

// =====================================================================
// 👩‍💻 พื้นที่ของ: ไตเติ้ล (Insert เวลาทำงานใหม่)
// =====================================================================
func (r *SlotRepository) AddSlot() {

}

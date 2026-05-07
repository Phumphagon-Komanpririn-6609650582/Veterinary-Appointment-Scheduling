package repositories

import (
	"database/sql"
	"veterinary-api/models"
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
func (r *SlotRepository) GetAvailableSlots(vetID string) ([]models.Slot, error) {
	// ค้นหาเวลาของหมอคนนั้น (vet_id) และเช็กว่ายังมีคิวว่างเหลืออยู่ (slot_limit > 0)
	rows, err := r.DB.Query("SELECT id, vet_id, date, time_period, slot_limit FROM slots WHERE vet_id = ? AND slot_limit > 0", vetID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var slots []models.Slot
	for rows.Next() {
		var slot models.Slot
		// สแกนข้อมูลใส่ Struct ให้ตรงกับที่ออกแบบไว้ใน Slot Model
		if err := rows.Scan(&slot.ID, &slot.VetID, &slot.Date, &slot.TimePeriod, &slot.SlotLimit); err != nil {
			return nil, err
		}
		slots = append(slots, slot)
	}
	return slots, nil
}

// =====================================================================
// 👩‍💻 พื้นที่ของ: ไตเติ้ล (Insert เวลาทำงานใหม่)
// =====================================================================
func (r *SlotRepository) AddSlot() {

}

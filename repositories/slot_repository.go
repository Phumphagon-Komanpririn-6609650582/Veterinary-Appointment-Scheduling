package repositories

import (
	"database/sql"
	"errors"
	"veterinary-api/models"
)

// ประกาศตัวแปร ErrNotFound เพื่อให้ Controller นำไปใช้เช็กและตอบ Status 404 ได้
var ErrNotFound = errors.New("record not found")

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
	query := `
		SELECT slots.id, slots.vet_id, users.name, slots.date, slots.time_period, slots.slot_limit 
		FROM slots 
		JOIN users ON slots.vet_id = users.id 
		WHERE slots.vet_id = ? AND slots.slot_limit > 0
	`
	rows, err := r.DB.Query(query, vetID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var slots []models.Slot
	for rows.Next() {
		var slot models.Slot
		if err := rows.Scan(&slot.ID, &slot.VetID, &slot.VetName, &slot.Date, &slot.TimePeriod, &slot.SlotLimit); err != nil {
			return nil, err
		}
		slots = append(slots, slot)
	}

	// ถ้า Query สำเร็จ แต่ไม่มีข้อมูลกลับมาเลย (array ว่าง) ให้ส่ง ErrNotFound ออกไป
	if len(slots) == 0 {
		return nil, ErrNotFound
	}

	return slots, nil
}

// =====================================================================
// 👩‍💻 พื้นที่ของ: ไตเติ้ล (Insert เวลาทำงานใหม่)
// =====================================================================
func (r *SlotRepository) AddSlot() {

}

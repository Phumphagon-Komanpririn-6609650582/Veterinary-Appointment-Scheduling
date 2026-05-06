package repositories

import (
	"database/sql"
	"veterinary-api/models"
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
func (r *VetRepository) GetAllVets() ([]models.User, error) {
	rows, err := r.DB.Query("SELECT id, username, name, role FROM users WHERE role = 'vet'")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var vets []models.User
	for rows.Next() {
		var vet models.User
		if err := rows.Scan(&vet.ID, &vet.Username, &vet.Name, &vet.Role); err != nil {
			return nil, err
		}
		vets = append(vets, vet)
	}
	return vets, nil
}

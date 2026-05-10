package services

import (
	"database/sql"
	"testing"
	"veterinary-api/repositories" // import repo เข้ามาด้วย

	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
)

func TestAppointmentService_CancelAppointment_Complete(t *testing.T) {
	db, _ := sql.Open("sqlite3", "../veterinary.db")
	repo := repositories.NewAppointmentRepository(db)
	service := NewAppointmentService(repo)

	// 1. เคส: ยกเลิกไม่ทัน (น้อยกว่า 2 ชม.)
	t.Run("Cancel_TooLate_Error", func(t *testing.T) {
		err := service.CancelAppointment("A-002") // มั่นใจว่า A-002 ใน DB เวลากระชั้นชิด
		assert.Error(t, err)
		// แก้ภาษาให้ตรงกับโค้ด Service ของเอลฟ์
		assert.Contains(t, err.Error(), "at least 2 hours in advance")
	})

	// 2. เคส: ยกเลิกสำเร็จ (ต้องหา ID ที่เวลาห่างๆ และสถานะยังไม่เป็น cancelled)
	t.Run("Cancel_Success", func(t *testing.T) {
		// แนะนำ: ไปแก้ A-003 ใน DB ให้เป็นปีหน้า และสถานะเป็น 'pending'
		err := service.CancelAppointment("A-003")
		assert.NoError(t, err)
	})

	// 3. เคส: หา ID ไม่เจอ
	t.Run("Error_NotFound", func(t *testing.T) {
		err := service.CancelAppointment("NON-EXISTENT-ID")
		assert.Error(t, err)
		assert.Equal(t, "This appointment was not found in the system.", err.Error())
	})

	// 4. เคส: ยกเลิกซ้ำ
	t.Run("Error_AlreadyCancelled", func(t *testing.T) {
		// ใช้ ID ที่เรารู้แน่ๆ ว่าเป็น cancelled (เช่น A-001 ที่ยกเลิกไปแล้ว)
		err := service.CancelAppointment("A-001")
		assert.Error(t, err)
		assert.Equal(t, "This appointment has been cancelled.", err.Error())
	})
}

// เพิ่มใน services/appointment_service_test.go

func TestAppointmentService_CancelAppointment_EdgeCases(t *testing.T) {
	db, _ := sql.Open("sqlite3", "../veterinary.db")
	repo := repositories.NewAppointmentRepository(db)
	service := NewAppointmentService(repo)

	t.Run("Error_AppointmentNotFound", func(t *testing.T) {
		// ส่ง ID มั่วๆ ที่ไม่มีใน DB
		err := service.CancelAppointment("NON-EXISTENT-ID")
		assert.Error(t, err)
		assert.Equal(t, "This appointment was not found in the system.", err.Error())
	})

	t.Run("Error_AlreadyCancelled", func(t *testing.T) {
		// ต้องมั่นใจว่าใน DB มี ID นี้และ status เป็น 'cancelled' อยู่แล้ว
		// สมมติเป็น ID A-001 ที่เอลฟ์เพิ่งรัน Test Success ไปก่อนหน้า
		err := service.CancelAppointment("A-001")

		// ถ้า Logic เช็ค status ของเอลฟ์ทำงาน บรรทัดสีแดงจะเขียวทันที
		assert.Error(t, err)
		assert.Equal(t, "This appointment has been cancelled.", err.Error())
	})
}

func TestAppointmentService_Cancel_Errors(t *testing.T) {
	db, _ := sql.Open("sqlite3", "../veterinary.db")
	repo := repositories.NewAppointmentRepository(db)
	service := NewAppointmentService(repo)

	t.Run("Error_NotFound", func(t *testing.T) {
		// 1. ส่ง ID ที่ไม่มีใน DB เพื่อให้ GetByID คืนค่า error (จะทำให้บรรทัดแดงแรกเขียว)
		err := service.CancelAppointment("NON-EXIST-ID")
		assert.Error(t, err)
		assert.Equal(t, "This appointment was not found in the system.", err.Error())
	})

	t.Run("Error_AlreadyCancelled", func(t *testing.T) {
		// 2. ส่ง ID ที่สถานะเป็น 'cancelled' อยู่แล้ว (จะทำให้บรรทัดแดงที่สองเขียว)
		// **ต้องมั่นใจว่าใน DB ของเอลฟ์มี ID A-001 ที่ status='cancelled' แล้วนะ
		err := service.CancelAppointment("A-001")
		assert.Error(t, err)
		assert.Equal(t, "This appointment has been cancelled.", err.Error())
	})
}

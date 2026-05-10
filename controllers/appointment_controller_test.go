package controllers

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"veterinary-api/models"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// --- ส่วนที่ 1: สร้างหุ่นจำลอง (Mock Object) ---
type MockAppointmentService struct {
	mock.Mock
}

func (m *MockAppointmentService) CreateAppointment(app *models.Appointment) error {
	args := m.Called(app)
	return args.Error(0)
}

func (m *MockAppointmentService) UpdateAppointment(app models.Appointment) error {
	args := m.Called(app)
	return args.Error(0)
}

func (m *MockAppointmentService) CancelAppointment(id string) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockAppointmentService) GetAppointments() ([]models.Appointment, error) {
	args := m.Called()
	return args.Get(0).([]models.Appointment), args.Error(1)
}

func (m *MockAppointmentService) UpdateStatus(id string, status string) error {
	args := m.Called(id, status)
	return args.Error(0)
}

// =====================================================================
// 👨‍💻 พื้นที่ของ: เอลฟ์ (DELETE /api/appointments/:id)
// =====================================================================
func TestAppointmentController_CancelAppointment_Mock(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Success_Case", func(t *testing.T) {
		mockService := new(MockAppointmentService)
		ctrl := NewAppointmentController(mockService)
		mockService.On("CancelAppointment", "A-001").Return(nil)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = []gin.Param{{Key: "id", Value: "A-001"}}

		ctrl.CancelAppointment(c)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "The appointment has been cancelled.")
		mockService.AssertExpectations(t)
	})

	t.Run("Fail_2Hours_Constraint", func(t *testing.T) {
		mockService := new(MockAppointmentService)
		ctrl := NewAppointmentController(mockService)
		errMsg := "This cannot be canceled as it must be canceled at least 2 hours in advance."
		mockService.On("CancelAppointment", "A-002").Return(errors.New(errMsg))

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = []gin.Param{{Key: "id", Value: "A-002"}}

		ctrl.CancelAppointment(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), errMsg)
	})
}

// =====================================================================
// 👩‍💻 พื้นที่ของ: นุช (PUT /api/appointments/:id)
// =====================================================================
func TestAppointmentController_UpdateAppointment_Mock(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// เคสที่ 1: อัปเดตสำเร็จ
	t.Run("Update_Success", func(t *testing.T) {
		mockService := new(MockAppointmentService)
		ctrl := NewAppointmentController(mockService)
		mockService.On("UpdateAppointment", mock.MatchedBy(func(app models.Appointment) bool {
			return app.ID == "A-001"
		})).Return(nil)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = []gin.Param{{Key: "id", Value: "A-001"}}
		jsonBody := `{"pet_name": "น้องแมวโชคดี", "note": "มาตามนัดใหม่"}`
		c.Request, _ = http.NewRequest("PUT", "/api/appointments/A-001", strings.NewReader(jsonBody))
		c.Request.Header.Set("Content-Type", "application/json")

		ctrl.UpdateAppointment(c)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "appointment updated successfully")
		mockService.AssertExpectations(t)
	})

	// เคสที่ 2: อัปเดตไม่สำเร็จเพราะ Service พ่น Error (เช่น โดนยกเลิกไปแล้ว)
	t.Run("Update_Fail_Cancelled_Appointment", func(t *testing.T) {
		mockService := new(MockAppointmentService)
		ctrl := NewAppointmentController(mockService)
		errMsg := "cannot edit cancelled appointment"
		mockService.On("UpdateAppointment", mock.Anything).Return(errors.New(errMsg))

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = []gin.Param{{Key: "id", Value: "A-001"}}
		jsonBody := `{"pet_name": "น้อนนน"}`
		c.Request, _ = http.NewRequest("PUT", "/api/appointments/A-001", strings.NewReader(jsonBody))
		c.Request.Header.Set("Content-Type", "application/json")

		ctrl.UpdateAppointment(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), errMsg)
	})

	// เคสที่ 3: ส่ง JSON ผิดรูปแบบ (คลุมโค้ดบรรทัด ShouldBindJSON ของนุช)
	t.Run("Update_Fail_Invalid_JSON", func(t *testing.T) {
		mockService := new(MockAppointmentService)
		ctrl := NewAppointmentController(mockService)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = []gin.Param{{Key: "id", Value: "A-001"}}
		invalidJson := `{"pet_name": "น้อนนน"` // จงใจไม่ใส่ปีกกาปิด
		c.Request, _ = http.NewRequest("PUT", "/api/appointments/A-001", strings.NewReader(invalidJson))
		c.Request.Header.Set("Content-Type", "application/json")

		ctrl.UpdateAppointment(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "Failed to update appointment")
	})
}
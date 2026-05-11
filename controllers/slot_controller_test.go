package controllers

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"veterinary-api/models"
	"veterinary-api/repositories"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockSlotService struct {
	mock.Mock
}

func (m *MockSlotService) GetAllAvailableSlots() ([]models.Slot, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.Slot), args.Error(1)
}

func (m *MockSlotService) GetAvailableSlots(vetID string) ([]models.Slot, error) {
	args := m.Called(vetID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.Slot), args.Error(1)
}

func (m *MockSlotService) AddSlot(slot models.Slot) error {
	args := m.Called(slot)
	return args.Error(0)
}

// =====================================================================
// Test: GetAllAvailableSlots Controller
// =====================================================================

// 200 OK
func TestGetAllAvailableSlots_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockService := new(MockSlotService)
	controller := NewSlotController(mockService)

	mockSlots := []models.Slot{
		{ID: "S-001", VetID: "U-001", VetName: "นสพ.สมชาย", Date: "2026-05-01", TimePeriod: "09:00-10:00", SlotLimit: 1},
	}
	mockService.On("GetAllAvailableSlots").Return(mockSlots, nil)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Set("role", "assistant")

	controller.GetAllAvailableSlots(c)

	assert.Equal(t, http.StatusOK, w.Code)
	mockService.AssertExpectations(t)
}

//404 Not Found

func TestGetAllAvailableSlots_NotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockService := new(MockSlotService)
	controller := NewSlotController(mockService)

	mockService.On("GetAllAvailableSlots").Return(nil, repositories.ErrNotFound)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Set("role", "assistant")

	controller.GetAllAvailableSlots(c)

	assert.Equal(t, http.StatusNotFound, w.Code)
	mockService.AssertExpectations(t)
}

// 500 Internal Server Error
func TestGetAllAvailableSlots_InternalError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockService := new(MockSlotService)
	controller := NewSlotController(mockService)

	mockService.On("GetAllAvailableSlots").Return(nil, errors.New("database connection lost"))

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Set("role", "assistant")

	controller.GetAllAvailableSlots(c)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	mockService.AssertExpectations(t)
}

// =====================================================================
// Test: GetAvailableSlots Controller
// =====================================================================

// 200 OK
func TestGetAvailableSlots_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockService := new(MockSlotService)
	controller := NewSlotController(mockService)

	vetID := "V-001"
	mockSlots := []models.Slot{
		{ID: "S-001", VetID: vetID, VetName: "นสพ.สมชาย", Date: "2026-05-01", TimePeriod: "09:00-10:00", SlotLimit: 1},
	}

	mockService.On("GetAvailableSlots", vetID).Return(mockSlots, nil)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	c.Params = gin.Params{{Key: "id", Value: vetID}}

	controller.GetAvailableSlots(c)

	assert.Equal(t, http.StatusOK, w.Code)
	mockService.AssertExpectations(t)
}

// 404 Not Found
func TestGetAvailableSlots_NotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockService := new(MockSlotService)
	controller := NewSlotController(mockService)

	vetID := "V-001"
	mockService.On("GetAvailableSlots", vetID).Return(nil, repositories.ErrNotFound)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = gin.Params{{Key: "id", Value: vetID}}

	controller.GetAvailableSlots(c)

	assert.Equal(t, http.StatusNotFound, w.Code)
	mockService.AssertExpectations(t)
}

// 500 Internal Server Error
func TestGetAvailableSlots_InternalError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockService := new(MockSlotService)
	controller := NewSlotController(mockService)

	vetID := "V-001"
	mockService.On("GetAvailableSlots", vetID).Return(nil, errors.New("unexpected database error"))

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = gin.Params{{Key: "id", Value: vetID}}

	controller.GetAvailableSlots(c)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	mockService.AssertExpectations(t)
}

// =====================================================================
//  Test: AddSlot Controller
// =====================================================================

// 201 Created (กรณีเพิ่มสำเร็จ)
func TestAddSlot_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockService := new(MockSlotService)
	controller := NewSlotController(mockService)

	// 1. จำลองข้อมูลที่จะส่งมา
	mockSlot := models.Slot{Date: "2026-06-01", TimePeriod: "10:00-11:00", SlotLimit: 5}

	// 2. แปลง mockSlot เป็น JSON
	jsonValue, _ := json.Marshal(mockSlot)

	// 3. บอก Mock ว่าถ้าเรียก AddSlot ให้รันผ่าน
	mockService.On("AddSlot", mock.Anything).Return(nil)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	// 4. จำลอง Role ผู้ช่วย
	c.Set("role", "assistant")

	// 5. จำลอง URL Params
	c.Params = gin.Params{{Key: "id", Value: "V-002"}}

	// 6. จำลองการส่ง Request พร้อม Body (JSON)
	req, _ := http.NewRequest("POST", "/api/vets/V-002/slots", bytes.NewBuffer(jsonValue))
	req.Header.Set("Content-Type", "application/json") // สำคัญ: บอกว่าส่ง JSON มานะ
	c.Request = req

	controller.AddSlot(c)

	assert.Equal(t, http.StatusCreated, w.Code)
	mockService.AssertExpectations(t)
}

// 403 Forbidden (กรณีไม่ใช่ผู้ช่วย)
func TestAddSlot_Forbidden(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockService := new(MockSlotService)
	controller := NewSlotController(mockService)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	// 📌 จำลองว่าคนนี้ Login มาเป็นหมอ (vet) ไม่ใช่ผู้ช่วย
	c.Set("role", "vet")

	controller.AddSlot(c)

	// ต้องโดนเด้งออก 403 Forbidden
	assert.Equal(t, http.StatusForbidden, w.Code)
}

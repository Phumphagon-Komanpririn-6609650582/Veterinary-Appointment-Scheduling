package controllers

import (
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

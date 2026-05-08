package controllers

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"veterinary-api/models"
	"veterinary-api/repositories"
	"veterinary-api/services"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockVetService struct {
	mock.Mock
	services.VetService // ฝัง Interface เพื่อกัน Go โวยวายว่าฟังก์ชันไม่ครบ
}

func (m *MockVetService) GetAllVets() ([]models.User, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.User), args.Error(1)
}

// =====================================================================
// Test: GetAllVets
// =====================================================================

// 200 OK
func TestGetAllVets_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockService := new(MockVetService)

	controller := NewVetController(mockService)

	mockVets := []models.User{
		{ID: "V-001", Username: "dr_somchai", Name: "นสพ.สมชาย", Role: "vet"},
	}

	mockService.On("GetAllVets").Return(mockVets, nil)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	controller.GetAllVets(c)

	assert.Equal(t, http.StatusOK, w.Code)
	mockService.AssertExpectations(t)
}

// 500 Internal Server Error
func TestGetAllVets_Error(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockService := new(MockVetService)
	controller := NewVetController(mockService)

	mockService.On("GetAllVets").Return(nil, errors.New("service error"))

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	controller.GetAllVets(c)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	mockService.AssertExpectations(t)
}

// 404 NotFound
func TestGetAllVets_NotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockService := new(MockVetService)
	controller := NewVetController(mockService)

	mockService.On("GetAllVets").Return(nil, repositories.ErrNotFound)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	controller.GetAllVets(c)

	assert.Equal(t, http.StatusNotFound, w.Code)
	mockService.AssertExpectations(t)
}

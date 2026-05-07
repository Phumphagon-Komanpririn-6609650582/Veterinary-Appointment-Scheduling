package services

import (
	"errors"
	"testing"
	"veterinary-api/repositories"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestGetAllVets_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	rows := sqlmock.NewRows([]string{"id", "username", "name", "role"}).
		AddRow("U-001", "vet01", "นสพ.สมชาย", "vet").
		AddRow("U-002", "vet02", "สพ.ญ.สมศรี", "vet")

	mock.ExpectQuery("SELECT id, username, name, role FROM users WHERE role = 'vet'").
		WillReturnRows(rows)

	repo := repositories.NewVetRepository(db)
	service := NewVetService(repo)

	vets, err := service.GetAllVets()

	assert.NoError(t, err)
	assert.Len(t, vets, 2)
	assert.Equal(t, "นสพ.สมชาย", vets[0].Name)
}

func TestGetAllVets_Error(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	mock.ExpectQuery("SELECT id, username, name, role FROM users WHERE role = 'vet'").
		WillReturnError(errors.New("database connection failed"))

	repo := repositories.NewVetRepository(db)
	service := NewVetService(repo)

	vets, err := service.GetAllVets()

	assert.Error(t, err)
	assert.Nil(t, vets)
	assert.Equal(t, "database connection failed", err.Error())
}

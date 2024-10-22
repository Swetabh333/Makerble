package test_models

import (
	"testing"

	"github.com/Swetabh333/Makerble/app/models"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestDoctorCreation(t *testing.T) {
	user := models.User{
		ID:   uuid.New(),
		Name: "TestUser",
	}

	doctor := models.Doctor{
		ID:     uuid.New(),
		UserID: user.ID,
		Name:   "TestDoctor",
	}

	assert.NotNil(t, doctor)
	assert.Equal(t, "TestDoctor", doctor.Name)
	assert.Equal(t, user.ID, doctor.UserID)
}

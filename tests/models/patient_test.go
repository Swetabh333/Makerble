package test_models

import (
	"testing"

	"github.com/Swetabh333/Makerble/app/models"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestPatientCreation(t *testing.T) {
	doctor := models.Doctor{
		ID:   uuid.New(),
		Name: "TestDoctor",
	}

	patient := models.Patient{
		ID:       uuid.New(),
		Name:     "TestPatient",
		Age:      30,
		Gender:   "male",
		DoctorID: doctor.ID,
	}

	assert.NotNil(t, patient)
	assert.Equal(t, "TestPatient", patient.Name)
	assert.Equal(t, 30, patient.Age)
	assert.Equal(t, "male", patient.Gender)
	assert.Equal(t, doctor.ID, patient.DoctorID)
}

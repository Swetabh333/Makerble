package test_models

import (
	"testing"

	"github.com/Swetabh333/Makerble/app/models"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestUserCreation(t *testing.T) {
	user := models.User{
		ID:       uuid.New(),
		Name:     "TestUser",
		Password: "TestPassword123",
		Role:     "doctor",
	}

	assert.NotNil(t, user)
	assert.Equal(t, "TestUser", user.Name)
	assert.Equal(t, "doctor", user.Role)
}

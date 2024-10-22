package test_routes

import (
	"testing"

	"github.com/Swetabh333/Makerble/app/routes"
	"github.com/stretchr/testify/assert"
)

func TestHashPassword(t *testing.T) {
	password := "TestPassword123"
	hashedPassword, err := routes.HashPassword(password)
	assert.NoError(t, err)
	assert.NotEqual(t, password, hashedPassword)
	assert.True(t, routes.CheckPasswordHash(password, hashedPassword))
}

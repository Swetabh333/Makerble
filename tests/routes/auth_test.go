package test_routes

import (
	"testing"

	"github.com/Swetabh333/Makerble/app/routes"
	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
)

func TestPasswordValidation(t *testing.T) {
	validate := validator.New()
	validate.RegisterValidation("password", routes.PasswordValidation)

	tests := []struct {
		password string
		valid    bool
	}{
		{"Valid1Password", true},
		{"invalidpassword", false},
		{"12345678", false},
		{"short", false},
		{"NoNumbersHere", false},
	}

	for _, tt := range tests {
		t.Run(tt.password, func(t *testing.T) {
			err := validate.Var(tt.password, "password")
			assert.Equal(t, tt.valid, err == nil)
		})
	}
}

package test_routes

import (
	"testing"

	"github.com/Swetabh333/Makerble/app/routes"
	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
)

func TestRoleValidation(t *testing.T) {
	validate := validator.New()
	validate.RegisterValidation("role", routes.RoleValidation)

	tests := []struct {
		role  string
		valid bool
	}{
		{"doctor", true},
		{"receptionist", true},
		{"admin", false},
		{"", false},
	}

	for _, tt := range tests {
		t.Run(tt.role, func(t *testing.T) {
			err := validate.Var(tt.role, "role")
			assert.Equal(t, tt.valid, err == nil)
		})
	}
}

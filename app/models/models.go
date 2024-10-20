package models

import "github.com/google/uuid"

type User struct {
	ID       uuid.UUID `gorm:"primary key"`
	Name     string    `gorm:"size:255;not null;unique"`
	Password string    `gorm:"size:255;not null"`
	Role     string    `gorm:"size:255;not null"`
}

type Patient struct {
	ID       uuid.UUID `gorm:"primary key"`
	Name     string    `gorm:"primary key;size:255;not null;unique"`
	Age      int       `gorm:"not null"`
	Gender   string
	DoctorID uuid.UUID `gorm:"not null"`
	User     User      `gorm:"foreignKey:DoctoID;constraint:OnUpdate:CASCADE;OnDelete:CASCADE"`
}

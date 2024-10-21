package models

import "github.com/google/uuid"

// Structure for  registering and storing users
type User struct {
	ID       uuid.UUID `gorm:"primary key"`
	Name     string    `gorm:"size:255;not null;unique"`
	Password string    `gorm:"size:255;not null"`
	Role     string    `gorm:"size:255;not null"`
}

// Structure for doctor information connected to the User table via UserID field
type Doctor struct {
	ID       uuid.UUID `gorm:"primary key"`
	UserID   uuid.UUID `gorm:"not null"`
	Name     string    `gorm:"not null"`
	Patients []Patient
	User     User `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE;OnDelete:CASCADE"`
}

// Structure for patient information connected to the Doctor table via DoctorID field.
type Patient struct {
	ID       uuid.UUID `gorm:"primary key"`
	Name     string    `gorm:"primary key;size:255;not null;unique"`
	Age      int       `gorm:"not null"`
	Gender   string
	DoctorID uuid.UUID `gorm:"not null"`
	Doctor   Doctor    `gorm:"foreignKey:DoctorID;constraint:OnUpdate:CASCADE;OnDelete:CASCADE"`
}

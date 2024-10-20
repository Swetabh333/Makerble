package models

import "github.com/google/uuid"

type User struct {
	ID       uuid.UUID `gorm:"primary key"`
	Name     string    `gorm:"size:255;not null;unique"`
	Password string    `gorm:"size:255;not null"`
	Role     string    `gorm:"size:255;not null"`
}

type Doctor struct {
	ID         uuid.UUID `gorm:"primary key"`
	UserID     uuid.UUID `gorm:"not null"`
	Name       string    `gorm:"not null"`
	Speciality string
	Patients   []Patient
	User       User `gorm:"foreignKey:UserID"`
}

type Patient struct {
	ID       uuid.UUID `gorm:"primary key"`
	Name     string    `gorm:"primary key;size:255;not null;unique"`
	Age      int       `gorm:"not null"`
	Gender   string
	DoctorID uuid.UUID `gorm:"not null"`
	Doctor   Doctor    `gorm:"foreignKey:DoctorID;constraint:OnUpdate:CASCADE;OnDelete:CASCADE"`
}

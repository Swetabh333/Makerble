package models

import "github.com/google/uuid"

type User struct {
	ID       uuid.UUID `gorm:"primary key"`
	Name     string    `gorm:"size:255;not null;unique"`
	Password string    `gorm:"size:255;not null"`
	Role     string    `gorm:"size:255;not null"`

}



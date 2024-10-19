package databases

import (
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

//This function returns a database instance after connecting to our postgres database , this instance can later be used to perform database operations.

func ConnectToDatabase() (*gorm.DB, error) {
	dsn := os.Getenv("DSN_STRING")
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	return db, nil
}

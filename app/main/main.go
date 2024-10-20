package main

import (
	"log"

	"github.com/Swetabh333/Makerble/app/databases"
	"github.com/Swetabh333/Makerble/app/models"
	"github.com/Swetabh333/Makerble/app/routes"
	//"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	//Load the env file
	godotenv.Load(".env")

	//Connect to the database and return db instance
	db, err := databases.ConnectToDatabase()
	if err != nil {
		log.Fatal("Could not connect to the database")
	}
	log.Print("connected to the database")

	//Migrating data to postgres
	err = db.AutoMigrate(&models.User{}, &models.Doctor{}, &models.Patient{})
	if err != nil {
		log.Fatal("Error migrating data to database")
	}

	//gin.SetMode(gin.ReleaseMode)

	//Create a router for routing requests
	router := routes.NewRouter()

	//routes
	router.GET("/ping", routes.HandlePong)

	//Start out http server at port 8080
	router.Run("0.0.0.0:8080")
}

package main

import (
	"log"

	"github.com/Swetabh333/Makerble/app/databases"
	"github.com/Swetabh333/Makerble/app/routes"
	"github.com/joho/godotenv"
)

func main() {
	//Load the env file
	godotenv.Load(".env")

	//Connect to the database and return db instance
	_, err := databases.ConnectToDatabase()
	if err != nil {
		log.Fatal("Could not connect to the database")
	}
	log.Print("connected to the database")

	//Create a router for routing requests
	router := routes.NewRouter()

	//Start out http server at port 8080
	router.Run("0.0.0.0:8080")
}

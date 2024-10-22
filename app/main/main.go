package main

import (
	"log"
	"time"

	"github.com/Swetabh333/Makerble/app/databases"
	"github.com/Swetabh333/Makerble/app/middleware"
	"github.com/Swetabh333/Makerble/app/models"
	"github.com/Swetabh333/Makerble/app/routes"
	"github.com/gin-contrib/cors"

	"github.com/joho/godotenv"
)

func main() {
	//Load the env file
	if err := godotenv.Load(".env"); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}
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

	//setting up CORS for added security

	config := cors.Config{
		AllowAllOrigins:  true,
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}

	router.Use(cors.New(config))
	//routes

	//ping to check if server is alive
	router.GET("/ping", routes.HandlePing)

	//Auth routes

	//for registering a new user
	router.POST("/auth/register", routes.RegisterHandler(db))
	//for logging in
	router.POST("/auth/login", routes.LoginHandler(db))

	// Patient routes

	router.POST("/patients/add", middleware.VerifyAuthentication, routes.AddPatient(db))

	//Start out http server at port 8080
	if err := router.Run("0.0.0.0:8080"); err != nil {
		log.Fatal("Error starting server: ", err)
	}
}

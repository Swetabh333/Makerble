package main

import (
	"log"
	"time"

	"github.com/Swetabh333/Makerble/app/databases"
	"github.com/Swetabh333/Makerble/app/middleware"
	"github.com/Swetabh333/Makerble/app/models"
	"github.com/Swetabh333/Makerble/app/routes"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

var db *gorm.DB
var redisClient *redis.Client

//init function runs before man to set everything up for us

func init() {

	//Connect to the database and return db instance
	db, err := databases.ConnectToDatabase()
	if err != nil {
		log.Fatal("Could not connect to the database")
	}
	log.Print("connected to the database")

	//Connect to redis
	redisClient, err = databases.ConnectToRedis()
	if err != nil {
		log.Fatal("Could not connect to redis ", err)
	}

	//Migrating data to postgres
	err = db.AutoMigrate(&models.User{}, &models.Doctor{}, &models.Patient{})
	if err != nil {
		log.Fatal("Error migrating data to database")
	}

}

func main() {

	gin.SetMode(gin.ReleaseMode)

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

	// route for creating a new patient - can only be done by a recptionist
	router.POST("/patient", middleware.VerifyAuthentication, routes.AddPatient(db))
	//route for getting a patients-data - can be accessed by both doctors and receptionists
	router.GET("/patient/:id", middleware.VerifyAuthentication, routes.GetPatient(db, redisClient))
	//route to fetch all patients - accessible by both doctor and receptionist
	router.GET("/patients", middleware.VerifyAuthentication, routes.GetAllPatients(db, redisClient))
	//route for deleting a patient using his id - accessibly by receptionist
	router.DELETE("/patient/:id", middleware.VerifyAuthentication, routes.DeletePatient(db))
	//route for updating a patients records - accessible by both doctors and receptionists
	router.PUT("/patient/:id", middleware.VerifyAuthentication, routes.UpdatePatient(db))

	//Start out http server at port 8080
	if err := router.Run("0.0.0.0:8080"); err != nil {
		log.Fatal("Error starting server: ", err)
	}
}

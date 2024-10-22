package routes

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/Swetabh333/Makerble/app/models"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type addRequest struct {
	Name   string `json:"name" validate:"required"`
	Age    int    `json:"age" validate:"required,gte=0"`
	Gender string `json:"gender" validate:"required,gender"`
	Doctor string `json:"doctorName" validate:"required"`
}

type updateRequest struct {
	Name      string `json:"name" validate:"required"`
	Age       int    `json:"age" validate:"required,gte=0"`
	Gender    string `json:"gender" validate:"required,gender"`
	Diagnosis string `json:"diagnosis"`

	Doctor string `json:"doctorName" validate:"required"`
}

var validGenders = []string{"male", "female", "other"}

func genderValidation(fl validator.FieldLevel) bool {
	gender := strings.ToLower(fl.Field().String())
	for _, validGender := range validGenders {
		if gender == validGender {
			return true
		}
	}
	return false
}

func AddPatient(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		role, ok := c.Get("Role")
		if !ok || role != "receptionist" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized to access this route"})
		}
		request := addRequest{}
		err := c.BindJSON(&request)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "Some internal error occured",
			})
			return

		}
		//To make sure the requuest is valid
		validate := validator.New()
		validate.RegisterValidation("gender", genderValidation)
		// Validate request data
		err = validate.Struct(&request)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}
		doctor := models.Doctor{}

		err = db.Where("name = ?", request.Doctor).Find(&doctor).Error
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "The doctor requested does not exist",
			})
		}

		existingPatient := models.Patient{}

		err = db.Where("name = ? AND age = ? AND gender = ?", request.Name, request.Age, request.Gender).Find(&existingPatient).Error
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "The patient already exists",
			})
		}

		patient := models.Patient{
			ID:       uuid.New(),
			Name:     request.Name,
			Age:      request.Age,
			Gender:   request.Gender,
			DoctorID: doctor.ID,
		}

		err = db.Create(&patient).Error
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to add patient",
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"message": "Patient added successfully", "patient": patient,
		})

	}
}

// Route to get a particular patients data using their id - using caching
func GetPatient(db *gorm.DB, rdb *redis.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		role, ok := c.Get("Role")
		if !ok || (role != "doctor" && role != "receptionist") {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized to access this route"})
			return
		}

		patientID := c.Param("id")
		cacheKey := fmt.Sprintf("patient:%s", patientID)

		// Try to get from redis first
		cachedData, err := rdb.Get(c.Request.Context(), cacheKey).Result()
		if err == nil {
			// Cache hit - found in redis
			patient := models.Patient{}
			err := json.Unmarshal([]byte(cachedData), &patient)
			if err == nil {
				c.JSON(http.StatusOK, gin.H{"patient": patient})
				return
			}
		}

		// Cache miss, get from postgres database
		patient := models.Patient{}
		err = db.Where("id = ?", patientID).First(&patient).Error
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Patient not found"})
			return
		}

		// Store in redis for future requests
		patientJson, err := json.Marshal(patient)
		if err == nil {
			rdb.Set(c.Request.Context(), cacheKey, patientJson, 30*time.Minute) // Cache for 30 minutes
		}

		c.JSON(http.StatusOK, gin.H{"patient": patient})
	}
}

// Route for getting all patients
func GetAllPatients(db *gorm.DB, rdb *redis.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		role, ok := c.Get("Role")
		if !ok || (role != "doctor" && role != "receptionist") {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized to access this route"})
			return
		}

		cacheKey := "all_patients"
		cachedData, err := rdb.Get(c.Request.Context(), cacheKey).Result()
		if err == nil {
			var patients []models.Patient
			if err := json.Unmarshal([]byte(cachedData), &patients); err == nil {
				if len(patients) == 0 {
					c.JSON(http.StatusOK, gin.H{"message": "No patients found"})
					return
				}
				c.JSON(http.StatusOK, gin.H{"patients": patients})
				return
			}
		}

		patients := []models.Patient{}
		err = db.Find(&patients).Error
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch patients"})
			return
		}

		if len(patients) == 0 {
			c.JSON(http.StatusOK, gin.H{"message": "No patients found"})
			return
		}

		if patientsJson, err := json.Marshal(patients); err == nil {
			rdb.Set(c.Request.Context(), cacheKey, patientsJson, 15*time.Minute) // Cache for 15 minutes
		}

		c.JSON(http.StatusOK, gin.H{"patients": patients})
	}
}

// Route to delete a patient's data
func DeletePatient(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		role, ok := c.Get("Role")
		if !ok || role != "receptionist" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized to delete this patient"})
			return
		}

		patientID := c.Param("id")
		patient := models.Patient{}

		err := db.Where("id = ?", patientID).First(&patient).Error
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Patient not found"})
			return
		}

		err = db.Delete(&patient).Error

		if err != nil {

			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete patient"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Patient deleted successfully"})
	}
}

// Function that accepts all the parameters and updates the ones that have been changed
func UpdatePatient(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		role, ok := c.Get("Role")
		if !ok || (role != "doctor" && role != "receptionist") {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized to access this route"})
			return
		}

		patientID := c.Param("id")
		patient := models.Patient{}

		err := db.Where("id = ?", patientID).First(&patient).Error
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Patient not found"})
			return
		}

		updateRequest := updateRequest{}
		if err := c.BindJSON(&updateRequest); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
			return
		}

		validate := validator.New()
		validate.RegisterValidation("gender", genderValidation)
		err = validate.Struct(&updateRequest)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		doctor := models.Doctor{}
		err = db.Where("name = ?", updateRequest.Doctor).First(&doctor).Error
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Doctor not found"})
			return
		}

		patient.Name = updateRequest.Name
		patient.Age = updateRequest.Age
		patient.Gender = updateRequest.Gender
		patient.DoctorID = doctor.ID
		patient.Diagnosis = updateRequest.Diagnosis
		err = db.Save(&patient).Error
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update patient"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Patient updated successfully", "patient": patient})
	}
}

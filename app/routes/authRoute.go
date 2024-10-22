package routes

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/Swetabh333/Makerble/app/helper"
	"github.com/Swetabh333/Makerble/app/models"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// defing register json payload
type register struct {
	Name     string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required,min=8,max=64,password"`
	Role     string `json:"role" validate:"required"`
}

type login struct {
	Name     string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

// function to validate if password is valid or not - password should have an uppercase letter ,  a lowercase letter , a number and be at least 8 characters
func passwordValidation(f1 validator.FieldLevel) bool {
	password := f1.Field().String()
	hasMinLen := len(password) >= 8
	hasUpper := false
	hasLower := false
	hasNumber := false
	for _, char := range password {
		switch {
		case 'A' <= char && char <= 'Z':
			hasUpper = true
		case 'a' <= char && char <= 'z':
			hasLower = true
		case '0' <= char && char <= '9':
			hasNumber = true
		}
	}
	return hasMinLen && hasUpper && hasLower && hasNumber
}

// function to encrypt the password before storing in database
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

// function to decrypt and compare the password while loggin in
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// For validating roles
var validRoles = []string{"doctor", "receptionist"}

func roleValidation(fl validator.FieldLevel) bool {
	role := strings.ToLower(fl.Field().String())
	for _, validRole := range validRoles {
		if role == validRole {
			return true
		}
	}
	return false
}

func RegisterHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {

		// extracting the body of request
		regBody := register{}
		validate := validator.New()
		validate.RegisterValidation("password", passwordValidation)
		validate.RegisterValidation("role", roleValidation)

		err := c.BindJSON(&regBody)
		if err != nil {
			fmt.Println("Error binding request body")
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Some internal error occured",
			})
			return
		}
		//password length validation
		if err = validate.Struct(&regBody); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "password should be at least 8 characters and contain at least 1 number, 1 uppercase character and 1 lower case character",
			})

			return
		}
		//encrypting password before storing
		hashedPassword, err := HashPassword(regBody.Password)
		if err != nil {
			fmt.Println("Error hashing password")
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Some internal error occured",
			})
			return

		}
		user := models.User{
			ID:       uuid.New(),
			Name:     regBody.Name,
			Role:     regBody.Role,
			Password: hashedPassword,
		}

		err = db.Create(&user).Error
		if err != nil {
			if err.Error() == `ERROR: duplicate key value violates unique constraint "uni_users_name" (SQLSTATE 23505)` {
				c.JSON(http.StatusBadRequest, gin.H{
					"error": "User already exists",
				})

				return
			}
			fmt.Printf("Error storing user in database: %s\n", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Some internal error occured",
			})
			return
		}
		if user.Role == "doctor" {
			doctor := models.Doctor{
				Name:   user.Name,
				ID:     uuid.New(),
				UserID: user.ID,
			}
			err = db.Create(&doctor).Error
			if err != nil {

				fmt.Printf("Error storing user in database: %s\n", err.Error())
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": "Some internal error occured",
				})
				return
			}
		}
		fmt.Println("User successfully registered")
		c.JSON(http.StatusOK, gin.H{
			"message": "User successfully created",
		})

	}
}

func LoginHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		login := login{}
		user := models.User{}
		err := c.BindJSON(&login)
		if err != nil {
			fmt.Println("Error binding request body")
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "Some internal error occured",
			})
			return
		}
		validate := validator.New()
		if err = validate.Struct(&login); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})

		err = db.Where("name = ?", login.Name).Find(&user).Error
		if err != nil {
			fmt.Println("User not found")
			c.JSON(http.StatusUnauthorized, gin.H{
				"message": "Username does not exist",
			})
			return
		}
		check := CheckPasswordHash(login.Password, user.Password)
		if !check {
			fmt.Println("Incorrect password")
			c.JSON(http.StatusUnauthorized, gin.H{
				"message": "username and password do not match",
			})
			return

		}
		token, err := helper.GenerateToken(user.Name, user.ID.String(), user.Role)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "Error generating cookies",
			})
		}
		refreshToken, err := helper.GenerateRefreshToken(user.Name, user.ID.String(), user.Role)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "Error generating cookies",
			})
		}
		c.SetCookie("token", token, 24*3600, "/", "", false, true)
		c.SetCookie("refreshToken", refreshToken, 30*24*3600, "/", "", false, true)
		fmt.Println("Logged in")

		c.JSON(http.StatusOK, gin.H{
			"message": "Logged in successfully",
		})
		return

	}
}

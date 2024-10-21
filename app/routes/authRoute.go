package routes

import (
	"fmt"
	"net/http"

	"github.com/Swetabh333/Makerble/app/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// defing register json payload
type register struct {
	Name     string `json:"username"`
	Password string `json:"password"`
	Role     string `json:"role"`
}

type login struct {
	Name     string `json:"username"`
	Password string `json:"password"`
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

func RegisterHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {

		// extracting the body of request
		regBody := register{}
		err := c.BindJSON(&regBody)
		if err != nil {
			fmt.Println("Error binding request body")
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Some internal error occured",
			})
			return
		}
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
		//login := login{}
		//user := models.User{}
	}
}

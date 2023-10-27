// main.go

package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type User struct {
	ID       uint   `json:"id" gorm:"primaryKey"`
	Username string `json:"username"`
	Password string `json:"-"`
}

var db *gorm.DB

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		os.Getenv("DB_USERNAME"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"),
	)

	var errDb error
	db, errDb = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if errDb != nil {
		log.Fatal(errDb)
	}

	db.AutoMigrate(&User{})

	router := gin.Default()

	api := router.Group("/api")
	{
		api.POST("/register", register)
		api.POST("/login", login)
		api.GET("/health", HealthCheckHandler)
	}

	port := "4000"
	router.Run(":" + port)
}

type HealthCheckResponse struct {
	Status string `json:"status"`
}

func HealthCheckHandler(c *gin.Context) {
	response := HealthCheckResponse{Status: "OK"}
	c.JSON(http.StatusOK, response)
}

func register(c *gin.Context) {
	var user User
	c.BindJSON(&user)

	var existingUser User
	result := db.Where("username = ?", user.Username).First(&existingUser)
	if result.RowsAffected != 0 {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Username already exists"})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Internal server error"})
		return
	}

	user.Password = string(hashedPassword)

	db.Create(&user)
	c.JSON(http.StatusOK, gin.H{"message": "Registration successful"})
}

func login(c *gin.Context) {
	var user User
	c.BindJSON(&user)

	var users []User
	result := db.Where("username = ?", user.Username).Find(&users)
	if result.RowsAffected == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid credentials"})
		return
	}

	for _, u := range users {
		if u.Password == user.Password {
			c.JSON(http.StatusOK, gin.H{"message": "Login successful"})
			return
		}
	}

	c.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid credentials"})
}

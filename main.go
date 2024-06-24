package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

func main() {
	var err error

	err = godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	// Load database credentials from environment variables
	dsn := os.Getenv("MYSQL_DSN")
	if dsn == "" {
		log.Fatal("MYSQL_DSN environment variable is required")
	}

	db, err = sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		log.Fatalf("Failed to ping the database: %v", err)
	}

	router := gin.Default()
	router.POST("/submit-form", submitFormHandler)

	// Use environment variable for port or default to 8080
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	if err := router.Run(":" + port); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}

type FormInput struct {
	Name     string `json:"name" binding:"required"`
	Whatsapp string `json:"whatsapp" binding:"required"`
	Message  string `json:"message" binding:"required"`
}

func submitFormHandler(c *gin.Context) {
	var input FormInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	stmt, err := db.Prepare("INSERT INTO submissions (name, whatsapp, message) VALUES (?, ?, ?)")
	if err != nil {
		log.Printf("Failed to prepare statement: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to submit data"})
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec(input.Name, input.Whatsapp, input.Message)
	if err != nil {
		log.Printf("Failed to execute statement: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to submit data"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Data submitted successfully!"})
}

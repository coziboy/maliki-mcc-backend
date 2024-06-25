package main

import (
	"database/sql"
	"log"
	"net/http"
	"time"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

func main() {
	var err error
	db, err = sql.Open("mysql", "root:MuLFjpjCHAKGLKBkOtXvIhWPbBIrdbAD@tcp(viaduct.proxy.rlwy.net:41263)/railway")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	router := gin.Default()

	// CORS middleware configuration
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"https://mcc-coziboy.up.railway.app"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	router.POST("/submit-form", submitFormHandler)
	router.Run(envPortOr("8080"))
}

// Returns PORT from environment if found, defaults to
// value in `port` parameter otherwise. The returned port
// is prefixed with a `:`, e.g. `":3000"`.
func envPortOr(port string) string {
	// If `PORT` variable in environment exists, return it
	if envPort := os.Getenv("PORT"); envPort != "" {
	  return ":" + envPort
	}
	// Otherwise, return the value of `port` variable from function argument
	return ":" + port
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
		log.Fatal(err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(input.Name, input.Whatsapp, input.Message)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to submit data"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Data submitted successfully!"})
}

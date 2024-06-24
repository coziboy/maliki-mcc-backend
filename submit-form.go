package main

import (
    "database/sql"
    "fmt"
    "log"
    "net/http"

    "github.com/gin-gonic/gin"
    _ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

func main() {
    var err error
    db, err = sql.Open("mysql", "username:password@tcp(your_server:port)/your_database")
    if err != nil {
        log.Fatal(err)
    }
    defer db.Close()

    err = db.Ping()
    if err != nil {
        log.Fatal(err)
    }

    router := gin.Default()
    router.POST("/submit-form", submitFormHandler)
    router.Run(":8080")
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

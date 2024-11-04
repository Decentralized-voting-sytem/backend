package database

import (
    "fmt"
    "log"
    "yourproject/models" // Import your models package

    "gorm.io/driver/postgres"
    "gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDatabase() {
    var err error
    dsn := "host=localhost user=postgres password=samyak dbname=mydbname port=5432 sslmode=disable"
    DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
    if err != nil {
        log.Fatalf("Failed to connect to database: %v", err)
    }

    err = DB.AutoMigrate(&models.Voter{}, &models.Candidate{}, &models.Vote{}, &models.Admin{})
    if err != nil {
        log.Fatalf("Failed to migrate database: %v", err)
    }

    fmt.Println("Database connection established and models migrated")
}

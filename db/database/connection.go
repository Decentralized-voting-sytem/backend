package database

import (
    "log"
    "github.com/Decentralized-voting-sytem/backend/db/models"
    "gorm.io/driver/postgres"
    "gorm.io/gorm"
)

var DB *gorm.DB

func Init() {
    var err error
    dsn := "host=localhost user=postgres password=samyak dbname=voting port=5432 sslmode=disable"
    DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
    if err != nil {
        log.Fatalf("Failed to connect to database: %v", err)
    }

    err = DB.AutoMigrate(&models.Vote{}, &models.Voter{}, &models.Candidate{}, &models.Admin{})
    if err != nil {
        log.Fatalf("Failed to migrate database: %v", err)
    }
}

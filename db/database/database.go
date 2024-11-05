package database

import (
    "sync"
    "log"
    "github.com/Decentralized-voting-sytem/backend/db/models"
    "gorm.io/driver/postgres"
    "gorm.io/gorm"
)

var DBLock sync.Mutex

func Init() *gorm.DB{
    DBLock.Lock()
	defer DBLock.Unlock()

	dsn := "host=localhost user=postgres password=samyak dbname=voting port=5432 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}

	db.AutoMigrate(&models.Vote{}, &models.Voter{}, &models.Candidate{}, &models.Admin{})

	return db
}


package db

import (
	 "gorm.io/gorm"
	 "github.com/Decentralized-voting-sytem/backend/db/database"
)

var DB *gorm.DB

func InitDB() {
	DB= database.Init()
}
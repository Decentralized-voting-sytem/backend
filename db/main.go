package main

import (
	"gorm.io/gorm"
"github.com/Decentralized-voting-sytem/backend/db/database"
)

var DB *gorm.DB

func main() {
	DB = database.Init()
}
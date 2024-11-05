package main

import (
    "github.com/Decentralized-voting-sytem/backend/db/database"
    "github.com/Decentralized-voting-sytem/backend/services/auth/controllers"
    "github.com/gin-contrib/cors"
    "github.com/gin-gonic/gin"
    "gorm.io/gorm"
)

var DB *gorm.DB

func main() {
    DB = database.Init()
    r := gin.Default()

    config := cors.DefaultConfig()
    config.AllowOrigins = []string{"http://localhost:3000"}
    config.AllowHeaders = []string{"Content-Type"}
    config.AllowCredentials = true
    r.Use(cors.New(config))

    r.POST("/login", controllers.Login)

    r.Run()
}

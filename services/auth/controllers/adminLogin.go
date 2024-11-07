package controllers

import (
		"fmt"
		"net/http"
		"github.com/Decentralized-voting-sytem/backend/db/database"
		"github.com/Decentralized-voting-sytem/backend/db/models"
		"github.com/gin-gonic/gin"
)
	

func AdminLogin(c *gin.Context) {
	var body struct {
		AdminID  string `json:"admin_id"`
		Name     string `json:"name"`
		Password string `json:"password"`
	}

	if c.Bind(&body) != nil {
		c.JSON(402, gin.H{
			"error": "Failed to read the body",
		})
		return
	}

	var admin models.Admin
	query := `SELECT * FROM admins WHERE admin_id = ? AND name = ? AND password = ?`
	res := database.DB.Raw(query, body.AdminID, body.Name, body.Password).First(&admin)
	if res.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": res.Error.Error()})
		return
	}

	if res.RowsAffected == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		fmt.Println(body)
		return
	}

	c.SetCookie(
		"admin_id",        
		body.AdminID,      
		3000,            
		"/",              
		"",                
		false,             
		true,              
	)

	c.JSON(200, gin.H{
		"message": "Admin Login successful",
	})
}
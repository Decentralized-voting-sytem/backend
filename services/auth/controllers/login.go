package controllers

import (
	"fmt"
	"net/http"
	"strings"
	"github.com/Decentralized-voting-sytem/backend/db/database"
	"github.com/Decentralized-voting-sytem/backend/db/models"
	"github.com/gin-gonic/gin"
)


func Login(c *gin.Context) {
	var body struct {
		VoterID  string `json:"voter_id"`
		Name     string `json:"name"`
		DOB      string `json:"dob"`
		Password string `json:"password"`
	}

	// Bind JSON body to struct
	if c.Bind(&body) != nil {
		c.JSON(402, gin.H{
			"error": "Failed to read the body",
		})
		return
	}

	var voter models.Voter

	// Query to check for matching voter credentials
	query1 := `SELECT * FROM voters WHERE voter_id = ? AND dob = ? AND password = ?`
	res := database.DB.Raw(query1, body.VoterID, body.DOB, body.Password).First(&voter)
	if res.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": res.Error.Error()})
		return
	}

	// Check for name match and handle unauthorized access
	if res.RowsAffected == 0 || strings.ToLower(voter.Name) != strings.ToLower(body.Name) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	// Set a cookie for the voter_id
	c.SetCookie(
		"voter_id",       
		body.VoterID,     
		300,              
		"/",              
		"",               
		false,            
		true,             
	)

	c.JSON(200, gin.H{
		"message": "Login successful",
	})
}

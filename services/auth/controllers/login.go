package controllers

import (
	"net/http"
	"strings"
	"github.com/Decentralized-voting-sytem/backend/db/database"
	"github.com/Decentralized-voting-sytem/backend/db/models"
	"github.com/Decentralized-voting-sytem/backend/services/auth/utils"
	"github.com/gin-gonic/gin"
)

func Login(c *gin.Context) {
	var body struct {
		VoterID  string `json:"voter_id"`
		Name     string `json:"name"`
		DOB      string `json:"dob"`
		Password string `json:"password"`
	}

	// Bind the incoming JSON request body to the body struct
	if c.Bind(&body) != nil {
		c.JSON(402, gin.H{
			"error": "Failed to read the body",
		})
		return
	}

	var voter models.Voter

	// Check credentials against the database
	query := `SELECT * FROM voters WHERE voter_id = ? AND dob = ? AND password = ?`
	res := database.DB.Raw(query, body.VoterID, body.DOB, body.Password).First(&voter)
	if res.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": res.Error.Error()})
		return
	}

	// Check if the name matches
	if res.RowsAffected == 0 || strings.ToLower(voter.Name) != strings.ToLower(body.Name) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	token, err := utils.GenerateVerificationToken(body.VoterID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate verification token"})
		return
	}

	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("Auth", token, 3600*24*30, "", "", false, false)

	// Return success response
	c.JSON(200, gin.H{
		"message": "Login successful",
	})
}

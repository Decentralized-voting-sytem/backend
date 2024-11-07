package controllers

import (
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

	if c.Bind(&body) != nil {
		c.JSON(402, gin.H{
			"error": "Failed to read the body",
		})
		return
	}

	var voter models.Voter

	query1 := `SELECT * FROM voters WHERE voter_id = ? AND dob = ? AND password = ?`
	res := database.DB.Raw(query1, body.VoterID, body.DOB, body.Password).First(&voter)
	if res.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": res.Error.Error()})
		return
	}

	if res.RowsAffected == 0 || strings.ToLower(voter.Name) != strings.ToLower(body.Name) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	query2 := `UPDATE voters SET logged = ? WHERE voter_id = ?`
	res2 := database.DB.Exec(query2, true, body.VoterID)
	if res2.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update login status"})
		return
	}

	c.JSON(200, gin.H{
		"message": "Login successful",
	})
}

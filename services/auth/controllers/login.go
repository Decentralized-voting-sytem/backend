package controllers

import (
	// "time"
	"github.com/gin-gonic/gin"
	"net/http"
	"github.com/Decentralized-voting-sytem/backend/db/models"
	"github.com/Decentralized-voting-sytem/backend/db/database"
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
	var vote models.Vote

	query1 := `SELECT * FROM voters WHERE voter_id = ? name = ? AND dob = ? password = ?`
	res := database.DB.Raw(query1, body.VoterID, body.Name, body.DOB, body.Password).Scan(&voter)
	if res.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": res.Error.Error()})
		return
	}

	if res.RowsAffected == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	query2 := `SELECT * FROM votes WHERE voter_id = ?`
	rep := database.DB.Raw(query2, body.VoterID).Scan(&vote)
	if rep.RowsAffected > 0 {
		c.JSON(http.StatusForbidden, gin.H{"error": "Voter has already voted"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Login successful",
		"voter": gin.H{
			"id":       voter.ID,
			"voter_id": voter.VoterID,
			"name":     voter.Name,
			"dob":      voter.DOB,
		},
	})
}

package controllers

import (
	"net/http"
	"time"
	"fmt"
	"github.com/Decentralized-voting-sytem/backend/db/database"
	"github.com/gin-gonic/gin"
)

func RegisterVote(c *gin.Context) {

	var body struct {
		CandidateID      string         `json:"candidate_id"`
	}

	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})
		return
	}

	query := `INSERT INTO votes (voter_id, candidate_id, created_at, updated_at) VALUES (?, ?, ?, ?)`
	result := database.DB.Exec(query, voterID, body.CandidateID, time.Now(), time.Now())

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to cast vote"})
		fmt.Println("Error:", result.Error)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Vote successfully cast"})
}
package controllers

import (
	"net/http"
	"github.com/Decentralized-voting-sytem/backend/db/database"
	// "github.com/Decentralized-voting-sytem/backend/db/models"
	// "github.com/Decentralized-voting-sytem/backend/services/auth/utils"
	"github.com/gin-gonic/gin"
)

// GetVoteInsights returns voting insights such as vote counts, most popular candidate, and vote percentages
func GetVoteInsights(ctx *gin.Context) {
	// Query 1: Get the total number of votes for each candidate (Aggregate Query with JOIN)
	var voteCounts []struct {
		CandidateName string `json:"candidate_name"`
		VoteCount     int    `json:"vote_count"`
	}
	query1 := `
		-- Aggregate Query with JOIN
		SELECT c.Name AS candidate_name, COUNT(v.ID) AS vote_count
		FROM votes v
		JOIN candidates c ON v.candidate_id = c.ID
		GROUP BY c.Name;
	`
	if err := database.DB.Raw(query1).Scan(&voteCounts).Error; err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Query 2: Get the most popular candidate (Aggregate Query with JOIN)
	var mostPopularCandidate struct {
		CandidateName string `json:"candidate_name"`
		VoteCount     int    `json:"vote_count"`
	}
	query2 := `
		-- Aggregate Query with JOIN
		SELECT c.Name AS candidate_name, COUNT(v.ID) AS vote_count
		FROM votes v
		JOIN candidates c ON v.candidate_id = c.ID
		GROUP BY c.Name
		ORDER BY vote_count DESC
		LIMIT 1;
	`
	if err := database.DB.Raw(query2).Scan(&mostPopularCandidate).Error; err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Query 3: Calculate the vote percentage for each candidate (Nested Query with Aggregate)
	var votePercentages []struct {
		CandidateName  string  `json:"candidate_name"`
		VotePercentage float64 `json:"vote_percentage"`
	}
	query3 := `
		-- Nested Query with Aggregate
		SELECT c.Name AS candidate_name, 
		       (COUNT(v.ID) * 100.0 / (SELECT COUNT(*) FROM votes)) AS vote_percentage
		FROM votes v
		JOIN candidates c ON v.candidate_id = c.ID
		GROUP BY c.Name;
	`
	if err := database.DB.Raw(query3).Scan(&votePercentages).Error; err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Prepare the insights to be returned
	insights := gin.H{
		"vote_counts":           voteCounts,
		"most_popular_candidate": mostPopularCandidate,
		"vote_percentages":      votePercentages,
	}

	// Return the insights as JSON response to the frontend
	ctx.JSON(http.StatusOK, insights)
}

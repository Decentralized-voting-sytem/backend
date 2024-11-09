package controllers

import (
	"github.com/golang-jwt/jwt/v5"
	"fmt"
	"net/http"
	"time"
	"github.com/Decentralized-voting-sytem/backend/db/database"
	"github.com/Decentralized-voting-sytem/backend/services/auth/utils"
	"github.com/gin-gonic/gin"
)

func RegisterVote(ctx *gin.Context) {

	var body struct {
		CandidateID  string `json:"candidate_id"`
	}

	token, err := ctx.Cookie("Auth")
	if err != nil {
		ctx.JSON(400, gin.H{"message": token})
		return
	}

	// Parse the JWT token
	parsedToken, err := utils.ParseToken(token)
	if err != nil {
		ctx.JSON(401, gin.H{"message": "Invalid token"})
		return
	}

	// Extract the voter ID from the token
	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	if !ok || !parsedToken.Valid {
		ctx.JSON(402, gin.H{"message": "Invalid token claims"})
		return
	}

	voterID, ok := claims["iss"].(string)
	if !ok {
		ctx.JSON(403, gin.H{"message": "Invalid voter ID in token"})
		return
	}

	// Bind the request body to the struct
	if err := ctx.Bind(&body); err != nil {
		ctx.JSON(401, gin.H{
			"error": fmt.Sprintf("Failed to read the body: %v", err),
		})
		fmt.Println("Bind error:", err)
		return
	}

	query := `INSERT INTO votes (voter_id, candidate_id, created_at, updated_at) VALUES (?, ?, ?, ?)`
	result := database.DB.Exec(query, voterID, body.CandidateID, time.Now(), time.Now())

	// Check for errors in the database query execution
	if result.Error != nil {
		ctx.JSON(500, gin.H{"error": "Failed to cast vote"})
		fmt.Println("Error:", result.Error)
		return
	}

	// Return success message
	ctx.JSON(http.StatusOK, gin.H{"message": "Vote successfully cast"})
}

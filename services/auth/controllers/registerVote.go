package controllers

import (
	"fmt"
	"net/http"
	"time"
	"github.com/Decentralized-voting-sytem/backend/db/database"
	"github.com/Decentralized-voting-sytem/backend/services/auth/utils"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func RegisterVote(ctx *gin.Context) {

	var body struct {
		CandidateID int `json:"candidate_id"`
	}

	// Retrieve the Auth token from the cookie
	token, err := ctx.Cookie("Auth")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Token not found"})
		return
	}

	// Parse the JWT token
	parsedToken, err := utils.ParseToken(token)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid token"})
		return
	}

	// Extract the voter ID from the token claims
	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	if !ok || !parsedToken.Valid {
		ctx.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid token claims"})
		return
	}

	voterID, ok := claims["iss"].(string)
	if !ok {
		ctx.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid voter ID in token"})
		return
	}

	// Bind the request body to the struct
	if err := ctx.Bind(&body); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Failed to read the body: %v", err)})
		return
	}

	// Insert the vote into the database
	query := `INSERT INTO votes (voter_id, candidate_id, created_at, updated_at) VALUES (?, ?, ?, ?)`
	result := database.DB.Exec(query, voterID, body.CandidateID, time.Now(), time.Now())

	// Check for errors in the database query execution
	if result.Error != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Voter has already voted"})
		return
	}

	// Delete the Auth cookie by setting it with an expired timestamp
	ctx.SetCookie("Auth", "", -1, "/", "", false, true)

	// Return success message
	ctx.JSON(http.StatusOK, gin.H{"message": "Vote successfully cast"})
}
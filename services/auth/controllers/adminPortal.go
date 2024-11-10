package controllers

import (
	"fmt"
	"log"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
	"github.com/Decentralized-voting-sytem/backend/db/database"
	"github.com/Decentralized-voting-sytem/backend/db/models"
	"github.com/Decentralized-voting-sytem/backend/services/auth/utils"
	"github.com/gin-gonic/gin"
)

func AdminPortal(ctx *gin.Context) {
	// Retrieve the Auth token from the cookie
	token, err := ctx.Cookie("AdminAuth")
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

	// Extract the admin ID from the token claims
	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	if !ok || !parsedToken.Valid {
		ctx.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid token claims"})
		return
	}

	adminID, ok := claims["iss"].(string)
	if !ok {
		ctx.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid admin ID in token"})
		return
	}

	// Define slices to hold the data from each table
	var votes []models.Vote
	var voters []models.Voter
	var candidates []models.Candidate
	var voteCounts []struct {
		Party     string `json:"party"`
		VoteCount int    `json:"vote_count"`
	}

	// Fetch data from the votes table
	query1 := `SELECT * FROM votes`
	res1 := database.DB.Raw(query1).Find(&votes)
	if res1.Error != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": res1.Error.Error()})
		return
	}

	// Fetch data from the voters table
	query2 := `SELECT * FROM voters`
	res2 := database.DB.Raw(query2).Find(&voters)
	if res2.Error != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": res2.Error.Error()})
		return
	}

	// Fetch data from the candidates table
	query3 := `SELECT * FROM candidates`
	res3 := database.DB.Raw(query3).Find(&candidates)
	if res3.Error != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": res3.Error.Error()})
		return
	}

	// Query 4: Aggregate query to get vote count by party
	query4 := `
		SELECT c.Name AS Party, COUNT(v.ID) AS VoteCount
		FROM votes v
		JOIN candidates c ON v.candidate_id = c.ID
		GROUP BY c.Name;
	`
	res4 := database.DB.Raw(query4).Scan(&voteCounts)
	if res4.Error != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": res4.Error.Error()})
		return
	}

	// Retrieve the admin name from the database using the adminID
	var adminName string
	queryAdmin := `SELECT name FROM admins WHERE admin_id = ?`
	resAdmin := database.DB.Raw(queryAdmin, adminID).Scan(&adminName)
	if resAdmin.Error != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": resAdmin.Error.Error()})
		return
	}

	// Grant the admin write privileges using the stored procedure
	grantQuery := fmt.Sprintf(`CALL grant_privileges_to_user('%s');`, adminName)

	// Execute the query to grant the write privileges
	resGrant := database.DB.Exec(grantQuery)
	if resGrant.Error != nil {
		// Log and return the error if the query fails
		log.Println("Error granting privileges:", resGrant.Error)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": resGrant.Error.Error()})
		return
	}

	// Send all data to the frontend
	ctx.JSON(http.StatusOK, gin.H{
		"admin_id":    adminID,
		"admin_name":  adminName,  // Add the admin name to the response
		"votes":       votes,
		"voters":      voters,
		"candidates":  candidates,
		"vote_counts": voteCounts, // Add vote counts to the response
	})
}


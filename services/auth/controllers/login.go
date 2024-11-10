package controllers

import (
	"net/http"
	"strings"
	"github.com/Decentralized-voting-sytem/backend/db/database"
	"github.com/Decentralized-voting-sytem/backend/db/models"
	"time"
	"github.com/golang-jwt/jwt/v5"
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
	if err := c.Bind(&body); err != nil {
		c.JSON(402, gin.H{"error": "Failed to read the body"})
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

	// Create the JWT token
	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"iss": body.VoterID,
		"exp": time.Now().Add(time.Hour * 24).Unix(),
	})

	secretKey := "samyak"
	if secretKey == "" {
		c.JSON(500, gin.H{"error": "Server configuration error"})
		return
	}

	token, err := claims.SignedString([]byte(secretKey))
	if err != nil {
		c.JSON(401, gin.H{"error": "Token generation error"})
		return
	}

	// Set cookie with the generated token
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("Auth", token, 3600*24*30, "", "", false, false)

	// Return success response
	c.JSON(200, gin.H{
		"message": "Login successful",
		"token":   token,
	})
}
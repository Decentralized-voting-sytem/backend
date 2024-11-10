package controllers

import (
		"strings"
		"time"
		"github.com/golang-jwt/jwt/v5"
		"net/http"
		"github.com/Decentralized-voting-sytem/backend/db/database"
		"github.com/Decentralized-voting-sytem/backend/db/models"
		"github.com/gin-gonic/gin"
)
	
func AdminLogin(c *gin.Context) {
	var body struct {
		AdminID  string `json:"admin_id"`
		Name     string `json:"name"`
		Password string `json:"password"`
	}

	// Bind the incoming JSON request body to the body struct
	if err := c.Bind(&body); err != nil {
		c.JSON(402, gin.H{"error": "Failed to read the body"})
		return
	}

	var admin models.Admin

	// Check credentials against the database
	query := `SELECT * FROM admins WHERE admin_id = ? AND password = ?`
	res := database.DB.Raw(query, body.AdminID, body.Password).First(&admin)
	if res.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": res.Error.Error()})
		return
	}

	// Check if the name matches
	if res.RowsAffected == 0 || strings.ToLower(admin.Name) != strings.ToLower(body.Name) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	// Create the JWT token
	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"iss": body.AdminID,
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
	c.SetCookie("AdminAuth", token, 3600*24, "", "", false, false)

	// Return success response
	c.JSON(200, gin.H{
		"message": "Login successful",
		"token":   token,
	})
}
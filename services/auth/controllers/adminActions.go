package controllers

import (
	// "github.com/golang-jwt/jwt/v5"
	"net/http"
	// "github.com/Decentralized-voting-sytem/backend/services/auth/utils"
	"fmt"
	"github.com/Decentralized-voting-sytem/backend/db/database"
	"github.com/gin-gonic/gin"
)

func AdminActions(ctx *gin.Context) {
	// // Retrieve the Auth token from the cookie
	// token, err := ctx.Cookie("AdminAuth")
	// if err != nil {
	// 	ctx.JSON(http.StatusBadRequest, gin.H{"message": "Token not found"})
	// 	return
	// }

	// // Parse the JWT token
	// parsedToken, err := utils.ParseToken(token)
	// if err != nil {
	// 	ctx.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid token"})
	// 	return
	// }

	// // Extract the admin ID from the token claims
	// claims, ok := parsedToken.Claims.(jwt.MapClaims)
	// if !ok || !parsedToken.Valid {
	// 	ctx.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid token claims"})
	// 	return
	// }

	// adminID,ok := claims["iss"].(string)
	// if !ok {
	// 	ctx.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid admin ID in token"})
	// 	return
	// }

	// fmt.Println(adminID)

	var input struct {
		Table   string                 `json:"table"`  // Table to perform action on
		Action  string                 `json:"action"` // Action type: insert, update, delete
		Data    map[string]interface{} `json:"data"`   // Data for insert or update
		Field   string                 `json:"field,omitempty"` // Optional: field to update
	}

	// Parse JSON input
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	var query string
	var args []interface{}

	// Construct the query based on the table and action
	switch input.Table {
	case "voters":
		switch input.Action {
		case "insert":
			query = "INSERT INTO voters (name, dob, voter_id, password, created_at, updated_at) VALUES (?, ?, ?, ?, NOW(), NOW())"
			args = append(args, input.Data["name"], input.Data["dob"], input.Data["voter_id"], input.Data["password"])
		case "update":
			switch input.Field {
			case "name":
				query = "UPDATE voters SET name = ? WHERE voter_id = ?"
				args = append(args, input.Data["value"], input.Data["voter_id"])
			case "password":
				query = "UPDATE voters SET password = ? WHERE voter_id = ?"
				args = append(args, input.Data["value"], input.Data["voter_id"])
			default:
				ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid field for update"})
				return
			}
		case "delete":
			query = "DELETE FROM voters WHERE voter_id = ?"
			args = append(args, input.Data["voter_id"])
		default:
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid action for voters table"})
			return
		}

	case "votes":
		switch input.Action {
		case "insert":
			query = "INSERT INTO votes (voter_id, candidate_id, created_at, updated_at) VALUES (?, ?, NOW(), NOW())"
			args = append(args, input.Data["voter_id"], input.Data["candidate_id"])
		case "update":
			query = "UPDATE votes SET candidate_id = ? WHERE voter_id = ?"
			args = append(args, input.Data["value"], input.Data["voter_id"])
		case "delete":
			query = "DELETE FROM votes WHERE voter_id = ?"
			args = append(args, input.Data["voter_id"])
		default:
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid action for votes table"})
			return
		}

	default:
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid table specified"})
		return
	}

	// Print the query and arguments for debugging
	fmt.Println("Executing query: %s\n", query)
	fmt.Println("With arguments: %v\n", args)

	// Execute the SQL query
	if query != "" {
		res := database.DB.Exec(query, args...)
		if res.Error != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": res.Error.Error()})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{"message": "Action executed successfully"})
	} else {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "No valid action found"})
	}
}

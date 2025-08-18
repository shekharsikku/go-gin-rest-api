package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (app *application) getAllUsers(ctx *gin.Context) {
	users, err := app.models.Users.GetAll()

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve users"})
	}

	ctx.JSON(http.StatusOK, users)
}

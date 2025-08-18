package main

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/shekharsikku/go-gin-rest-api/internal/database"
	"golang.org/x/crypto/bcrypt"
)

type registerRequest struct {
	Name     string `json:"name" binding:"required,min=3"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

type loginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=1"`
}

type loginResponse struct {
	Token  string `json:"token"`
	UserId int    `json:"uid"`
}

func (app *application) registerUser(ctx *gin.Context) {
	var register registerRequest

	if err := ctx.ShouldBindJSON(&register); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(register.Password), bcrypt.DefaultCost)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Something went wrong"})
		return
	}

	register.Password = string(hashed)

	user := database.User{
		Name:     register.Name,
		Email:    register.Email,
		Password: register.Password,
	}

	err = app.models.Users.Insert(&user)

	if err != nil {
		if sqliteErr, ok := err.(interface{ Error() string }); ok &&
			strings.Contains(sqliteErr.Error(), "UNIQUE constraint failed") {
			ctx.JSON(http.StatusConflict, gin.H{"error": "Email already exists"})
			return
		}

		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Could not register user"})
		return
	}

	ctx.JSON(http.StatusCreated, user)
}

func (app *application) loginUser(ctx *gin.Context) {
	var auth loginRequest

	if err := ctx.ShouldBindJSON(&auth); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := app.models.Users.GetByEmail(auth.Email)

	fmt.Println(user)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return
		}
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Something went wrong"})
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(auth.Password))

	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"uid": user.Id,
		"exp": time.Now().Add(time.Hour * 72).Unix(),
	})

	tokenString, err := token.SignedString([]byte(app.jwtSecret))

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error generating token"})
		return
	}

	ctx.JSON(http.StatusOK, loginResponse{Token: tokenString, UserId: user.Id})
}

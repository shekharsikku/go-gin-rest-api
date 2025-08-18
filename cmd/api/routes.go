package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (app *application) routes() http.Handler {
	router := gin.Default()

	router.GET("/hello", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{"message": "Hello, World!"})
	})

	v1 := router.Group("/api/v1")

	{
		v1.POST("/events", app.createEvent)
		v1.GET("/events", app.getEvents)
		v1.GET("/events/:id", app.getEvent)
	}

	return router
}

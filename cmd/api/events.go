package main

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/shekharsikku/go-gin-rest-api/internal/database"
)

func (app *application) createEvent(ctx *gin.Context) {
	var event database.Event

	if err := ctx.ShouldBindJSON(&event); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user := app.GetUserFromContext(ctx)
	event.Owner = user.Id
	err := app.models.Events.Insert(&event)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create event"})
		return
	}

	ctx.JSON(http.StatusCreated, event)
}

func (app *application) getEvents(ctx *gin.Context) {
	events, err := app.models.Events.GetAll()

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve events"})
		return
	}

	ctx.JSON(http.StatusOK, events)
}

func (app *application) getEvent(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid event id"})
		return
	}

	event, err := app.models.Events.Get(id)

	if event == nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Event not found"})
		return
	}

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve event"})
		return
	}

	ctx.JSON(http.StatusOK, event)
}

func (app *application) updateEvent(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid event id"})
		return
	}

	user := app.GetUserFromContext(ctx)
	existingEvent, err := app.models.Events.Get(id)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve event"})
		return
	}

	if existingEvent == nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Event not found"})
		return
	}

	if existingEvent.Owner != user.Id {
		ctx.JSON(http.StatusForbidden, gin.H{"error": "You are not authorized to update this event"})
		return
	}

	updatedEvent := &database.Event{}

	if err := ctx.ShouldBindJSON(updatedEvent); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updatedEvent.Id = id

	if err := app.models.Events.Update(updatedEvent); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update event"})
		return
	}

	updatedEvent.Owner = user.Id
	ctx.JSON(http.StatusOK, updatedEvent)
}

func (app *application) deleteEvent(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid event id"})
		return
	}

	user := app.GetUserFromContext(ctx)
	existingEvent, err := app.models.Events.Get(id)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"Error": "Failed to retrieve event"})
		return
	}

	if existingEvent == nil {
		ctx.JSON(http.StatusNotFound, gin.H{"Error": "Event not found"})
		return
	}

	if existingEvent.Owner != user.Id {
		ctx.JSON(http.StatusForbidden, gin.H{"error": "You are not authorized to delete this event"})
		return
	}

	if err := app.models.Events.Delete(id); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete event"})
	}

	ctx.JSON(http.StatusNoContent, nil)
}

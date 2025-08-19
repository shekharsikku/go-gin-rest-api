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

/** Routes Handles for Events Attendees */

func (app *application) addAttendeeToEvent(ctx *gin.Context) {
	eventId, err := strconv.Atoi(ctx.Param("id"))

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid event id"})
		return
	}

	userId, err := strconv.Atoi(ctx.Param("uid"))

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user id"})
		return
	}

	event, err := app.models.Events.Get(eventId)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve event"})
		return
	}

	if event == nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Event not found"})
		return
	}

	userToAdd, err := app.models.Users.Get(userId)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve user"})
		return
	}

	if userToAdd == nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	user := app.GetUserFromContext(ctx)

	if event.Owner != user.Id {
		ctx.JSON(http.StatusForbidden, gin.H{"error": "You are not authorized to add an attendee"})
		return
	}

	existingAttendee, err := app.models.Attendees.GetByEventAndAttendee(event.Id, userToAdd.Id)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve attendee"})
		return
	}

	if existingAttendee != nil {
		ctx.JSON(http.StatusConflict, gin.H{"error": "Attendee already exists"})
		return
	}

	attendee := database.Attendee{
		EventId: event.Id,
		UserId:  userToAdd.Id,
	}

	_, err = app.models.Attendees.Insert(&attendee)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add attendee"})
		return
	}

	ctx.JSON(http.StatusCreated, attendee)
}

func (app *application) deleteAttendeeFromEvent(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid event id"})
		return
	}

	userId, err := strconv.Atoi(ctx.Param("uid"))

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user id"})
		return
	}

	event, err := app.models.Events.Get(id)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Something went wrong"})
		return
	}

	if event == nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Event not found"})
		return
	}

	user := app.GetUserFromContext(ctx)

	if event.Owner != user.Id {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "You are not authorized to delete an attendee from event"})
		return
	}

	err = app.models.Attendees.Delete(userId, id)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete attendee"})
		return
	}

	ctx.JSON(http.StatusNoContent, nil)
}

func (app *application) getAttendeesForEvent(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid event id"})
		return
	}

	users, err := app.models.Attendees.GetAttendeesByEvent(id)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to to retrieve attendees for events"})
		return
	}

	ctx.JSON(http.StatusOK, users)
}

func (app *application) getEventsByAttendee(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid attendee id"})
		return
	}

	events, err := app.models.Attendees.GetEventsByAttendee(id)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get events"})
		return
	}

	ctx.JSON(http.StatusOK, events)
}

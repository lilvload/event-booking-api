package routes

import (
	"database/sql"
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/lilvload/event-booking-api/models"
)

// getEvents godoc
// @Summary      List events
// @Description  Returns all events. An empty list is returned as [].
// @Tags         events
// @Produce      json
// @Success      200  {array}   models.Event
// @Failure      500  {object}  map[string]string  "Could not fetch events"
// @Router       /events [get]
func getEvents(context *gin.Context) {
	events, err := models.GetAllEvents()

	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Could not fetch events. Try again later."})
		return
	}

	context.JSON(http.StatusOK, events)
}

// getEvent godoc
// @Summary      Get an event
// @Tags         events
// @Produce      json
// @Param        id   path      int  true  "Event ID"
// @Success      200  {object}  models.Event
// @Failure      400  {object}  map[string]string  "The ID is not a number"
// @Failure      404  {object}  map[string]string  "Event not found"
// @Failure      500  {object}  map[string]string  "Could not fetch the event"
// @Router       /events/{id} [get]
func getEvent(context *gin.Context) {
	eventId, err := strconv.ParseInt(context.Param("id"), 10, 64)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Could not parse event id."})
		return
	}

	event, err := models.GetEventByID(eventId)

	if errors.Is(err, sql.ErrNoRows) {
		context.JSON(http.StatusNotFound, gin.H{"message": "Event not found."})
		return
	}
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Could not fetch event. Try again later."})
		return
	}

	context.JSON(http.StatusOK, event)
}

// createEvent godoc
// @Summary      Create an event
// @Description  Creates an event owned by the current user. All four fields are required; dateTime is an RFC 3339 timestamp.
// @Tags         events
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        event  body      models.Event            true  "Event data (id and userId are ignored)"
// @Success      201    {object}  map[string]interface{}  "message and the created event"
// @Failure      400    {object}  map[string]string       "Invalid JSON or a missing field"
// @Failure      401    {object}  map[string]string       "No valid token"
// @Failure      500    {object}  map[string]string       "Could not create the event"
// @Router       /events [post]
func createEvent(context *gin.Context) {
	var event models.Event
	err := context.ShouldBindJSON(&event)

	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Could not parse request data."})
		return
	}

	userId := context.GetInt64("userId")
	event.UserID = userId

	err = event.Save()

	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Could not create event. Try again later."})
		return
	}

	context.JSON(http.StatusCreated, gin.H{"message": "Event created!", "event": event})
}

// updateEvent godoc
// @Summary      Update an event
// @Description  Replaces all fields of an event. Only the creator of the event can update it.
// @Tags         events
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id     path      int                true  "Event ID"
// @Param        event  body      models.Event       true  "New event data (id and userId are ignored)"
// @Success      200    {object}  map[string]string  "message"
// @Failure      400    {object}  map[string]string  "Invalid ID, invalid JSON or a missing field"
// @Failure      401    {object}  map[string]string  "No valid token"
// @Failure      403    {object}  map[string]string  "The current user is not the creator"
// @Failure      404    {object}  map[string]string  "Event not found"
// @Failure      500    {object}  map[string]string  "Could not fetch or update the event"
// @Router       /events/{id} [put]
func updateEvent(context *gin.Context) {
	eventId, err := strconv.ParseInt(context.Param("id"), 10, 64)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Could not parse event id."})
		return
	}

	userId := context.GetInt64("userId")
	event, err := models.GetEventByID(eventId)

	if errors.Is(err, sql.ErrNoRows) {
		context.JSON(http.StatusNotFound, gin.H{"message": "Event not found."})
		return
	}
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Could not fetch the event."})
		return
	}

	if event.UserID != userId {
		context.JSON(http.StatusForbidden, gin.H{"message": "You have no permission to update this event."})
		return
	}

	var updatedEvent models.Event
	err = context.ShouldBindJSON(&updatedEvent)

	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Could not parse request data."})
		return
	}

	updatedEvent.ID = eventId
	err = updatedEvent.Update()

	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Could not update event."})
		return
	}

	context.JSON(http.StatusOK, gin.H{"message": "Event updated successfully!"})
}

// deleteEvent godoc
// @Summary      Delete an event
// @Description  Deletes an event and all registrations for it. Only the creator of the event can delete it.
// @Tags         events
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      int                true  "Event ID"
// @Success      200  {object}  map[string]string  "message"
// @Failure      400  {object}  map[string]string  "Invalid ID"
// @Failure      401  {object}  map[string]string  "No valid token"
// @Failure      403  {object}  map[string]string  "The current user is not the creator"
// @Failure      404  {object}  map[string]string  "Event not found"
// @Failure      500  {object}  map[string]string  "Could not fetch or delete the event"
// @Router       /events/{id} [delete]
func deleteEvent(context *gin.Context) {
	eventId, err := strconv.ParseInt(context.Param("id"), 10, 64)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Could not parse event id."})
		return
	}

	userId := context.GetInt64("userId")
	event, err := models.GetEventByID(eventId)

	if errors.Is(err, sql.ErrNoRows) {
		context.JSON(http.StatusNotFound, gin.H{"message": "Event not found."})
		return
	}
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Could not fetch the event."})
		return
	}

	if event.UserID != userId {
		context.JSON(http.StatusForbidden, gin.H{"message": "You have no permission to delete this event."})
		return
	}

	err = event.Delete()
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Could not delete the event."})
		return
	}

	context.JSON(http.StatusOK, gin.H{"message": "Event deleted successfully!"})
}

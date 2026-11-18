package routes

import (
	"database/sql"
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/lilvload/event-booking-api/models"
)

// registerForEvent godoc
// @Summary      Register for an event
// @Description  Registers the current user for the event. No request body.
// @Tags         registrations
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      int                true  "Event ID"
// @Success      201  {object}  map[string]string  "message"
// @Failure      400  {object}  map[string]string  "Invalid ID"
// @Failure      401  {object}  map[string]string  "No valid token"
// @Failure      404  {object}  map[string]string  "Event not found"
// @Failure      409  {object}  map[string]string  "Already registered for this event"
// @Failure      500  {object}  map[string]string  "Could not save the registration"
// @Router       /events/{id}/register [post]
func registerForEvent(context *gin.Context) {
	userId := context.GetInt64("userId")
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
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Could not fetch event."})
		return
	}

	err = event.Register(userId)
	if err != nil {
		if strings.Contains(err.Error(), "UNIQUE constraint failed") {
			context.JSON(http.StatusConflict, gin.H{"message": "Already registered for this event."})
			return
		}
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Could not register user for event."})
		return
	}

	context.JSON(http.StatusCreated, gin.H{"message": "Registered!"})
}

// cancelRegistration godoc
// @Summary      Cancel a registration
// @Description  Cancels the current user's registration for the event. No request body.
// @Tags         registrations
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      int                true  "Event ID"
// @Success      200  {object}  map[string]string  "message"
// @Failure      400  {object}  map[string]string  "Invalid ID"
// @Failure      401  {object}  map[string]string  "No valid token"
// @Failure      404  {object}  map[string]string  "The user is not registered for this event"
// @Failure      500  {object}  map[string]string  "Could not cancel the registration"
// @Router       /events/{id}/register [delete]
func cancelRegistration(context *gin.Context) {
	userId := context.GetInt64("userId")
	eventId, err := strconv.ParseInt(context.Param("id"), 10, 64)

	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Could not parse event id."})
		return
	}

	var event models.Event
	event.ID = eventId

	err = event.CancelRegistration(userId)

	if errors.Is(err, sql.ErrNoRows) {
		context.JSON(http.StatusNotFound, gin.H{"message": "Registration not found."})
		return
	}
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Could not cancel registration."})
		return
	}

	context.JSON(http.StatusOK, gin.H{"message": "Cancelled!"})
}

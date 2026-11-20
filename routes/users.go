package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lilvload/event-booking-api/db"
	"github.com/lilvload/event-booking-api/models"
	"github.com/lilvload/event-booking-api/utils"
)

// signup godoc
// @Summary      Create a user
// @Description  Registers a new user. The password is stored as a bcrypt hash. No token is returned; call POST /login afterwards.
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        user  body      models.SignupRequest        true  "Email and password"
// @Success      201   {object}  map[string]string  "message"
// @Failure      400   {object}  map[string]string  "Invalid JSON, a missing field or validation error"
// @Failure      409   {object}  map[string]string  "Email already registered"
// @Failure      500   {object}  map[string]string  "Could not save the user"
// @Router       /signup [post]
func signup(context *gin.Context) {
	var req models.SignupRequest

	err := context.ShouldBindJSON(&req)

	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": validationMessage(err)})
		return
	}

	user := models.User{
		Email:    req.Email,
		Password: req.Password,
	}

	err = user.Save()
	if err != nil {
		if db.IsUniqueViolation(err) {
			context.JSON(http.StatusConflict, gin.H{"message": "Email already registered."})
			return
		}
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Could not save user."})
		return
	}

	context.JSON(http.StatusCreated, gin.H{"message": "User created successfully!"})
}

// login godoc
// @Summary      Log in
// @Description  Checks the credentials and returns a JWT that is valid for 2 hours.
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        user  body      models.LoginRequest        true  "Email and password"
// @Success      200   {object}  map[string]string  "message and token"
// @Failure      400   {object}  map[string]string  "Invalid JSON or a missing field"
// @Failure      401   {object}  map[string]string  "Credentials invalid"
// @Failure      500   {object}  map[string]string  "Could not generate the token"
// @Router       /login [post]
func login(context *gin.Context) {
	var req models.LoginRequest

	err := context.ShouldBindJSON(&req)

	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": validationMessage(err)})
		return
	}

	user := models.User{
		Email:    req.Email,
		Password: req.Password,
	}

	err = user.ValidateCredentials()

	if err != nil {
		context.JSON(http.StatusUnauthorized, gin.H{"message": err.Error()})
		return
	}

	token, err := utils.GenerateToken(user.Email, user.ID)

	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Could not authenticate user."})
		return
	}

	context.JSON(http.StatusOK, gin.H{"message": "Login successful!", "token": token})
}

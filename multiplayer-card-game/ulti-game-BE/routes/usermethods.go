package routes

import (
	"errors"
	"log"
	"net/http"
	"strings"

	"exmaple.com/ulti-restapi/models"
	"exmaple.com/ulti-restapi/utility"
	"github.com/gin-gonic/gin"
)

func authenticationMethod(context *gin.Context) {
	var user models.User
	err := context.ShouldBindJSON(&user)

	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Could not parse user data the first occurence"})
	}

	if user.State == "signup" {
		err = signup(user)
		if err != nil {
			context.JSON(http.StatusInternalServerError, gin.H{"message": "Couldnot save user"})
			return
		}
		context.JSON(http.StatusCreated, gin.H{"message": "User saved.", "user": user})

	} else if user.State == "login" {

		token, err := login(user)
		if err != nil {
			context.JSON(http.StatusUnauthorized, gin.H{"message": err.Error()})
			return
		}

		context.JSON(http.StatusOK, gin.H{"message": "Logged in", "token": token})

	} else {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Unrecognisable data."})
	}
}

func signup(user models.User) error {

	err := user.Save()
	if err != nil {
		err = errors.New("credentials is invalid, at signup")
		return err
	}
	return err

}

func login(user models.User) (string, error) {

	err := user.ValidateUser()

	if err != nil {
		err = errors.New("credentials is invalid at login")
		return "", err
	}

	if user.IsLoggedIn {
		log.Println("User is already logged in")
		return "", errors.New("user already logged in from another device")
	}

	// mark as logged in
	if err := user.SetLoggedIn(true); err != nil {
		return "", errors.New("could not update login state")
	}

	token, err := utility.GenerateToken(user.Email, user.ID)

	if err != nil {
		err = errors.New("credentials is invalid didn't get token")
		return "", err
	}

	return token, nil
}

func logout(context *gin.Context) {
	header := context.GetHeader("Authorization")

	if header == "" {
		context.JSON(http.StatusUnauthorized, gin.H{"message": "missing Authorization header"})
	}

	parts := strings.SplitN(header, " ", 2)
	if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
		context.JSON(http.StatusUnauthorized, gin.H{"message": "invalid Authorization header"})
		return
	}

	token := parts[1]

	userID, err := utility.VerifyToken(token)
	if err != nil {
		context.JSON(http.StatusUnauthorized, gin.H{"message": "invalid or expired token"})
		return
	}

	user, err := models.GetUserByID(userID)
	if err != nil {
		log.Println("GetUserByID error:", err)
		context.JSON(http.StatusInternalServerError, gin.H{"message": "user not found"})
		return
	}

	if err := user.SetLoggedIn(false); err != nil {
		log.Println("SetLoggedIn error:", err)

		context.JSON(http.StatusInternalServerError, gin.H{"message": "could not update logout state"})
		return
	}

	context.JSON(http.StatusOK, gin.H{"message": "logged out"})

}

func getUsers(context *gin.Context) {
	users, err := models.GetAllUsers()
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Could not fetch"})
		return
	}
	context.JSON(http.StatusOK, users)
}

func getProfile(context *gin.Context) {
	header := context.GetHeader("Authorization")

	if header == "" {
		context.JSON(http.StatusUnauthorized, gin.H{"message": "missing Authorization header"})
		return
	}

	parts := strings.SplitN(header, " ", 2)
	if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
		context.JSON(http.StatusUnauthorized, gin.H{"message": "invalid Authorization header"})
		return
	}

	token := parts[1]
	userID, err := utility.VerifyToken(token)
	if err != nil {
		context.JSON(http.StatusUnauthorized, gin.H{"message": "invalid or expired token"})
		return
	}

	user, err := models.GetUserByID(userID)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "user not found"})
		return
	}

	context.JSON(http.StatusOK, gin.H{
		"id":    user.ID,
		"name":  user.Name,
		"email": user.Email,
	})
}

func changeEmail(context *gin.Context) {
	header := context.GetHeader("Authorization")

	if header == "" {
		context.JSON(http.StatusUnauthorized, gin.H{"message": "missing Authorization header"})
		return
	}

	parts := strings.SplitN(header, " ", 2)
	if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
		context.JSON(http.StatusUnauthorized, gin.H{"message": "invalid Authorization header"})
		return
	}

	token := parts[1]
	userID, err := utility.VerifyToken(token)
	if err != nil {
		context.JSON(http.StatusUnauthorized, gin.H{"message": "invalid or expired token"})
		return
	}

	var requestData struct {
		NewEmail string `json:"newEmail" binding:"required,email"`
		Password string `json:"password" binding:"required"`
	}

	if err := context.ShouldBindJSON(&requestData); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Invalid request data"})
		return
	}

	user, err := models.GetUserByID(userID)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "user not found"})
		return
	}

	// Verify current password
	if !utility.CheckPasswordHash(requestData.Password, user.Password) {
		context.JSON(http.StatusUnauthorized, gin.H{"message": "incorrect password"})
		return
	}

	// Update email
	if err := user.UpdateEmail(requestData.NewEmail); err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "could not update email"})
		return
	}

	context.JSON(http.StatusOK, gin.H{"message": "email updated successfully"})
}

func changePassword(context *gin.Context) {
	header := context.GetHeader("Authorization")

	if header == "" {
		context.JSON(http.StatusUnauthorized, gin.H{"message": "missing Authorization header"})
		return
	}

	parts := strings.SplitN(header, " ", 2)
	if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
		context.JSON(http.StatusUnauthorized, gin.H{"message": "invalid Authorization header"})
		return
	}

	token := parts[1]
	userID, err := utility.VerifyToken(token)
	if err != nil {
		context.JSON(http.StatusUnauthorized, gin.H{"message": "invalid or expired token"})
		return
	}

	var requestData struct {
		CurrentPassword string `json:"currentPassword" binding:"required"`
		NewPassword     string `json:"newPassword" binding:"required"`
	}

	if err := context.ShouldBindJSON(&requestData); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Invalid request data"})
		return
	}

	user, err := models.GetUserByID(userID)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "user not found"})
		return
	}

	// Verify current password
	if !utility.CheckPasswordHash(requestData.CurrentPassword, user.Password) {
		context.JSON(http.StatusUnauthorized, gin.H{"message": "incorrect current password"})
		return
	}

	// Hash new password
	hashedPassword, err := utility.HashPassword(requestData.NewPassword)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "could not hash password"})
		return
	}

	// Update password
	if err := user.UpdatePassword(hashedPassword); err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "could not update password"})
		return
	}

	context.JSON(http.StatusOK, gin.H{"message": "password updated successfully"})
}

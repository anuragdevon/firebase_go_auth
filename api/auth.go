package api

import (
	"log"
	"net/http"

	"firebase.google.com/go/auth"
	"github.com/gin-gonic/gin"

	firebase_conn "firebase_go_auth/firebase_conn"
)

type NewUser struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`

	// Should be given as base64 encoded image
	// Image   string `json:"image"`
}

type LoginUser struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func UserSignUp(c *gin.Context) {

	// Extract Input
	var input NewUser
	if err := c.ShouldBindJSON(&input); err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"response": "Invalid Input!"})
		return
	}

	ctx, client, err := firebase_conn.FirebaseInit()
	if err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"response": "Internal Server Error!"})
		return
	}

	// Check if user already exists
	_, err = client.GetUserByEmail(ctx, input.Email)
	if err == nil {
		log.Println(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"response": "User already exists!"})
		return
	}

	defaultPhotoUrl := "data:image/jpeg;base64,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"

	params := (&auth.UserToCreate{}).
		Email(input.Email).
		EmailVerified(false).
		Password(input.Password).
		DisplayName(input.Name).
		PhotoURL(defaultPhotoUrl).
		Disabled(false)

	_, err = client.CreateUser(ctx, params)

	if err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"response": "Unable to create user!"})
		return
	}

	// Send Email Verification
	err = firebase_conn.EmailVerification(input.Email, client, ctx)
	if err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"response": "Verification Email Sending Failed!"})
		return

	}

	c.JSON(http.StatusOK, gin.H{"response": "User Created Successfully!"})
}

func UserSignIn(c *gin.Context) {

	// Extract Input
	var input NewUser
	if err := c.ShouldBindJSON(&input); err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"response": "Invalid Input!"})
		return
	}

	ctx, client, err := firebase_conn.FirebaseInit()
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"response": "Internal Server Error!"})
		return
	}

	// Check if user exists
	user, err := client.GetUserByEmail(ctx, input.Email)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"response": "No Such User Exists!"})
		return
	}

	// Check if user is verified
	if !user.EmailVerified {
		c.JSON(http.StatusBadRequest, gin.H{"response": "Email Not Verified!"})
		return
	}

	resp, err := firebase_conn.SignInWithEmailPassword(input.Email, input.Password)
	if err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"response": "Internal Server Error!"})
		return
	}
	log.Println("\nLOGIN DATA: ", resp.Body)

	// Check

	if err != nil {
		c.JSON(http.StatusOK, gin.H{"response": "User Created Successfully!"})
	}
}

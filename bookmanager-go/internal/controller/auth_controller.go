package controller

import (
	"bookmanager-go/internal/model"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AuthHTMLController struct {
	DB *gorm.DB
}

// ShowRegisterPage displays the registration form
func (ac *AuthHTMLController) ShowRegisterPage(c *gin.Context) {
	bc := BookController{DB: ac.DB}
	bc.renderHTML(c, http.StatusOK, "register.html", gin.H{
		"Title":       "Register",
		"PageTitle":   "Register",
		"Description": "Create a new user account.",
	})
}

// RegisterUser handles the registration form submission
func (ac *AuthHTMLController) RegisterUser(c *gin.Context) {
	email := c.PostForm("email")
	password := c.PostForm("password")

	// Validate input fields
	if email == "" || password == "" {
		bc := BookController{DB: ac.DB} // HTML renderer
		bc.renderHTML(c, http.StatusBadRequest, "register.html", gin.H{
			"PageTitle":   "Register",
			"Message":     "Email and password are required.",
			"MessageType": "warning",
		})
		return
	}

	// Hash password
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("Error hashing password: %v\n", err)
		bc := BookController{DB: ac.DB} // HTML renderer
		bc.renderHTML(c, http.StatusInternalServerError, "register.html", gin.H{
			"PageTitle":   "Register",
			"Message":     "Internal error.",
			"MessageType": "danger",
		})
		return
	}

	// Create user struct
	user := model.User{
		Email:    email,
		Password: string(hashed),
	}

	// Save user to database
	if err := ac.DB.Create(&user).Error; err != nil {
		log.Printf("Error creating user: %v\n", err)
		bc := BookController{DB: ac.DB} // HTML renderer
		bc.renderHTML(c, http.StatusBadRequest, "register.html", gin.H{
			"PageTitle":   "Register",
			"Message":     "User already exists.",
			"MessageType": "warning",
		})
		return
	}

	// Redirect to login page on success
	c.Redirect(http.StatusSeeOther, "/login")
}

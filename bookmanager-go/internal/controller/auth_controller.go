package controller

import (
	"bookmanager-go/internal/auth"
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
			"Title":       "Register",
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
			"Title":       "Register",
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
			"Title":       "Register",
			"Message":     "User already exists.",
			"MessageType": "warning",
		})
		return
	}

	// Redirect to login page on success
	c.Redirect(http.StatusSeeOther, "/login")
}

// ShowLoginPage displays the login form
func (ac *AuthHTMLController) ShowLoginPage(c *gin.Context) {
	bc := BookController{DB: ac.DB}

	// Render login HTML page
	bc.renderHTML(c, http.StatusOK, "login.html", gin.H{
		"Title":       "Login",
		"PageTitle":   "Login",
		"Description": "Enter your email and password.",
	})
}

// LoginUser handles login form submission (step 1: input + user + password)
func (ac *AuthHTMLController) LoginUser(c *gin.Context) {
	email := c.PostForm("email")
	password := c.PostForm("password")

	bc := BookController{DB: ac.DB} // for renderHTML()

	// Validate input
	if email == "" || password == "" {
		bc.renderHTML(c, http.StatusBadRequest, "login.html", gin.H{
			"Title":       "Login",
			"PageTitle":   "Login",
			"Message":     "Email and password are required.",
			"MessageType": "warning",
		})
		return
	}

	// Look up user
	var user model.User
	if err := ac.DB.Where("email = ?", email).First(&user).Error; err != nil {
		log.Printf("Login error: user not found (%s)\n", email)
		bc.renderHTML(c, http.StatusUnauthorized, "login.html", gin.H{
			"Title":       "Login",
			"PageTitle":   "Login",
			"Message":     "User not found.",
			"MessageType": "danger",
		})
		return
	}

	// Validate password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		log.Printf("Login error: invalid password for user %s\n", email)
		bc.renderHTML(c, http.StatusUnauthorized, "login.html", gin.H{
			"Title":       "Login",
			"PageTitle":   "Login",
			"Message":     "Invalid password.",
			"MessageType": "danger",
		})
		return
	}
	// Create JWT token
	token, err := auth.CreateToken(user.ID)
	if err != nil {
		log.Printf("Error creating JWT: %v\n", err)
		bc.renderHTML(c, http.StatusInternalServerError, "login.html", gin.H{
			"Title":       "Login",
			"PageTitle":   "Login",
			"Message":     "Internal server error.",
			"MessageType": "danger",
		})
		return
	}

	// Set JWT cookie (1 day)
	c.SetCookie(
		"jwt",
		token,
		3600*24,
		"/",
		"",
		false, // not HTTPS only
		true,  // httpOnly
	)

	// Success → redirect to book list
	c.Redirect(http.StatusSeeOther, "/books/list")
}

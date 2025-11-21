package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
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

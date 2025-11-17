package controller

import "github.com/gin-gonic/gin"

// renderHTML is a helper function that renders an HTML template
func (bc *BookController) renderHTML(c *gin.Context, status int, templateName string, data gin.H) {
	// Ensure the map is initialized
	if data == nil {
		data = gin.H{}
	}

	// Global default values
	if _, exists := data["Title"]; !exists {
		data["Title"] = "BookManager"
	}
	if _, exists := data["PageTitle"]; !exists {
		if status >= 400 {
			data["PageTitle"] = "Error Situation"
		} else {
			data["PageTitle"] = "BookManager"
		}
	}
	if _, exists := data["Description"]; !exists {
		if status >= 400 {
			data["Description"] = "An unexpected error occurred. Please try again later."
		} else {
			data["Description"] = "Manage your personal library — add, edit, and organize your favorite books."
		}
	}
	if _, exists := data["MessageType"]; !exists {
		data["MessageType"] = ""
	}
	if _, exists := data["Message"]; !exists {
		data["Message"] = ""
	}

	// Render page
	c.HTML(status, templateName, data)
}

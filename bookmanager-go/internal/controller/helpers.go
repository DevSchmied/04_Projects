package controller

import (
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
)

// parseIDParam extracts and validates the "id" URL parameter.
func parseIDParam(c *gin.Context) (uint, error) {
	idParam := c.Param("id")
	if idParam == "" {
		return 0, fmt.Errorf("missing ID parameter")
	}

	id, err := strconv.Atoi(idParam)
	if err != nil || id <= 0 {
		return 0, fmt.Errorf("invalid ID: %s", idParam)
	}
	return uint(id), nil
}

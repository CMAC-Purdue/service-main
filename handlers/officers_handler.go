package handlers

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"service-main/db"
	"service-main/util"
)

type createOfficerRequest struct {
	Name          string  `json:"name" binding:"required"`
	LinkedinPhoto *string `json:"linkedin_photo"`
	ImageURI      *string `json:"image_uri"`
}

type officerResponse struct {
	ID            int64   `json:"id"`
	Name          string  `json:"name"`
	LinkedinPhoto *string `json:"linkedin_photo"`
	ImageURI      *string `json:"image_uri"`
}

func CreateOfficerHandler(queries *db.Queries) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req createOfficerRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		officer, err := queries.CreateOfficer(c.Request.Context(), db.CreateOfficerParams{
			Name:          strings.TrimSpace(req.Name),
			LinkedinPhoto: util.ToNullString(req.LinkedinPhoto),
			ImageUri:      util.ToNullString(req.ImageURI),
		})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create officer"})
			return
		}

		c.JSON(http.StatusCreated, officerResponse{
			ID:            officer.ID,
			Name:          officer.Name,
			LinkedinPhoto: util.NullStringToPointer(officer.LinkedinPhoto),
			ImageURI:      util.NullStringToPointer(officer.ImageUri),
		})
	}
}

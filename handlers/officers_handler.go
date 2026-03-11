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

type errorResponse struct {
	Error string `json:"error"`
}

// CreateOfficerHandler godoc
// @Summary Create officer
// @Description Create a new officer record.
// @Tags officers
// @Accept json
// @Produce json
// @Param payload body createOfficerRequest true "Officer payload"
// @Success 201 {object} officerResponse
// @Failure 400 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Router /officers [post]
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
			c.JSON(http.StatusInternalServerError, errorResponse{Error: "Failed to create new officer"})
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

type getOfficersResponse struct {
	Officers []officerResponse `json:"officers"`
}

// GetOfficersHandler godoc
// @Summary List all officers
// @Description Returns a list of every officer
// @Tags officers
// @Produce json
// @Success 200 {object} getOfficersResponse
// @Failure 500 {object} errorResponse
// @Router /officers [get]
func GetOfficersHandler(queries *db.Queries) gin.HandlerFunc {
	return func(c *gin.Context) {
		officers, err := queries.ListOfficers(c.Request.Context())

		fixedOfficers := make([]officerResponse, len(officers))

		for i, v := range officers {
			fixedOfficers[i] = officerResponse{
				ID:            v.ID,
				Name:          v.Name,
				LinkedinPhoto: util.NullStringToPointer(v.LinkedinPhoto),
				ImageURI:      util.NullStringToPointer(v.ImageUri),
			}
		}

		if err != nil {
			c.JSON(http.StatusInternalServerError, errorResponse{Error: err.Error()})
			return
		}

		c.JSON(http.StatusOK, getOfficersResponse{
			Officers: fixedOfficers,
		})
	}
}

package handlers

import (
	"net/http"
	"service-main/auth"

	"github.com/gin-gonic/gin"
)

type AdminSessionRequest struct {
	Password string `json:"password" binding:"required"`
}

type messageResponse struct {
	Message string `json:"message"`
}

// DisplaySessionsHandler godoc
// @Summary List active sessions
// @Description Returns all active sessions. Requires an authenticated `sessid` cookie; this cannot be set via Swagger Authorize.
// @Tags admin
// @Produce json
// @Success 200 {object} map[string]string
// @Failure 401 {object} errorResponse
// @Router /auth/sessions [get]
func DisplaySessionsHandler(store *auth.SessionStore) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, store.DisplaySessions())
	}
}

// AdminSessionLogin godoc
// @Summary Admin login
// @Description Authenticate with the admin password and receive a session cookie
// @Tags admin
// @Accept json
// @Produce json
// @Param payload body AdminSessionRequest true "Admin credentials"
// @Success 200 {object} messageResponse
// @Failure 400 {object} errorResponse
// @Failure 401 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Router /opme [post]
func AdminSessionLogin(store *auth.SessionStore, password string) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req AdminSessionRequest
		err := c.ShouldBindJSON(&req)

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "bad request"})
			return
		}

		if password != req.Password {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "stop trying please"})
			return
		}

		session := auth.NewSession("Admin")
		err = store.AddSessionWithCtx(c, session)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Success!"})
	}
}

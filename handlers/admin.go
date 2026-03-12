package handlers

import (
	"net/http"
	"service-main/auth"

	"github.com/gin-gonic/gin"
)

type AdminSessionRequest struct {
	Password string `json:"password" binding:"required"`
}

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

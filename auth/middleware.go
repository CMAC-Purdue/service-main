package auth

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (store *SessionStore) SessionGuard() gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.GetHeader("sessid")

		if header == "" {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		session, exists := store.GetSession(header)

		if !exists {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		c.Set("session", session)
	}
}

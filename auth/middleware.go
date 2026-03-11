package auth

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (store *SessionStore) SessionGuard() gin.HandlerFunc {
	return func(c *gin.Context) {
		auth, err := c.Cookie("sessid")

		if err != nil {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		if auth == "" {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		session, exists := store.GetSession(auth)

		if !exists {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		c.Set("session", session)
	}
}

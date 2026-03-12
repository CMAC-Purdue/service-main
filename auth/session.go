package auth

import (
	"net/http"
	"service-main/util"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

const sessionTTL = 30 * time.Minute

type Session struct {
	User   string
	Expiry time.Time
}

type SessionStore struct {
	Sessions map[string]Session
	mu       sync.RWMutex
}

func (s *SessionStore) AddSession(sess Session) (string, error) {
	id, err := util.GenerateSessionID()

	if err != nil {
		return "", err
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	s.Sessions[id] = sess

	return id, nil
}

func NewSession(user string) Session {
	return Session{User: user, Expiry: newExpiry()}
}

func newExpiry() time.Time {
	return time.Now().UTC().Add(sessionTTL)
}

func isExpired(expiresAt time.Time) bool {
	return !time.Now().UTC().Before(expiresAt)
}

func (s *SessionStore) SessionCleanJob() {
	for {
		time.Sleep(5 * time.Minute)

		s.mu.Lock()
		for k, v := range s.Sessions {
			if isExpired(v.Expiry) {
				delete(s.Sessions, k)
			}
		}
		s.mu.Unlock()
	}
}

func (s *SessionStore) GetSession(id string) (Session, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	sess, exists := s.Sessions[id]

	return sess, exists
}

func (s *SessionStore) AddSessionWithCtx(c *gin.Context, session Session) error {
	id, err := s.AddSession(session)

	if err != nil {
		return err
	}

	http.SetCookie(c.Writer, &http.Cookie{
		Name:     "sessid",
		Value:    id,
		MaxAge:   int(sessionTTL.Seconds()),
		Path:     "/",
		Secure:   gin.Mode() == gin.ReleaseMode,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	})

	return nil
}

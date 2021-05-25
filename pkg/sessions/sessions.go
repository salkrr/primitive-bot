// Package sessions implements types for working with telegram sessions.
package sessions

import (
	"sync"

	"github.com/lazy-void/primitive-bot/pkg/menu"

	"github.com/lazy-void/primitive-bot/pkg/primitive"
	"github.com/lazy-void/primitive-bot/pkg/tg"
)

// Session represents one telegram session.
type Session struct {
	UserID        int64
	MenuMessageID int64
	InChan        chan<- tg.Message
	ImgPath       string
	Menu          menu.Menu
	Config        primitive.Config
}

// NewSession initializes new instance of Session object.
func NewSession(userID, menuMessageID int64, imgPath string) Session {
	c := primitive.NewConfig()

	return Session{
		UserID:        userID,
		MenuMessageID: menuMessageID,
		ImgPath:       imgPath,
		Menu:          menu.New(c),
		Config:        c,
	}
}

// ActiveSessions represents list of all active telegram sessions.
type ActiveSessions struct {
	data map[int64]Session
	mu   sync.Mutex
}

// NewActiveSessions initializes new instance of ActiveSessions object.
func NewActiveSessions() *ActiveSessions {
	return &ActiveSessions{
		data: make(map[int64]Session),
	}
}

// Set adds new or replaces existing session.
func (as *ActiveSessions) Set(userID int64, s Session) {
	as.mu.Lock()
	defer as.mu.Unlock()

	as.data[userID] = s
}

// Get returns session of user with specified ID. If the session
// doesn't exist, second parameter will be equal to false.
func (as *ActiveSessions) Get(userID int64) (Session, bool) {
	as.mu.Lock()
	defer as.mu.Unlock()

	s, ok := as.data[userID]
	if !ok {
		return Session{}, false
	}

	return s, true
}

// Delete removes session of user with specified ID.
// If the session doesn't exist, nothing will happens.
func (as *ActiveSessions) Delete(userID int64) {
	as.mu.Lock()
	defer as.mu.Unlock()

	delete(as.data, userID)
}

package sessions

import (
	"reflect"
	"testing"
	"time"

	"golang.org/x/text/language"
	"golang.org/x/text/message"

	"github.com/lazy-void/primitive-bot/pkg/menu"
	"github.com/lazy-void/primitive-bot/pkg/primitive"
)

func TestNewSession(t *testing.T) {
	var userID, menuMessageID int64 = 123456789, 987654321
	imgPath := "path/to/image.png"
	expectedConfig := primitive.NewConfig()
	expectedMenu := menu.New(expectedConfig)
	menu.InitText(*message.NewPrinter(language.English))

	s := NewSession(userID, menuMessageID, imgPath)

	switch {
	case s.UserID != userID:
		t.Errorf("session.UserID = %d; want %d", s.UserID, userID)
	case s.MenuMessageID != menuMessageID:
		t.Errorf("session.menuMessageID = %d; want %d", s.MenuMessageID, menuMessageID)
	case s.State != InMenu:
		t.Errorf("session.State = %v; want %v", s.State, InMenu)
	case s.Input == nil:
		t.Error("session.Input must be non-nil")
	case s.QuitInput == nil:
		t.Error("session.QuitInput must be non-nil")
	case s.ImgPath != imgPath:
		t.Errorf("session.ImgPath = %v; want %v", s.ImgPath, imgPath)
	case reflect.DeepEqual(s.Menu, expectedMenu):
		t.Errorf("session.Menu = %+v; want %+v", s.Menu, expectedMenu)
	case s.Config != expectedConfig:
		t.Errorf("session.Config = %+v; want %+v", s.Config, expectedConfig)
	}
}

func TestNewActiveSessionsStartsTimeouterThatTerminatesInactiveSessions(t *testing.T) {
	timeout := 10 * time.Millisecond
	frequency := 5 * time.Millisecond
	var userID int64 = 123456789
	session := NewSession(userID, 123, "img.png")

	// create
	as := NewActiveSessions(timeout, frequency)

	// add session
	as.Set(userID, session)

	// wait
	time.Sleep(timeout + frequency)

	// check that session is terminated
	if _, ok := as.Get(userID); ok {
		t.Error("inactive session wasn't terminated.")
	}
}

func TestNewActiveSessionsStartsTimeouterThatDoesNotTerminateActiveSessions(t *testing.T) {
	timeout := 50 * time.Millisecond
	frequency := 10 * time.Millisecond
	var userID int64 = 123456789
	session := NewSession(userID, 123, "img.png")

	// create
	as := NewActiveSessions(timeout, frequency)

	// add session
	as.Set(userID, session)

	// wait
	time.Sleep(timeout / 2)

	// update time of last request
	as.Get(userID)

	// wait
	time.Sleep(frequency + timeout/2)

	// check that session wasn't terminated
	if _, ok := as.Get(userID); !ok {
		t.Error("active session was terminated.")
	}
}

func TestNewActiveSessionsWhenTerminatedSessionIsInInputMenuState(t *testing.T) {
	timeout := time.Millisecond
	frequency := time.Millisecond
	var userID int64 = 123456789
	session := NewSession(userID, 123, "img.png")
	session.State = InInputDialog

	// create
	as := NewActiveSessions(timeout, frequency)

	// add session
	as.Set(userID, session)

	// wait for signal from quit channel
	after := time.After(100 * time.Millisecond)
	select {
	case <-after:
		t.Error("signal on quit channel was not sent.")
	case <-session.QuitInput:
	}
}

func TestActiveSessions_Set(t *testing.T) {
	timeout := 100 * time.Second
	frequency := 100 * time.Second
	var userID int64 = 123456789
	session := NewSession(userID, 123, "img.png")

	as := NewActiveSessions(timeout, frequency)

	// add session
	as.Set(userID, session)

	s, ok := as.sessions[userID]
	if !ok {
		t.Error("session isn't in the active sessions.")
	} else if !reflect.DeepEqual(s, session) {
		t.Errorf("session = %+v; want %+v ", s, session)
	}

	// update session
	session.Config.Shape = primitive.ShapeEllipse
	as.Set(userID, session)

	s, ok = as.sessions[userID]
	if !ok {
		t.Error("session isn't in the active sessions.")
	} else if !reflect.DeepEqual(s, session) {
		t.Errorf("session = %+v; want %+v ", s, session)
	}
}

func TestActiveSessions_Get(t *testing.T) {
	timeout := 100 * time.Second
	frequency := 100 * time.Second
	var userID int64 = 123456789
	session := NewSession(userID, 123, "img.png")

	as := NewActiveSessions(timeout, frequency)

	// when session is not in the active session
	_, ok := as.Get(userID)
	if ok {
		t.Error("session mustn't be in the active sessions.")
	}

	// add session
	as.Set(userID, session)

	s, ok := as.Get(userID)
	if !ok {
		t.Error("session must be in the active sessions.")
	} else if reflect.DeepEqual(s, session) {
		t.Errorf("session = %+v; want %+v ", s, session)
	}
}

package sessions

import (
	"log"
	"reflect"
	"sync"
	"testing"
	"time"

	"golang.org/x/text/language"
	"golang.org/x/text/message"

	"github.com/lazy-void/primitive-bot/pkg/menu"
	"github.com/lazy-void/primitive-bot/pkg/primitive"
)

type out struct {
	sync.Mutex
	data string
}

func (o *out) Write(p []byte) (n int, err error) {
	o.Lock()
	defer o.Unlock()

	o.data = string(p)
	return len(p), nil
}

func (o *out) Read() string {
	o.Lock()
	defer o.Unlock()

	return o.data
}

func TestNewSession(t *testing.T) {
	var userID, menuMessageID int64 = 123456789, 987654321
	imgPath := "path/to/image.png"
	expectedConfig := primitive.New(1)
	expectedMenu := menu.New(expectedConfig)
	menu.InitText(message.NewPrinter(language.English))

	s := NewSession(userID, menuMessageID, imgPath, 1)

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
	session := NewSession(userID, 123, "img.png", 1)

	// create
	as := NewActiveSessions(timeout, frequency, nil)

	// add session
	as.Set(userID, session, false)

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
	session := NewSession(userID, 123, "img.png", 1)

	// create
	as := NewActiveSessions(timeout, frequency, nil)

	// add session
	as.Set(userID, session, false)

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
	session := NewSession(userID, 123, "img.png", 1)
	session.State = InInputDialog

	// create
	as := NewActiveSessions(timeout, frequency, nil)

	// add session
	as.Set(userID, session, false)

	// wait for signal from quit channel
	after := time.After(100 * timeout)
	select {
	case <-after:
		t.Error("signal on quit channel was not sent.")
	case <-session.QuitInput:
	}
}

func TestTimeouterWhenTerminatedSessionIsInInputMenuStateButNobodyListensQuitChannelLogsAnError(t *testing.T) {
	timeout := time.Millisecond
	frequency := time.Millisecond
	var userID int64 = 123456789
	session := NewSession(userID, 123, "img.png", 1)
	session.State = InInputDialog
	expected := "nobody listens on the QuitInput channel\n"

	// create and add session
	as := &ActiveSessions{
		sessions: make(map[int64]Session),
		timeout:  timeout,
	}
	as.Set(userID, session, false)

	logOut := &out{}
	go as.timeouter(frequency, log.New(logOut, "", 0))
	time.Sleep(frequency * 10)

	res := logOut.Read()
	if res != expected {
		t.Errorf("got log message: %q; want %q", res, expected)
	}
}

func TestActiveSessions_Set(t *testing.T) {
	timeout := 100 * time.Second
	frequency := 100 * time.Second
	var userID int64 = 123456789
	session := NewSession(userID, 123, "img.png", 1)

	as := NewActiveSessions(timeout, frequency, nil)

	// add session
	as.Set(userID, session, false)

	s, ok := as.sessions[userID]
	if !ok {
		t.Error("session isn't in the active sessions.")
	} else if !reflect.DeepEqual(s, session) {
		t.Errorf("session = %+v; want %+v ", s, session)
	}

	// update session
	session.Config.Shape = primitive.ShapeEllipse
	as.Set(userID, session, false)

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
	session := NewSession(userID, 123, "img.png", 1)

	as := NewActiveSessions(timeout, frequency, nil)

	// when session is not in the active session
	_, ok := as.Get(userID)
	if ok {
		t.Error("session mustn't be in the active sessions.")
	}

	// add session
	as.Set(userID, session, false)

	s, ok := as.Get(userID)
	if !ok {
		t.Error("session must be in the active sessions.")
	} else if reflect.DeepEqual(s, session) {
		t.Errorf("session = %+v; want %+v ", s, session)
	}
}

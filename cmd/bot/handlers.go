package main

import (
	"fmt"

	"github.com/lazy-void/primitive-bot/pkg/primitive"
	"github.com/lazy-void/primitive-bot/pkg/queue"

	"github.com/lazy-void/primitive-bot/pkg/menu"

	"github.com/lazy-void/primitive-bot/pkg/sessions"
	"github.com/lazy-void/primitive-bot/pkg/tg"
)

func (app *application) showRootMenuView(s sessions.Session) {
	app.showMenuView(s.UserID, s.MenuMessageID, s.Menu.RootView)
}

func (app *application) handleCreateButton(s sessions.Session) {
	n := app.queue.GetNumOperations(s.UserID)
	if n >= app.operationsLimit {
		if _, err := app.bot.SendMessage(s.UserID, operationsLimitMessage); err != nil {
			app.serverError(s.UserID, err)
		}
		return
	}

	app.infoLog.Printf(enqueuedLogMessage, s.UserID, s.ImgPath, s.Config.Iterations, s.Config.Shape,
		s.Config.Alpha, s.Config.Repeat, s.Config.OutputSize, s.Config.Extension)
	pos := app.queue.Enqueue(queue.Operation{
		UserID:  s.UserID,
		ImgPath: s.ImgPath,
		Config:  s.Config,
	})

	_, err := app.bot.SendMessage(s.UserID, app.createStatusMessage(s.Config, pos))
	if err != nil {
		app.serverError(s.UserID, err)
	}
}

func (app *application) showShapesMenuView(s sessions.Session) {
	app.showMenuView(s.UserID, s.MenuMessageID, s.Menu.ShapesView)
}

func (app *application) handleShapesButton(s sessions.Session, n int) {
	s.Config.Shape = primitive.Shape(n)
	app.sessions.Set(s.UserID, s)

	// update menu
	selected := fmt.Sprintf("%s/%d", menu.ShapesViewCallback, s.Config.Shape)
	s.Menu.ShapesView = menu.NewMenuView(menu.ShapesViewTmpl, selected)
	app.sessions.Set(s.UserID, s)

	app.showMenuView(s.UserID, s.MenuMessageID, s.Menu.ShapesView)
}

func (app *application) showIterMenuView(s sessions.Session) {
	app.showMenuView(s.UserID, s.MenuMessageID, s.Menu.IterView)
}

func (app *application) handleIterButton(s sessions.Session, n int) {
	if n > 5000 {
		return
	}
	s.Config.Iterations = n
	app.sessions.Set(s.UserID, s)

	// update menu
	selected := fmt.Sprintf("%s/%d", menu.IterViewCallback, s.Config.Iterations)
	s.Menu.IterView = menu.NewMenuView(menu.IterViewTmpl, selected)
	app.sessions.Set(s.UserID, s)

	app.showMenuView(s.UserID, s.MenuMessageID, s.Menu.IterView)
}

func (app *application) handleIterInput(s sessions.Session) {
	in := make(chan tg.Message)
	out := make(chan int)

	// Get user input
	s.InChan = in
	app.sessions.Set(s.UserID, s)

	go app.getInputFromUser(s.UserID, s.MenuMessageID, 1, 5000, in, out)

	s.Config.Iterations = <-out
	s.InChan = nil
	app.sessions.Set(s.UserID, s)
	close(in)

	buttonText := fmt.Sprintf("%s (%d)", menu.OtherButtonText, s.Config.Iterations)
	s.Menu.IterView = menu.NewMenuView(
		menu.IterViewTmpl, menu.IterInputCallback, buttonText)
	app.sessions.Set(s.UserID, s)

	app.showMenuView(s.UserID, s.MenuMessageID, s.Menu.IterView)
}

func (app *application) showRepMenuView(s sessions.Session) {
	app.showMenuView(s.UserID, s.MenuMessageID, s.Menu.RepView)
}

func (app *application) handleRepButton(s sessions.Session, n int) {
	s.Config.Repeat = n
	app.sessions.Set(s.UserID, s)

	// update menu
	selected := fmt.Sprintf("%s/%d", menu.RepViewCallback, s.Config.Repeat)
	s.Menu.RepView = menu.NewMenuView(menu.RepViewTmpl, selected)
	app.sessions.Set(s.UserID, s)

	app.showMenuView(s.UserID, s.MenuMessageID, s.Menu.RepView)
}

func (app *application) showAlphaMenuView(s sessions.Session) {
	app.showMenuView(s.UserID, s.MenuMessageID, s.Menu.AlphaView)
}

func (app *application) handleAlphaButton(s sessions.Session, n int) {
	app.showMenuView(s.UserID, s.MenuMessageID, s.Menu.AlphaView)

	if n < 0 || n > 255 {
		return
	}
	s.Config.Alpha = n
	app.sessions.Set(s.UserID, s)

	// update menu
	selected := fmt.Sprintf("%s/%d", menu.AlphaViewCallback, s.Config.Alpha)
	s.Menu.AlphaView = menu.NewMenuView(menu.AlphaViewTmpl, selected)
	app.sessions.Set(s.UserID, s)

	app.showMenuView(s.UserID, s.MenuMessageID, s.Menu.AlphaView)
}

func (app *application) handleAlphaInput(s sessions.Session) {
	in := make(chan tg.Message)
	out := make(chan int)

	// Get user input
	s.InChan = in
	app.sessions.Set(s.UserID, s)

	go app.getInputFromUser(s.UserID, s.MenuMessageID, 1, 255, in, out)

	s.Config.Alpha = <-out
	s.InChan = nil
	app.sessions.Set(s.UserID, s)
	close(in)

	buttonText := fmt.Sprintf("%s (%d)", menu.OtherButtonText, s.Config.Alpha)
	s.Menu.AlphaView = menu.NewMenuView(
		menu.AlphaViewTmpl, menu.AlphaInputCallback, buttonText)
	app.sessions.Set(s.UserID, s)

	app.showMenuView(s.UserID, s.MenuMessageID, s.Menu.AlphaView)
}

func (app *application) showExtMenuView(s sessions.Session) {
	app.showMenuView(s.UserID, s.MenuMessageID, s.Menu.ExtView)
}

func (app *application) handleExtButton(s sessions.Session, ext string) {
	s.Config.Extension = ext
	app.sessions.Set(s.UserID, s)

	// update menu
	selected := fmt.Sprintf("%s/%s", menu.ExtViewCallback, s.Config.Extension)
	s.Menu.ExtView = menu.NewMenuView(menu.ExtViewTmpl, selected)
	app.sessions.Set(s.UserID, s)

	app.showMenuView(s.UserID, s.MenuMessageID, s.Menu.ExtView)
}

func (app *application) showSizeMenuView(s sessions.Session) {
	app.showMenuView(s.UserID, s.MenuMessageID, s.Menu.SizeView)
}

func (app *application) handleSizeButton(s sessions.Session, n int) {
	if n < 256 || n > 1920 {
		return
	}
	s.Config.OutputSize = n
	app.sessions.Set(s.UserID, s)

	// update menu
	selected := fmt.Sprintf("%s/%d", menu.SizeViewCallback, s.Config.OutputSize)
	s.Menu.SizeView = menu.NewMenuView(menu.SizeViewTmpl, selected)
	app.sessions.Set(s.UserID, s)

	app.showMenuView(s.UserID, s.MenuMessageID, s.Menu.SizeView)
}

func (app *application) handleSizeInput(s sessions.Session) {
	in := make(chan tg.Message)
	out := make(chan int)

	// Get user input
	s.InChan = in
	app.sessions.Set(s.UserID, s)

	go app.getInputFromUser(s.UserID, s.MenuMessageID, 256, 3840, in, out)

	s.Config.OutputSize = <-out
	s.InChan = nil
	app.sessions.Set(s.UserID, s)
	close(in)

	buttonText := fmt.Sprintf("%s (%d)", menu.OtherButtonText, s.Config.OutputSize)
	s.Menu.SizeView = menu.NewMenuView(
		menu.SizeViewTmpl, menu.SizeInputCallback, buttonText)
	app.sessions.Set(s.UserID, s)

	app.showMenuView(s.UserID, s.MenuMessageID, s.Menu.SizeView)
}

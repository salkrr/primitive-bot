package main

import (
	"fmt"

	"github.com/lazy-void/primitive-bot/pkg/primitive"
	"github.com/lazy-void/primitive-bot/pkg/queue"

	"github.com/lazy-void/primitive-bot/pkg/menu"

	"github.com/lazy-void/primitive-bot/pkg/sessions"
	"github.com/lazy-void/primitive-bot/pkg/tg"
)

func (app *application) showRootMenuActivity(s sessions.Session) {
	app.showMenuActivity(s.UserID, s.MenuMessageID, s.Menu.RootActivity)
}

func (app *application) handleCreateButton(s sessions.Session) {
	n := app.queue.GetNumOperations(s.UserID)
	if n >= app.operationsLimit {
		if _, err := app.bot.SendMessage(s.UserID, operationsLimitMessage); err != nil {
			app.serverError(s.UserID, err)
		}
		return
	}

	pos := app.queue.Enqueue(queue.Operation{
		ChatID:  s.UserID,
		ImgPath: s.ImgPath,
		Config:  s.Config,
	})
	_, err := app.bot.SendMessage(s.UserID, app.createStatusMessage(s.Config, pos))
	if err != nil {
		app.serverError(s.UserID, err)
	}
}

func (app *application) showShapesMenuActivity(s sessions.Session) {
	app.showMenuActivity(s.UserID, s.MenuMessageID, s.Menu.ShapesActivity)
}

func (app *application) handleShapesButton(s sessions.Session, n int) {
	s.Config.Shape = primitive.Shape(n)
	app.sessions.Set(s.UserID, s)

	// update menu
	selected := fmt.Sprintf("%s/%d", menu.ShapesActivityCallback, s.Config.Shape)
	s.Menu.ShapesActivity = menu.NewMenuActivity(menu.ShapesActivityTmpl, selected)
	app.sessions.Set(s.UserID, s)

	app.showMenuActivity(s.UserID, s.MenuMessageID, s.Menu.ShapesActivity)
}

func (app *application) showIterMenuActivity(s sessions.Session) {
	app.showMenuActivity(s.UserID, s.MenuMessageID, s.Menu.IterActivity)
}

func (app *application) handleIterButton(s sessions.Session, n int) {
	if n > 5000 {
		return
	}
	s.Config.Iterations = n
	app.sessions.Set(s.UserID, s)

	// update menu
	selected := fmt.Sprintf("%s/%d", menu.IterActivityCallback, s.Config.Iterations)
	s.Menu.IterActivity = menu.NewMenuActivity(menu.IterActivityTmpl, selected)
	app.sessions.Set(s.UserID, s)

	app.showMenuActivity(s.UserID, s.MenuMessageID, s.Menu.IterActivity)
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
	s.Menu.IterActivity = menu.NewMenuActivity(
		menu.IterActivityTmpl, menu.IterInputCallback, buttonText)
	app.sessions.Set(s.UserID, s)

	app.showMenuActivity(s.UserID, s.MenuMessageID, s.Menu.IterActivity)
}

func (app *application) showRepMenuActivity(s sessions.Session) {
	app.showMenuActivity(s.UserID, s.MenuMessageID, s.Menu.RepActivity)
}

func (app *application) handleRepButton(s sessions.Session, n int) {
	s.Config.Repeat = n
	app.sessions.Set(s.UserID, s)

	// update menu
	selected := fmt.Sprintf("%s/%d", menu.RepActivityCallback, s.Config.Repeat)
	s.Menu.RepActivity = menu.NewMenuActivity(menu.RepActivityTmpl, selected)
	app.sessions.Set(s.UserID, s)

	app.showMenuActivity(s.UserID, s.MenuMessageID, s.Menu.RepActivity)
}

func (app *application) showAlphaMenuActivity(s sessions.Session) {
	app.showMenuActivity(s.UserID, s.MenuMessageID, s.Menu.AlphaActivity)
}

func (app *application) handleAlphaButton(s sessions.Session, n int) {
	app.showMenuActivity(s.UserID, s.MenuMessageID, s.Menu.AlphaActivity)

	if n < 0 || n > 255 {
		return
	}
	s.Config.Alpha = n
	app.sessions.Set(s.UserID, s)

	// update menu
	selected := fmt.Sprintf("%s/%d", menu.AlphaActivityCallback, s.Config.Alpha)
	s.Menu.AlphaActivity = menu.NewMenuActivity(menu.AlphaActivityTmpl, selected)
	app.sessions.Set(s.UserID, s)

	app.showMenuActivity(s.UserID, s.MenuMessageID, s.Menu.AlphaActivity)
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
	s.Menu.AlphaActivity = menu.NewMenuActivity(
		menu.AlphaActivityTmpl, menu.AlphaInputCallback, buttonText)
	app.sessions.Set(s.UserID, s)

	app.showMenuActivity(s.UserID, s.MenuMessageID, s.Menu.AlphaActivity)
}

func (app *application) showExtMenuActivity(s sessions.Session) {
	app.showMenuActivity(s.UserID, s.MenuMessageID, s.Menu.ExtActivity)
}

func (app *application) handleExtButton(s sessions.Session, ext string) {
	s.Config.Extension = primitive.Extension(ext)
	app.sessions.Set(s.UserID, s)

	// update menu
	selected := fmt.Sprintf("%s/%s", menu.ExtActivityCallback, s.Config.Extension)
	s.Menu.ExtActivity = menu.NewMenuActivity(menu.ExtActivityTmpl, selected)
	app.sessions.Set(s.UserID, s)

	app.showMenuActivity(s.UserID, s.MenuMessageID, s.Menu.ExtActivity)
}

func (app *application) showSizeMenuActivity(s sessions.Session) {
	app.showMenuActivity(s.UserID, s.MenuMessageID, s.Menu.SizeActivity)
}

func (app *application) handleSizeButton(s sessions.Session, n int) {
	if n < 256 || n > 1920 {
		return
	}
	s.Config.OutputSize = n
	app.sessions.Set(s.UserID, s)

	// update menu
	selected := fmt.Sprintf("%s/%d", menu.SizeActivityCallback, s.Config.OutputSize)
	s.Menu.SizeActivity = menu.NewMenuActivity(menu.SizeActivityTmpl, selected)
	app.sessions.Set(s.UserID, s)

	app.showMenuActivity(s.UserID, s.MenuMessageID, s.Menu.SizeActivity)
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
	s.Menu.SizeActivity = menu.NewMenuActivity(
		menu.SizeActivityTmpl, menu.SizeInputCallback, buttonText)
	app.sessions.Set(s.UserID, s)

	app.showMenuActivity(s.UserID, s.MenuMessageID, s.Menu.SizeActivity)
}

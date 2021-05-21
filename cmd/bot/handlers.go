package main

import (
	"fmt"
	"os"

	"github.com/lazy-void/primitive-bot/pkg/menu"

	"github.com/lazy-void/primitive-bot/pkg/primitive"
	"github.com/lazy-void/primitive-bot/pkg/queue"
	"github.com/lazy-void/primitive-bot/pkg/sessions"
	"github.com/lazy-void/primitive-bot/pkg/telegram"
)

func (app *application) handleMessage(m telegram.Message) {
	if m.Photo != nil {
		app.handlePhotoMessage(m)
		return
	}

	// Handle user input if they are inside input form
	s, ok := app.sessions.Get(m.From.ID)
	if ok && s.InChan != nil {
		s.InChan <- m
		return
	}

	if m.Text == "/status" {
		operations, positions := app.queue.GetOperations(s.ChatID)
		if len(operations) == 0 {
			_, err := app.bot.SendMessage(m.Chat.ID, StatusEmptyMessage)
			if err != nil {
				app.serverError(m.Chat.ID, err)
			}
			return
		}

		for i, op := range operations {
			_, err := app.bot.SendMessage(m.Chat.ID, app.createStatusMessage(op.Config, positions[i]))
			if err != nil {
				app.serverError(m.Chat.ID, err)
				return
			}
		}

		return
	}

	// Send help message
	_, err := app.bot.SendMessage(m.Chat.ID, HelpMessage)
	if err != nil {
		app.serverError(m.Chat.ID, err)
	}
}

func (app *application) handlePhotoMessage(m telegram.Message) {
	// If we already have session - delete it's menu
	s, ok := app.sessions.Get(m.From.ID)
	if ok {
		err := app.bot.DeleteMessage(s.ChatID, s.MenuMessageID)
		if err != nil {
			app.serverError(m.Chat.ID, err)
			return
		}
	}

	// Choose smallest image with dimensions >= 256
	var file telegram.PhotoSize
	for _, photo := range m.Photo {
		file = photo
		if photo.Width >= 256 && photo.Height >= 256 {
			break
		}
	}
	if file.FileID == "" {
		app.serverError(m.Chat.ID, fmt.Errorf("no image files in %v", m.Photo))
		return
	}

	path := fmt.Sprintf("%s/%s.jpg", app.inDir, file.FileUniqueID)
	if _, err := os.Stat(path); os.IsNotExist(err) {
		img, err := app.bot.DownloadFile(file.FileID)
		if err != nil {
			app.serverError(m.Chat.ID, fmt.Errorf("couldn't download image: %s", err))
			return
		}

		if err := os.WriteFile(path, img, 0644); err != nil {
			app.serverError(m.Chat.ID, fmt.Errorf("couldn't save image: %s", err))
			return
		}
	}

	// Create session
	msg, err := app.bot.SendMessage(m.Chat.ID, menu.RootActivityTmpl.Text, menu.RootActivityTmpl.Keyboard)
	if err != nil {
		app.serverError(m.Chat.ID, err)
		return
	}
	app.sessions.Set(m.From.ID, sessions.NewSession(m.Chat.ID, msg.MessageID, path))
}

func (app *application) handleCallbackQuery(q telegram.CallbackQuery) {
	defer func() {
		if err := app.bot.AnswerCallbackQuery(q.ID, ""); err != nil {
			app.errorLog.Printf("Answer callback query error: %s", err)
		}
	}()

	s, ok := app.sessions.Get(q.From.ID)
	if !ok || q.Message.MessageID != s.MenuMessageID {
		err := app.bot.DeleteMessage(q.Message.Chat.ID, q.Message.MessageID)
		if err != nil {
			app.errorLog.Printf("Deleting message error: %s", err)
		}
		return
	}

	var num int
	var slug string
	switch {
	case match(q.Data, menu.RootMenuCallback):
		app.showActivity(s.ChatID, s.MenuMessageID, s.Menu.RootActivity)
	case match(q.Data, menu.CreateButtonCallback):
		n := app.queue.GetNumOperations(s.ChatID)
		if n >= app.operationsLimit {
			if _, err := app.bot.SendMessage(s.ChatID, OperationsLimitMessage); err != nil {
				app.serverError(s.ChatID, err)
			}
			return
		}

		pos := app.queue.Enqueue(queue.Operation{
			ChatID:  s.ChatID,
			ImgPath: s.ImgPath,
			Config:  s.Config,
		})
		_, err := app.bot.SendMessage(s.ChatID, app.createStatusMessage(s.Config, pos))
		if err != nil {
			app.serverError(s.ChatID, err)
		}
	case match(q.Data, menu.ShapesButtonCallback):
		app.showActivity(s.ChatID, s.MenuMessageID, s.Menu.ShapesActivity)
	case match(q.Data, fmt.Sprintf("%s/([0-8])", menu.ShapesButtonCallback), &num):
		s.Config.Shape = primitive.Shape(num)
		app.sessions.Set(q.From.ID, s)

		// update menu
		selected := fmt.Sprintf("%s/%d", menu.ShapesButtonCallback, s.Config.Shape)
		s.Menu.ShapesActivity = menu.NewMenuActivity(menu.ShapesActivityTmpl, selected)
		app.sessions.Set(q.From.ID, s)

		app.showActivity(s.ChatID, s.MenuMessageID, s.Menu.ShapesActivity)
	case match(q.Data, menu.IterButtonCallback):
		app.showActivity(s.ChatID, s.MenuMessageID, s.Menu.IterActivity)
	case match(q.Data, fmt.Sprintf("%s/([0-9]+)", menu.IterButtonCallback), &num):
		if num > 5000 {
			return
		}
		s.Config.Iterations = num
		app.sessions.Set(q.From.ID, s)

		// update menu
		// update menu
		selected := fmt.Sprintf("%s/%d", menu.IterButtonCallback, s.Config.Iterations)
		s.Menu.IterActivity = menu.NewMenuActivity(menu.IterActivityTmpl, selected)
		app.sessions.Set(q.From.ID, s)

		app.showActivity(s.ChatID, s.MenuMessageID, s.Menu.IterActivity)
	case match(q.Data, menu.IterInputCallback):
		app.handleIterInput(q, s)
	case match(q.Data, menu.RepButtonCallback):
		app.showActivity(s.ChatID, s.MenuMessageID, s.Menu.RepActivity)
	case match(q.Data, fmt.Sprintf("%s/([1-6])", menu.RepButtonCallback), &num):
		s.Config.Repeat = num
		app.sessions.Set(q.From.ID, s)

		// update menu
		selected := fmt.Sprintf("%s/%d", menu.RepButtonCallback, s.Config.Repeat)
		s.Menu.RepActivity = menu.NewMenuActivity(menu.RepActivityTmpl, selected)
		app.sessions.Set(q.From.ID, s)

		app.showActivity(s.ChatID, s.MenuMessageID, s.Menu.RepActivity)
	case match(q.Data, menu.AlphaButtonCallback):
		app.showActivity(s.ChatID, s.MenuMessageID, s.Menu.AlphaActivity)
	case match(q.Data, fmt.Sprintf("%s/([0-9]+)", menu.AlphaButtonCallback), &num):
		if num < 0 || num > 255 {
			return
		}
		s.Config.Alpha = num
		app.sessions.Set(q.From.ID, s)

		// update menu
		selected := fmt.Sprintf("%s/%d", menu.AlphaButtonCallback, s.Config.Alpha)
		s.Menu.AlphaActivity = menu.NewMenuActivity(menu.AlphaActivityTmpl, selected)
		app.sessions.Set(q.From.ID, s)

		app.showActivity(s.ChatID, s.MenuMessageID, s.Menu.AlphaActivity)
	case match(q.Data, menu.AlphaInputCallback):
		app.handleAlphaInput(q, s)
	case match(q.Data, menu.ExtButtonCallback):
		app.showActivity(s.ChatID, s.MenuMessageID, s.Menu.ExtActivity)
	case match(q.Data, fmt.Sprintf("%s/(jpg|png|svg|gif)", menu.ExtButtonCallback), &slug):
		s.Config.Extension = primitive.Extension(slug)
		app.sessions.Set(q.From.ID, s)

		// update menu
		selected := fmt.Sprintf("%s/%s", menu.ExtButtonCallback, s.Config.Extension)
		s.Menu.ExtActivity = menu.NewMenuActivity(menu.ExtActivityTmpl, selected)
		app.sessions.Set(q.From.ID, s)

		app.showActivity(s.ChatID, s.MenuMessageID, s.Menu.ExtActivity)
	case match(q.Data, menu.SizeButtonCallback):
		app.showActivity(s.ChatID, s.MenuMessageID, s.Menu.SizeActivity)
	case match(q.Data, fmt.Sprintf("%s/([0-9]+)", menu.SizeButtonCallback), &num):
		if num < 256 || num > 1920 {
			return
		}
		s.Config.OutputSize = num
		app.sessions.Set(q.From.ID, s)

		// update menu
		selected := fmt.Sprintf("%s/%d", menu.SizeButtonCallback, s.Config.OutputSize)
		s.Menu.SizeActivity = menu.NewMenuActivity(menu.SizeActivityTmpl, selected)
		app.sessions.Set(q.From.ID, s)

		app.showActivity(s.ChatID, s.MenuMessageID, s.Menu.SizeActivity)
	case match(q.Data, menu.SizeInputCallback):
		app.handleSizeInput(q, s)
	}
}

func (app *application) handleIterInput(
	q telegram.CallbackQuery,
	s sessions.Session,
) {
	in := make(chan telegram.Message)
	out := make(chan int)

	// Get user input
	s.InChan = in
	app.sessions.Set(q.From.ID, s)

	go app.getInputFromUser(q.Message.Chat.ID, q.Message.MessageID, 1, 5000, in, out)

	s.Config.Iterations = <-out
	s.InChan = nil
	app.sessions.Set(q.From.ID, s)
	close(in)

	buttonText := fmt.Sprintf("%s (%d)", menu.OtherButtonText, s.Config.Iterations)
	s.Menu.IterActivity = menu.NewMenuActivity(
		menu.IterActivityTmpl, menu.IterInputCallback, buttonText)
	app.sessions.Set(q.From.ID, s)

	app.showActivity(q.Message.Chat.ID, q.Message.MessageID, s.Menu.IterActivity)
}

func (app *application) handleAlphaInput(
	q telegram.CallbackQuery,
	s sessions.Session,
) {
	in := make(chan telegram.Message)
	out := make(chan int)

	// Get user input
	s.InChan = in
	app.sessions.Set(q.From.ID, s)

	go app.getInputFromUser(q.Message.Chat.ID, q.Message.MessageID, 1, 255, in, out)

	s.Config.Alpha = <-out
	s.InChan = nil
	app.sessions.Set(q.From.ID, s)
	close(in)

	buttonText := fmt.Sprintf("%s (%d)", menu.OtherButtonText, s.Config.Alpha)
	s.Menu.AlphaActivity = menu.NewMenuActivity(
		menu.AlphaActivityTmpl, menu.AlphaInputCallback, buttonText)
	app.sessions.Set(q.From.ID, s)

	app.showActivity(q.Message.Chat.ID, q.Message.MessageID, s.Menu.AlphaActivity)
}

func (app *application) handleSizeInput(
	q telegram.CallbackQuery,
	s sessions.Session,
) {
	in := make(chan telegram.Message)
	out := make(chan int)

	// Get user input
	s.InChan = in
	app.sessions.Set(q.From.ID, s)

	go app.getInputFromUser(q.Message.Chat.ID, q.Message.MessageID, 256, 3840, in, out)

	s.Config.OutputSize = <-out
	s.InChan = nil
	app.sessions.Set(q.From.ID, s)
	close(in)

	buttonText := fmt.Sprintf("%s (%d)", menu.OtherButtonText, s.Config.OutputSize)
	s.Menu.SizeActivity = menu.NewMenuActivity(
		menu.SizeActivityTmpl, menu.SizeInputCallback, buttonText)
	app.sessions.Set(q.From.ID, s)

	app.showActivity(q.Message.Chat.ID, q.Message.MessageID, s.Menu.SizeActivity)
}

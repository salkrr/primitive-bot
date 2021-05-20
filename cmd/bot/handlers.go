package main

import (
	"fmt"
	"os"

	"github.com/lazy-void/primitive-bot/pkg/primitive"
	"github.com/lazy-void/primitive-bot/pkg/queue"
	"github.com/lazy-void/primitive-bot/pkg/sessions"
	"github.com/lazy-void/primitive-bot/pkg/telegram"
)

const (
	rootMenuCallback     = "/"
	createButtonCallback = "/create"
	shapesButtonCallback = "/shape"
	iterButtonCallback   = "/iter"
	iterInputCallback    = "/iter/input"
	repButtonCallback    = "/rep"
	alphaButtonCallback  = "/alpha"
	alphaInputCallback   = "/alpha/input"
	extButtonCallback    = "/ext"
	sizeButtonCallback   = "/size"
	sizeInputCallback    = "/size/input"
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
			_, err := app.bot.SendMessage(m.Chat.ID, statusEmptyMessage)
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
	_, err := app.bot.SendMessage(m.Chat.ID, helpMessage)
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
	msg, err := app.bot.SendMessage(m.Chat.ID, rootMenuText, rootKeyboard)
	if err != nil {
		app.serverError(m.Chat.ID, err)
		return
	}
	app.sessions.Set(m.From.ID, sessions.Session{
		ChatID:        m.Chat.ID,
		MenuMessageID: msg.MessageID,
		InChan:        nil,
		ImgPath:       path,
		Config:        primitive.NewConfig(),
	})
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
	case match(q.Data, rootMenuCallback):
		err := app.bot.EditMessageText(q.Message.Chat.ID, q.Message.MessageID, rootMenuText, rootKeyboard)
		if err != nil {
			app.serverError(s.ChatID, err)
		}
	case match(q.Data, createButtonCallback):
		n := app.queue.GetNumOperations(s.ChatID)
		if n >= app.operationsLimit {
			if _, err := app.bot.SendMessage(s.ChatID, operationsLimitMessage); err != nil {
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
	case match(q.Data, shapesButtonCallback):
		selected := fmt.Sprintf("%s/%d", shapesButtonCallback, s.Config.Shape)
		app.showMenu(q.Message.Chat.ID, q.Message.MessageID, selected, shapesMenuText, shapesKeyboard)
	case match(q.Data, fmt.Sprintf("%s/([0-8])", shapesButtonCallback), &num):
		s.Config.Shape = primitive.Shape(num)
		app.sessions.Set(q.From.ID, s)

		// update menu
		selected := fmt.Sprintf("%s/%d", shapesButtonCallback, s.Config.Shape)
		app.showMenu(q.Message.Chat.ID, q.Message.MessageID, selected, shapesMenuText, shapesKeyboard)
	case match(q.Data, iterButtonCallback):
		selected := fmt.Sprintf("%s/%d", iterButtonCallback, s.Config.Iterations)
		app.showMenu(q.Message.Chat.ID, q.Message.MessageID, selected, iterMenuText, iterKeyboard)
	case match(q.Data, fmt.Sprintf("%s/([0-9]+)", iterButtonCallback), &num):
		if num > 5000 {
			return
		}
		s.Config.Iterations = num
		app.sessions.Set(q.From.ID, s)

		// update menu
		selected := fmt.Sprintf("%s/%d", iterButtonCallback, s.Config.Iterations)
		app.showMenu(q.Message.Chat.ID, q.Message.MessageID, selected, iterMenuText, iterKeyboard)
	case match(q.Data, iterInputCallback):
		app.handleIterInput(q, s)
	case match(q.Data, repButtonCallback):
		selected := fmt.Sprintf("%s/%d", repButtonCallback, s.Config.Repeat)
		app.showMenu(q.Message.Chat.ID, q.Message.MessageID, selected, repMenuText, repKeyboard)
	case match(q.Data, fmt.Sprintf("%s/([1-6])", repButtonCallback), &num):
		s.Config.Repeat = num
		app.sessions.Set(q.From.ID, s)

		// update menu
		selected := fmt.Sprintf("%s/%d", repButtonCallback, s.Config.Repeat)
		app.showMenu(q.Message.Chat.ID, q.Message.MessageID, selected, repMenuText, repKeyboard)
	case match(q.Data, alphaButtonCallback):
		selected := fmt.Sprintf("%s/%d", alphaButtonCallback, s.Config.Alpha)
		app.showMenu(q.Message.Chat.ID, q.Message.MessageID, selected, alphaMenuText, alphaKeyboard)
	case match(q.Data, fmt.Sprintf("%s/([0-9]+)", alphaButtonCallback), &num):
		if num < 0 || num > 255 {
			return
		}
		s.Config.Alpha = num
		app.sessions.Set(q.From.ID, s)

		// update menu
		selected := fmt.Sprintf("%s/%d", alphaButtonCallback, s.Config.Alpha)
		app.showMenu(q.Message.Chat.ID, q.Message.MessageID, selected, alphaMenuText, alphaKeyboard)
	case match(q.Data, alphaInputCallback):
		app.handleAlphaInput(q, s)
	case match(q.Data, extButtonCallback):
		selected := fmt.Sprintf("%s/%s", extButtonCallback, s.Config.Extension)
		app.showMenu(q.Message.Chat.ID, q.Message.MessageID, selected, extMenuText, extKeyboard)
	case match(q.Data, fmt.Sprintf("%s/(jpg|png|svg|gif)", extButtonCallback), &slug):
		s.Config.Extension = primitive.Extension(slug)
		app.sessions.Set(q.From.ID, s)

		// update menu
		selected := fmt.Sprintf("%s/%s", extButtonCallback, s.Config.Extension)
		app.showMenu(q.Message.Chat.ID, q.Message.MessageID, selected, extMenuText, extKeyboard)
	case match(q.Data, sizeButtonCallback):
		selected := fmt.Sprintf("%s/%d", sizeButtonCallback, s.Config.OutputSize)
		app.showMenu(q.Message.Chat.ID, q.Message.MessageID, selected, sizeMenuText, sizeKeyboard)
	case match(q.Data, fmt.Sprintf("%s/([0-9]+)", sizeButtonCallback), &num):
		if num < 256 || num > 1920 {
			return
		}
		s.Config.OutputSize = num
		app.sessions.Set(q.From.ID, s)

		// update menu
		selected := fmt.Sprintf("%s/%d", sizeButtonCallback, s.Config.OutputSize)
		app.showMenu(q.Message.Chat.ID, q.Message.MessageID, selected, sizeMenuText, sizeKeyboard)
	case match(q.Data, sizeInputCallback):
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

	buttonText := fmt.Sprintf("%s (%d)", otherButtonText, s.Config.Iterations)
	app.showMenu(q.Message.Chat.ID, q.Message.MessageID, iterInputCallback,
		iterMenuText, iterKeyboard, buttonText)
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

	buttonText := fmt.Sprintf("%s (%d)", otherButtonText, s.Config.Alpha)
	app.showMenu(q.Message.Chat.ID, q.Message.MessageID, alphaInputCallback,
		alphaMenuText, alphaKeyboard, buttonText)
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

	buttonText := fmt.Sprintf("%s (%d)", otherButtonText, s.Config.OutputSize)
	app.showMenu(q.Message.Chat.ID, q.Message.MessageID, sizeInputCallback,
		sizeMenuText, sizeKeyboard, buttonText)
}

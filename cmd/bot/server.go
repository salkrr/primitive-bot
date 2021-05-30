package main

import (
	"fmt"
	"os"
	"time"

	"github.com/lazy-void/primitive-bot/pkg/menu"
	"github.com/lazy-void/primitive-bot/pkg/sessions"
	"github.com/lazy-void/primitive-bot/pkg/tg"
)

func (app *application) listenAndServe() {
	go app.worker()

	offset := int64(0)
	for {
		updates, err := app.bot.GetUpdates(offset, 100, 20, []string{"message", "callback_query"})
		if err != nil {
			app.errorLog.Print(err)
			continue
		}

		numUpdates := len(updates)
		if numUpdates == 0 {
			continue
		}

		for _, u := range updates {
			if u.Message.MessageID > 0 {
				app.infoLog.Printf("Message: text '%s' from the user '%s' with the ID '%d'",
					u.Message.Text, u.Message.From.FirstName, u.Message.From.ID)
				go app.processMessage(u.Message)
				continue
			}

			app.infoLog.Printf("Callback Query: data '%s' from the user '%s' with the ID '%d'",
				u.CallbackQuery.Data, u.CallbackQuery.From.FirstName, u.CallbackQuery.From.ID)
			go app.processCallbackQuery(u.CallbackQuery)
		}

		offset = updates[numUpdates-1].UpdateID + 1
	}
}

func (app *application) worker() {
	for {
		// If we'll delete the operation from the queue with Dequeue
		// then if user will use command '/status', we won't be able to
		// send him any information about this operation (it is not in the
		// queue so we have no idea if it exists or not)
		op, ok := app.queue.Peek()
		if !ok {
			time.Sleep(1 * time.Second)
			continue
		}

		// create primitive
		start := time.Now()
		outputPath := fmt.Sprintf("%s/%d_%d.%s", app.outDir, op.UserID, start.Unix(), op.Config.Extension)
		app.infoLog.Printf(creatingLogMessage, op.UserID, op.ImgPath, outputPath, op.Config.Iterations, op.Config.Shape,
			op.Config.Alpha, op.Config.Repeat, op.Config.OutputSize, op.Config.Extension)

		err := op.Config.Create(op.ImgPath, outputPath)
		if err != nil {
			app.serverError(op.UserID, err)
			return
		}
		app.infoLog.Printf(finishedLogMessage, op.UserID, op.ImgPath, outputPath, time.Since(start).Seconds())

		// send output to the user
		err = app.bot.SendDocument(op.UserID, outputPath)
		if err != nil {
			app.serverError(op.UserID, err)
			return
		}
		app.infoLog.Printf(sentLogMessage, op.UserID, outputPath)

		// remove operation from the queue
		app.queue.Dequeue()
	}
}

func (app *application) processMessage(m tg.Message) {
	if m.Photo != nil {
		app.processPhoto(m)
		return
	}

	// Handle user input if they are inside the input form
	s, ok := app.sessions.Get(m.From.ID)
	if ok && s.State == sessions.InInputDialog {
		s.Input <- m
		return
	}

	if m.Text == "/status" {
		operations, positions := app.queue.GetOperations(m.From.ID)
		if len(operations) == 0 {
			_, err := app.bot.SendMessage(m.Chat.ID,
				app.printer.Sprint("There aren't any operations in the queue."))
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
	_, err := app.bot.SendMessage(m.Chat.ID,
		app.printer.Sprintf("Send me some image."))
	if err != nil {
		app.serverError(m.Chat.ID, err)
	}
}

func (app *application) processPhoto(m tg.Message) {
	// If we already have session - delete it's menu
	s, ok := app.sessions.Get(m.From.ID)
	if ok {
		err := app.bot.DeleteMessage(s.UserID, s.MenuMessageID)
		if err != nil {
			app.serverError(m.Chat.ID, err)
			return
		}
	}

	path, err := app.downloadPhoto(m.Photo)
	if err != nil {
		app.serverError(m.Chat.ID, err)
		return
	}

	// Create session
	msg, err := app.bot.SendMessage(m.Chat.ID, menu.RootViewTmpl.Text, menu.RootViewTmpl.Keyboard)
	if err != nil {
		app.serverError(m.Chat.ID, err)
		return
	}
	app.sessions.Set(m.From.ID, sessions.NewSession(m.From.ID, msg.MessageID, path, app.workers))
}

func (app *application) downloadPhoto(photos []tg.PhotoSize) (string, error) {
	// Choose smallest image with dimensions >= 256
	var file tg.PhotoSize
	for _, photo := range photos {
		file = photo
		if photo.Width >= 256 && photo.Height >= 256 {
			break
		}
	}
	if file.FileID == "" {
		return "", fmt.Errorf("no image files in %v", photos)
	}

	path := fmt.Sprintf("%s/%s.jpg", app.inDir, file.FileUniqueID)
	// Download the file only if we don't have it
	if _, err := os.Stat(path); os.IsNotExist(err) {
		img, err := app.bot.DownloadFile(file.FileID)
		if err != nil {
			return "", fmt.Errorf("couldn't download image: %w", err)
		}

		if err := os.WriteFile(path, img, 0600); err != nil {
			return "", fmt.Errorf("couldn't save image: %w", err)
		}
	}

	return path, nil
}

func (app *application) processCallbackQuery(q tg.CallbackQuery) {
	defer func() {
		err := app.bot.AnswerCallbackQuery(q.ID, "")
		if err != nil {
			app.errorLog.Printf("Error answering callback query: %s", err)
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
	case match(q.Data, menu.RootViewCallback):
		app.showRootMenuView(s)
	case match(q.Data, menu.CreateButtonCallback):
		app.handleCreateButton(s)
	case match(q.Data, menu.ShapesViewCallback):
		app.showShapesMenuView(s)
	case match(q.Data, menu.ShapesButtonCallback, &num):
		app.handleShapesButton(s, num)
	case match(q.Data, menu.IterViewCallback):
		app.showIterMenuView(s)
	case match(q.Data, menu.IterButtonCallback, &num):
		app.handleIterButton(s, num)
	case match(q.Data, menu.IterInputCallback):
		app.handleIterInput(s)
	case match(q.Data, menu.RepViewCallback):
		app.showRepMenuView(s)
	case match(q.Data, menu.RepButtonCallback, &num):
		app.handleRepButton(s, num)
	case match(q.Data, menu.AlphaViewCallback):
		app.showAlphaMenuView(s)
	case match(q.Data, menu.AlphaButtonCallback, &num):
		app.handleAlphaButton(s, num)
	case match(q.Data, menu.AlphaInputCallback):
		app.handleAlphaInput(s)
	case match(q.Data, menu.ExtViewCallback):
		app.showExtMenuView(s)
	case match(q.Data, menu.ExtButtonCallback, &slug):
		app.handleExtButton(s, slug)
	case match(q.Data, menu.SizeViewCallback):
		app.showSizeMenuView(s)
	case match(q.Data, menu.SizeButtonCallback, &num):
		app.handleSizeButton(s, num)
	case match(q.Data, menu.SizeInputCallback):
		app.handleSizeInput(s)
	}
}

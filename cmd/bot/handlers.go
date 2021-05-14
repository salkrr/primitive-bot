package main

import (
	"fmt"
	"os"
	"runtime"
	"time"

	"github.com/fogleman/primitive/primitive"
	"github.com/lazy-void/primitive-bot/pkg/telegram"
)

type operation struct {
	chatID    int64
	imagePath string
	config    config
}

func (app *application) handleMessage(m telegram.Message, out chan operation) {
	if m.Photo == nil {
		err := app.bot.SendMessage(m.Chat.ID, messageHelp)
		if err != nil {
			app.serverError(m.Chat.ID, err)
		}
		return
	}

	// Download image
	img, err := app.bot.DownloadFile(m.Photo[2].FileID)
	if err != nil {
		app.serverError(m.Chat.ID, fmt.Errorf("couldn't download image: %s", err))
		return
	}
	path := fmt.Sprintf("./inputs/%d_%d_orig.jpg", m.From.ID, time.Now().Unix())
	os.WriteFile(path, img, 0644)
	out <- operation{
		chatID:    m.Chat.ID,
		imagePath: path,
		config: config{
			workers:    runtime.NumCPU() - 1,
			outputSize: 1920,
			shape:      primitive.ShapeTypeAny,
			iterations: 100,
			repeat:     1,
			alpha:      128,
			ext:        png,
		},
	}
	app.bot.SendMessage(m.Chat.ID, messageQueued)
}

func (app *application) primitiveWorker(in chan operation) {
	var op operation
	for {
		// get next image from queue
		op = <-in

		app.infoLog.Printf("Starting to create primitive from the image '%s'", op.imagePath)
		// create primitive
		outputPath := fmt.Sprintf("./outputs/%d_%d.%s", op.chatID, time.Now().Unix(), op.config.ext)
		err := op.config.create(op.imagePath, outputPath)
		if err != nil {
			app.serverError(op.chatID, err)
			return
		}
		app.infoLog.Printf("Finished creating primitive from the image '%s'", op.imagePath)

		// send output to the user
		err = app.bot.SendPhoto(op.chatID, outputPath)
		if err != nil {
			app.serverError(op.chatID, err)
			return
		}
		app.infoLog.Printf("Sent result to the %d", op.chatID)
	}
}

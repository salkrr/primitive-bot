package main

import (
	"flag"
	"log"
	"os"

	"github.com/lazy-void/primitive-bot/pkg/queue"
	"github.com/lazy-void/primitive-bot/pkg/sessions"
	"github.com/lazy-void/primitive-bot/pkg/telegram"
)

type application struct {
	infoLog         *log.Logger
	errorLog        *log.Logger
	inDir           string
	outDir          string
	operationsLimit int
	bot             *telegram.Bot
	sessions        *sessions.ActiveSessions
	queue           *queue.Queue
}

func main() {
	token := flag.String("token", "", "The token for the Telegram Bot")
	inDir := flag.String("i", "inputs", "Name of the directory where inputs will be stored")
	outDir := flag.String("o", "outputs", "Name of the directory where outputs will be stored")
	operationsLimit := flag.Int("limit", 5, "The number of operations that the user can add to the queue.")
	flag.Parse()

	if *token == "" {
		log.Fatal("You need to provide the token for the Telegram Bot!")
	}

	if err := os.MkdirAll(*inDir, 0664); err != nil {
		log.Fatal(err)
	}
	if err := os.MkdirAll(*outDir, 0664); err != nil {
		log.Fatal(err)
	}

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	app := application{
		infoLog:         infoLog,
		errorLog:        errorLog,
		inDir:           *inDir,
		outDir:          *outDir,
		operationsLimit: *operationsLimit,
		bot:             &telegram.Bot{Token: *token},
		sessions:        sessions.NewActiveSessions(),
		queue:           queue.New(),
	}

	infoLog.Printf("Starting to listen for updates...")
	app.listenAndServe()
}

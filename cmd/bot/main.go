package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/lazy-void/primitive-bot/pkg/primitive"

	"github.com/lazy-void/primitive-bot/pkg/queue"
	"github.com/lazy-void/primitive-bot/pkg/sessions"
	"github.com/lazy-void/primitive-bot/pkg/tg"
)

type application struct {
	infoLog         *log.Logger
	errorLog        *log.Logger
	inDir           string
	outDir          string
	operationsLimit int
	bot             *tg.Bot
	sessions        *sessions.ActiveSessions
	queue           *queue.Queue
}

func main() {
	token := flag.String("token", "", "The token for the Telegram Bot.")
	inDir := flag.String("i", "inputs", "Name of the directory where inputs will be stored.")
	outDir := flag.String("o", "outputs", "Name of the directory where outputs will be stored.")
	logPath := flag.String("log", "", "Path to the previous log file. It is used to restore queue.")
	operationsLimit := flag.Int("limit", 5, "The number of operations that the user can add to the queue.")
	flag.Parse()

	if *token == "" {
		log.Fatal("You need to provide the token for the Telegram Bot!")
	}

	q := queue.New()
	if *logPath != "" {
		err := restoreQueue(*logPath, q)
		if err != nil {
			log.Fatalf("Error restoring queue from the log: %s", err)
		}
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
		bot:             &tg.Bot{Token: *token},
		sessions:        sessions.NewActiveSessions(),
		queue:           q,
	}

	infoLog.Printf("Starting to listen for updates...")
	app.listenAndServe()
}

func restoreQueue(logPath string, q *queue.Queue) (err error) {
	f, err := os.Open(filepath.Clean(logPath))
	if err != nil {
		return err
	}
	defer func() {
		err = f.Close()
	}()

	regex := regexp.MustCompile(`INFO\t\b.*?\b \b.*?\b (.*)`)
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		// Remove log prefix
		matches := regex.FindStringSubmatch(scanner.Text())
		if len(matches) < 2 {
			continue
		}

		msg := matches[1]
		if strings.HasPrefix(msg, "Enqueued:") {
			op := queue.Operation{Config: primitive.NewConfig()}
			_, err := fmt.Sscanf(msg, enqueuedLogMessage, &op.UserID, &op.ImgPath,
				&op.Config.Iterations, &op.Config.Shape, &op.Config.Alpha,
				&op.Config.Repeat, &op.Config.OutputSize, &op.Config.Extension)
			if err != nil {
				return err
			}

			q.Enqueue(op)
			continue
		}

		if strings.HasPrefix(msg, "Finished:") {
			q.Dequeue()
		}
	}

	return scanner.Err()
}

package main

import (
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/fuzzingbits/canbot/pkg/canbot"
	"github.com/fuzzingbits/forge-wip/pkg/gol"
	"github.com/rollbar/rollbar-go"
)

var wg sync.WaitGroup

func main() {
	// Listen for OS signals
	signalChan := make(chan os.Signal)
	signal.Notify(signalChan, syscall.SIGINT)
	signal.Notify(signalChan, os.Interrupt)
	go func() {
		<-signalChan
		wg.Wait()
		os.Exit(0)
	}()

	// Get the primary logger
	logger := getLogger()

	// Create a new app
	app, err := canbot.NewApp(logger)
	if err != nil {
		logger.Error(err)
		os.Exit(1)
	}

	// First run of the app
	run(app)

	// Run at an interval
	for range time.NewTicker(time.Minute * 1).C {
		run(app)
	}
}

func run(app *canbot.App) {
	wg.Add(1)
	defer wg.Done()

	app.Interval()
}

func getLogger() gol.Logger {
	logger := log.New(os.Stderr, "", log.LstdFlags)

	if rollbarToken := os.Getenv("ROLLBAR_TOKEN"); rollbarToken != "" {
		return &gol.RollbarLogger{
			Logger:  logger,
			Rollbar: rollbar.New(rollbarToken, "prod", "", "", ""),
		}
	}

	return &gol.LogLogger{
		Logger: logger,
	}
}

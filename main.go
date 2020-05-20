package main

import (
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/fuzzingbits/canbot/pkg/canbot"
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

	// Create a new app
	app, err := canbot.NewApp()
	if err != nil {
		log.Fatal(err)
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

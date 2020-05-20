package canbot

import (
	"errors"
	"strings"

	"github.com/fuzzingbits/canbot/pkg/internal/slack"
	"github.com/fuzzingbits/forge-wip/pkg/config"
)

// NewApp reads env vars and creates a new App incstance
func NewApp() (*App, error) {
	configInstance := config.Config{
		Providers: []config.Provider{
			config.ProviderEnvironment{},
		},
	}

	app := &App{
		// Define defaults
		SlackUsername:  "CanBot",
		SlackIconEmoji: ":flushed",
		StatePath:      "state.json",
		// More setup
		state: &state{},
	}

	_ = configInstance.Unmarshal(app)

	app.SlackTargetsExpanded = strings.Split(app.SlackTargets, ",")
	app.slackAPI = &slack.Client{Token: app.SlackToken}
	app.state.filename = app.StatePath

	// Read in existing app state
	app.state.Read()

	// Validate the app
	if err := app.validate(); err != nil {
		return &App{}, err
	}

	return app, nil
}

func (app *App) validate() error {
	if len(app.SlackTargets) < 1 {
		return errors.New("SlackTargets can not be blank")
	}

	if len(app.SlackToken) < 1 {
		return errors.New("SlackToken can not be blank")
	}

	return nil
}

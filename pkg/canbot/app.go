package canbot

import (
	"fmt"

	"github.com/fuzzingbits/canbot/pkg/internal/slack"
	"github.com/fuzzingbits/forge-wip/pkg/gol"
)

// App is the canbot app
type App struct {
	SlackToken           string `env:"SLACK_TOKEN"`
	SlackUsername        string `env:"SLACK_USERNAME"`
	SlackIconEmoji       string `env:"SLACK_ICON_EMOJI"`
	SlackTargets         string `env:"SLACK_TARGETS"`
	SlackTargetsExpanded []string
	StatePath            string `env:"STATE_FILE"`
	slackAPI             slack.Service
	state                *state
	logger               gol.Logger
}

// Interval checks for new users and logs any errors
func (app *App) Interval() {
	err := app.checkUsers()
	if err != nil {
		app.logger.Error(err)
	}
}

func (app *App) checkUsers() error {
	users, err := app.slackAPI.UsersList()
	if err != nil {
		return err
	}

	app.logger.Log("Users Query Complete: %d found", len(users))

	for _, user := range users {
		// Skip all bot users
		if user.IsBot {
			continue
		}

		wasPreviouslyDeleted := app.state.IsUserAlreadyDeleted(user)

		// Skip deleted users that were deleted last run
		if user.Deleted && wasPreviouslyDeleted {
			continue
		}

		// Skip non-deleted that were not deleted last run
		if !user.Deleted && !wasPreviouslyDeleted {
			continue
		}

		// At this point we know that the users status has changed since the last run

		// For when users are newly un-deleted
		if !user.Deleted {
			app.state.AddPendingAlert(user)
			app.state.RemoveDeletedUser(user)
			app.logger.Log("New Un-Deleted User Found: %s", user.ID)

			continue
		}

		// For when users are newly deleted
		if user.Deleted {
			app.state.AddDeletedUser(user)

			if app.state.FirstRunComplete {
				app.state.AddPendingAlert(user)
				app.logger.Log("New Deleted User Found: %s", user.ID)
			}
		}
	}

	if err := app.sendAlert(); err != nil {
		return err
	}

	app.state.FirstRunComplete = true
	_ = app.state.Write()

	return nil
}

func (app *App) sendAlert() error {
	pendingAlerts := app.state.GetAllPendingAlerts()
	if len(pendingAlerts) > 0 && app.state.FirstRunComplete {
		app.logger.Log("About to alert on new Deleted User(s): %d", len(pendingAlerts))
	}

	for _, user := range pendingAlerts {
		for _, target := range app.SlackTargetsExpanded {
			if _, err := app.slackAPI.ChatPostMessage(slack.Message{
				Channel:   target,
				Text:      app.generateUserMessage(user),
				Username:  app.SlackUsername,
				IconEmoji: app.SlackIconEmoji,
			}); err != nil {
				// TODO: check the type of error
				// channel/user not found: log error and continue
				// any other error: return the error
				return err
			}

			app.logger.Log("Alerted %s about user: %s", target, user.ID)

			// Remove from pending alerts so we don't try to alert again
			app.state.RemovePendingAlert(user)
		}
	}

	return nil
}

func (app *App) generateUserMessage(user slack.User) string {
	if !user.Deleted {
		return fmt.Sprintf("Slack account reactivated: <@%s>", user.ID)
	}

	return fmt.Sprintf("Slack account deactivated: <@%s>", user.ID)
}

package canbot

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"testing"

	"github.com/fuzzingbits/canbot/pkg/internal/slack"
	"github.com/fuzzingbits/forge-wip/pkg/gol"
)

type mockSlack struct {
	UserListResponse        func() ([]slack.User, error)
	ChatPostMessageResponse func() (slack.Message, error)
}

func (m *mockSlack) UsersList() ([]slack.User, error) {
	return m.UserListResponse()
}

func (m *mockSlack) ChatPostMessage(message slack.Message) (slack.Message, error) {
	return m.ChatPostMessageResponse()
}

var golLogger = &gol.LogLogger{
	Logger: log.New(os.Stderr, "", log.LstdFlags),
}

func TestMain(t *testing.T) {
	resetTestState()
	defer testStateCleanup()

	os.Setenv("STATE_FILE", "state_git_excluded.json")
	os.Setenv("SLACK_TOKEN", "fake-token")
	os.Setenv("SLACK_TARGETS", "target1,target2")
	os.Setenv("SLACK_ICON_EMOJI", ":lol:")

	var postMessageError error
	sentMessageCount := 0
	slackUsers := []slack.User{
		{ID: "user1", Deleted: true},
		{ID: "user2", Deleted: false},
		{ID: "user3", Deleted: false},
		{ID: "user4", Deleted: true, IsBot: true},
	}

	slackAPI := &mockSlack{
		UserListResponse: func() ([]slack.User, error) {
			return slackUsers, nil
		},
		ChatPostMessageResponse: func() (slack.Message, error) {
			if postMessageError != nil {
				return slack.Message{}, postMessageError
			}

			sentMessageCount++
			return slack.Message{}, nil
		},
	}

	app, err := NewApp(golLogger)
	if err != nil {
		t.Error(err)
	}

	if err := app.state.Clear(); err != nil {
		t.Error(err)
	}

	app.slackAPI = slackAPI

	{ // First Run
		app.Interval()
		if sentMessageCount > 0 {
			t.Errorf("No messages should be sent on the first run, not %d", sentMessageCount)
		}
		sentMessageCount = 0
	}

	{ // Second Run
		slackUsers[1].Deleted = true
		app.Interval()
		if sentMessageCount != 2 {
			t.Errorf("2 messages should be sent on the second run, not %d", sentMessageCount)
		}
		sentMessageCount = 0
	}

	{ // Third Run
		app.Interval()
		if sentMessageCount > 0 {
			t.Errorf("No messages should be sent on the third run, not %d", sentMessageCount)
		}
		sentMessageCount = 0
	}

	{ // Fourth Run
		slackUsers[1].Deleted = false
		app.Interval()
		if sentMessageCount != 2 {
			t.Errorf("2 messages should be sent on the fourth run, not %d", sentMessageCount)
		}
		sentMessageCount = 0
	}

	{ // Fifth Run
		app.Interval()
		if sentMessageCount > 0 {
			t.Errorf("No messages should be sent on the fifth run, not %d", sentMessageCount)
		}
		sentMessageCount = 0
	}

	{ // Sixth Run
		postMessageError = errors.New("foobar")
		slackUsers[1].Deleted = true
		app.Interval()
		if sentMessageCount > 0 {
			t.Errorf("No messages should be sent on the sixth run, not %d", sentMessageCount)
		}
		sentMessageCount = 0
	}
	{ // Seventh Run
		postMessageError = nil
		app.Interval()
		if sentMessageCount != 2 {
			t.Errorf("2 messages should be sent on the seventh run, not %d", sentMessageCount)
		}
		sentMessageCount = 0
	}

	slackUsers[1].Deleted = false

	newApp, _ := NewApp(golLogger)
	newApp.slackAPI = slackAPI
	{ // Eighth Run
		postMessageError = nil
		newApp.Interval()
		if sentMessageCount != 2 {
			t.Errorf("2 messages should be sent on the eighth run, not %d", sentMessageCount)
		}
		sentMessageCount = 0
	}
}

func TestBadStateFile(t *testing.T) {
	resetTestState()
	defer testStateCleanup()
	jsonFilePath := "malformed_state_git_excluded.json"
	os.Setenv("STATE_FILE", jsonFilePath)
	ioutil.WriteFile(jsonFilePath, []byte("{true: true},"), 0644)

	_, err := NewApp(golLogger)
	if err == nil {
		t.Errorf("There should have been an error")
	}
}

func TestDefaults(t *testing.T) {
	resetTestState()
	defer testStateCleanup()

	os.Setenv("SLACK_TOKEN", "fake-token")
	os.Setenv("SLACK_TARGETS", "target1,target2")

	app, err := NewApp(golLogger)
	if err != nil {
		t.Error("there should not have been an error")
	}

	if app.SlackUsername == "" {
		t.Error("Missing default for SlackUsername")
	}

	if app.SlackIconEmoji == "" {
		t.Error("Missing default for SlackIconEmoji")
	}

	if app.StatePath == "" {
		t.Error("Missing default for StatePath")
	}
}

func TestMainNoTokenError(t *testing.T) {
	resetTestState()
	defer testStateCleanup()

	os.Setenv("SLACK_TARGETS", "target1,target2")

	_, err := NewApp(golLogger)
	if err == nil {
		t.Error("there should have been an error")
	}
}

func TestMainNoTargetsError(t *testing.T) {
	resetTestState()
	defer testStateCleanup()

	os.Setenv("SLACK_TOKEN", "fake-token")

	_, err := NewApp(golLogger)
	if err == nil {
		t.Error("there should have been an error")
	}
}

func TestMainBlankTokenError(t *testing.T) {
	resetTestState()
	defer testStateCleanup()

	os.Setenv("SLACK_TOKEN", "")

	_, err := NewApp(golLogger)
	if err == nil {
		t.Error("there should have been an error")
	}
}

func TestMainUsersError(t *testing.T) {
	resetTestState()
	defer testStateCleanup()

	os.Setenv("SLACK_TOKEN", "fake-token")
	os.Setenv("SLACK_TARGETS", "target1,target2")

	app, err := NewApp(golLogger)
	if err != nil {
		t.Error(err)
	}

	slackAPI := &mockSlack{
		UserListResponse: func() ([]slack.User, error) {
			return []slack.User{}, fmt.Errorf("Error pulling users")
		},
	}

	app.slackAPI = slackAPI
	app.Interval()
}

func resetTestState() {
	testStateCleanup()
	os.Unsetenv("SLACK_ICON_EMOJI")
	os.Unsetenv("SLACK_USERNAME")
	os.Unsetenv("STATE_FILE")
	os.Unsetenv("SLACK_TOKEN")
	os.Unsetenv("SLACK_TARGETS")
}

func testStateCleanup() {
	os.Remove("state.json")
}

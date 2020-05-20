package canbot

import (
	"encoding/json"
	"io/ioutil"
	"os"

	"github.com/fuzzingbits/canbot/pkg/internal/slack"
)

type state struct {
	DeletedUsers     []string
	PendingAlerts    map[string]slack.User
	FirstRunComplete bool
	filename         string
}

func (s *state) IsUserAlreadyDeleted(user slack.User) bool {
	for _, userID := range s.DeletedUsers {
		if userID == user.ID {
			return true
		}
	}

	return false
}

func (s *state) AddDeletedUser(user slack.User) {
	s.DeletedUsers = append(s.DeletedUsers, user.ID)

	s.Write()
}

func (s *state) RemoveDeletedUser(user slack.User) {
	for i, userID := range s.DeletedUsers {
		if userID == user.ID {
			s.DeletedUsers = append(s.DeletedUsers[:i], s.DeletedUsers[i+1:]...)
			break
		}
	}

	s.Write()
}

func (s *state) AddPendingAlert(user slack.User) {
	s.PendingAlerts[user.ID] = user

	s.Write()
}

func (s *state) GetAllPendingAlerts() []slack.User {
	users := []slack.User{}
	for _, user := range s.PendingAlerts {
		users = append(users, user)
	}

	return users
}

func (s *state) RemovePendingAlert(user slack.User) {
	delete(s.PendingAlerts, user.ID)
	s.Write()
}

func (s *state) Read() error {
	configFileBytes, err := ioutil.ReadFile(s.filename)
	if os.IsNotExist(err) {
		// If no file exists, clear the state. This creates the file and writes a blank state to it
		return s.Clear()
	}

	err = json.Unmarshal(configFileBytes, s)
	if err != nil {
		return s.Clear()
	}

	return s.Write()
}

func (s *state) Clear() error {
	s.FirstRunComplete = false
	s.DeletedUsers = []string{}
	s.PendingAlerts = make(map[string]slack.User)

	return s.Write()
}

func (s *state) Write() error {
	configBytes, _ := json.Marshal(s)

	return ioutil.WriteFile(s.filename, configBytes, 0644)
}

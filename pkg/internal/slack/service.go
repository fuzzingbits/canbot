package slack

import "net/http"

// Service is the Slack API Service definition
type Service interface {
	UsersList() ([]User, error)
	ChatPostMessage(message Message) (Message, error)
}

// UsersList returns a list of all Slack Users
func (c *Client) UsersList() ([]User, error) {
	target := usersResponse{}

	err := c.call(http.MethodGet, "users.list", nil, &target)
	if err != nil {
		return []User{}, err
	}

	return target.Members, nil
}

// ChatPostMessage posts a chat message
func (c *Client) ChatPostMessage(message Message) (Message, error) {
	target := messageResponse{}

	err := c.call(http.MethodPost, "chat.postMessage", message, &target)
	if err != nil {
		return Message{}, err
	}

	return target.Message, nil
}

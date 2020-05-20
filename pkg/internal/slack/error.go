package slack

import "fmt"

// ErrorResponse is the structure of a Slack API error response
type ErrorResponse struct {
	Ok           bool   `json:"ok"`
	ErrorMessage string `json:"error"`
}

// Error satisfies the error interface
func (e ErrorResponse) Error() string {
	return fmt.Sprintf("Slack API Error: [%s]", e.ErrorMessage)
}

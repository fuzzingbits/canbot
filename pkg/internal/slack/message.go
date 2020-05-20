package slack

// Message is a Slack Message
type Message struct {
	Username  string `json:"username"`
	IconEmoji string `json:"icon_emoji"`
	Channel   string `json:"channel"`
	Text      string `json:"text"`
}

type messageResponse struct {
	OK      bool    `json:"ok"`
	Channel string  `json:"channel"`
	Ts      string  `json:"ts"`
	Message Message `json:"message"`
}

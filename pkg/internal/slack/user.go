package slack

// User is a Slack User
// https://api.slack.com/types/user
type User struct {
	ID                string  `json:"id"`
	TeamID            string  `json:"team_id"`
	Name              string  `json:"name"`
	Deleted           bool    `json:"deleted"`
	Color             string  `json:"color"`
	RealName          string  `json:"real_name"`
	Tz                string  `json:"tz"`
	TzLabel           string  `json:"tz_label"`
	TzOffset          int     `json:"tz_offset"`
	Profile           profile `json:"profile"`
	IsAdmin           bool    `json:"is_admin"`
	IsOwner           bool    `json:"is_owner"`
	IsPrimaryOwner    bool    `json:"is_primary_owner"`
	IsRestricted      bool    `json:"is_restricted"`
	IsUltraRestricted bool    `json:"is_ultra_restricted"`
	IsBot             bool    `json:"is_bot"`
	Updated           int     `json:"updated"`
	IsAppUser         bool    `json:"is_app_user"`
	Has2Fa            bool    `json:"has_2fa"`
}

type profile struct {
	AvatarHash            string `json:"avatar_hash"`
	Title                 string `json:"title"`
	StatusText            string `json:"status_text"`
	StatusEmoji           string `json:"status_emoji"`
	RealName              string `json:"real_name"`
	DisplayName           string `json:"display_name"`
	RealNameNormalized    string `json:"real_name_normalized"`
	DisplayNameNormalized string `json:"display_name_normalized"`
	Email                 string `json:"email"`
	Image24               string `json:"image_24"`
	Image32               string `json:"image_32"`
	Image48               string `json:"image_48"`
	Image72               string `json:"image_72"`
	Image192              string `json:"image_192"`
	Image512              string `json:"image_512"`
	Team                  string `json:"team"`
}

type usersResponse struct {
	OK      bool   `json:"ok"`
	Members []User `json:"members"`
	CacheTs int    `json:"cache_ts"`
}

package slack

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/fuzzingbits/forge-wip/pkg/rest"
)

// Client is the main Slack API client
type Client struct {
	Token      string
	RestClient *rest.Client
}

func (c *Client) call(method string, endpoint string, payload interface{}, target interface{}) error {
	c.setup()

	urlEndpoint := &url.URL{
		Scheme: "https",
		Host:   "slack.com",
		Path:   fmt.Sprintf("/api/%s", endpoint),
	}

	return c.RestClient.CurlSimple(method, urlEndpoint, payload, target)
}

func (c *Client) setup() {
	if c.RestClient == nil {
		c.RestClient = &rest.Client{
			Service: restService{token: c.Token},
			HTTPClient: &http.Client{
				Transport: nil,
			},
		}
	}
}

type restService struct {
	token string
}

func (s restService) ModRequest(request *http.Request) error {
	request.Header.Set(
		"Authorization",
		fmt.Sprintf("Bearer %s", s.token),
	)

	request.Header.Set(
		"Content-Type",
		"application/json;charset=utf-8",
	)

	return nil
}

func (s restService) ErrorCheck(responseBytes []byte) error {
	slackError := ErrorResponse{}
	if err := json.Unmarshal(responseBytes, &slackError); err != nil {
		return err
	}

	if slackError.ErrorMessage != "" {
		return slackError
	}

	return nil
}

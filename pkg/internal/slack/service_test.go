package slack

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"reflect"
	"testing"

	"github.com/fuzzingbits/forge-wip/pkg/rest"
)

func TestClient_UsersList(t *testing.T) {
	testClient.setup()
	testClient.RestClient = &rest.Client{
		HTTPClient: testHTTPClient,
		Service:    testService,
	}

	tests := []struct {
		name    string
		fields  *Client
		want    []User
		wantErr bool
		setup   func()
	}{
		{
			name:   "Primary test",
			fields: testClient,
			want: []User{
				{
					ID:   "ABC123",
					Name: "Aaron",
				},
			},
			wantErr: false,
			setup: func() {
				transport.RoundTripFunc = func(req *http.Request) (*http.Response, error) {
					return &http.Response{
						StatusCode: http.StatusOK,
						Body:       ioutil.NopCloser(bytes.NewBufferString(`{"members": [{"id": "ABC123", "name": "Aaron"}]}`)),
						Header:     make(http.Header),
					}, nil
				}
			},
		},
		{
			name:    "Slack Error test",
			fields:  testClient,
			want:    []User{},
			wantErr: true,
			setup: func() {
				transport.RoundTripFunc = func(req *http.Request) (*http.Response, error) {
					return &http.Response{
						StatusCode: http.StatusOK,
						Body:       ioutil.NopCloser(bytes.NewBufferString(`{"error": "not authed or something"}`)),
						Header:     make(http.Header),
					}, nil
				}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setup != nil {
				tt.setup()
			}

			got, err := tt.fields.UsersList()
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.UsersList() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Client.UsersList() = %v, want %v", got, tt.want)
			}

			{ // Cleanup/Reset
				transport.RoundTripFunc = nil
			}
		})
	}
}

func TestClient_ChatPostMessage(t *testing.T) {
	testClient.setup()
	testClient.RestClient = &rest.Client{
		HTTPClient: testHTTPClient,
		Service:    testService,
	}

	type args struct {
		message Message
	}
	tests := []struct {
		name    string
		fields  *Client
		args    args
		want    Message
		wantErr bool
		setup   func()
	}{
		{
			name:   "Primary test",
			fields: testClient,
			args: args{
				message: Message{},
			},
			want:    Message{Channel: "FAKECHANEL"},
			wantErr: false,
			setup: func() {
				transport.RoundTripFunc = func(req *http.Request) (*http.Response, error) {
					return &http.Response{
						StatusCode: http.StatusOK,
						Body:       ioutil.NopCloser(bytes.NewBufferString(`{"message": {"channel": "FAKECHANEL"}}`)),
						Header:     make(http.Header),
					}, nil
				}
			},
		},
		{
			name:   "Slack Error test",
			fields: testClient,
			args: args{
				message: Message{},
			},
			want:    Message{},
			wantErr: true,
			setup: func() {
				transport.RoundTripFunc = func(req *http.Request) (*http.Response, error) {
					return &http.Response{
						StatusCode: http.StatusOK,
						Body:       ioutil.NopCloser(bytes.NewBufferString(`{"error": "not authed or something"}`)),
						Header:     make(http.Header),
					}, nil
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setup != nil {
				tt.setup()
			}

			got, err := tt.fields.ChatPostMessage(tt.args.message)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.ChatPostMessage() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Client.ChatPostMessage() = %v, want %v", got, tt.want)
			}

			{ // Cleanup/Reset
				transport.RoundTripFunc = nil
			}
		})
	}
}
